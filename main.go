package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	colorize "github.com/fatih/color"
	ignore "github.com/sabhiram/go-gitignore"
)

var (
	ignoredColor  = colorize.FgGreen
	includedColor = colorize.FgHiRed

	// defaultNoColor captures fatih/color's TTY/NO_COLOR/TERM auto-detection
	// at package init, so CLICOLOR_FORCE can override it and later restore it
	// on the next invocation (relevant across repeated realMain() calls in tests).
	defaultNoColor = colorize.NoColor
)

type Config struct {
	Color        bool
	Debug        bool
	Excluded     bool
	Fullpath     bool
	Ignorefile   string
	Included     bool
	Invertcolors bool
	Nocolor      bool
	Nofullpath   bool
	Verbose      bool
}

func (c *Config) fillDefaults() {

	c.Color = true
	c.Fullpath = true
}

// main function is a wrapper on the realMain function and emits OS exit code based on wrapped function
func main() {
	os.Exit(realMain())
}

func realMain() int {

	var config Config
	config.fillDefaults()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var verbose bool
	flag.BoolVar(&verbose, "verbose", config.Verbose, "emit verbose output")
	flag.BoolVar(&verbose, "v", config.Verbose, "emit verbose output")

	var debug bool
	flag.BoolVar(&debug, "debug", config.Debug, "emit debug messages")

	var stdin bool
	flag.BoolVar(&stdin, "stdin", false, "read from ignore file from STDIN")
	flag.BoolVar(&stdin, "s", false, "read from ignore file from STDIN")

	var color bool
	flag.BoolVar(&color, "color", config.Color, "enable colors")
	flag.BoolVar(&color, "c", config.Color, "enable colors")

	var nocolor bool
	flag.BoolVar(&nocolor, "nocolor", config.Nocolor, "disable use of colors")
	flag.BoolVar(&nocolor, "n", config.Nocolor, "disable use of colors")

	var ignorefile string
	flag.StringVar(&ignorefile, "ignorefile", "", "a path to a specific ignore file")
	flag.StringVar(&ignorefile, "i", "", "a path to a specific ignore file")

	var excluded bool
	flag.BoolVar(&excluded, "excluded", config.Excluded, "only output excluded files")
	flag.BoolVar(&excluded, "x", config.Excluded, "only output excluded files")

	var included bool
	flag.BoolVar(&included, "included", config.Included, "only output included files")

	var invertcolors bool
	flag.BoolVar(&invertcolors, "invertcolors", config.Invertcolors, "inverts the used color")

	var fullpath bool
	flag.BoolVar(&fullpath, "fullpath", config.Fullpath, "emits files and directories with full path")
	flag.BoolVar(&fullpath, "f", config.Fullpath, "emits files and directories with full path")

	var nofullpath bool
	flag.BoolVar(&nofullpath, "nofullpath", config.Nofullpath, "emits files and directories without full path")

	// NO_COLOR (https://no-color.org/): disables color when present and
	// non-empty, regardless of its value — unless overridden by a
	// higher-precedence signal (an explicit --color/-c flag, or
	// CLICOLOR_FORCE); see resolveColorConfig for the full precedence.
	origNoColorEnv, origNoColorSet := os.LookupEnv("NO_COLOR")
	noColorEnv := origNoColorEnv
	// Restore NO_COLOR to its original state before returning, since an
	// override below may clear it for the duration of this call only.
	defer func() {
		if origNoColorSet {
			_ = os.Setenv("NO_COLOR", origNoColorEnv)
		} else {
			_ = os.Unsetenv("NO_COLOR")
		}
	}()
	// CLICOLOR / CLICOLOR_FORCE (https://bixense.com/clicolors/):
	// CLICOLOR=0 disables color; CLICOLOR_FORCE!=0 forces color even when not a TTY.
	cliColorEnv := os.Getenv("CLICOLOR")
	cliColorForceEnv := os.Getenv("CLICOLOR_FORCE")

	flag.Parse()

	globalConfigDir := os.Getenv("HOME") + "/.config"
	XDGConfigHome := os.Getenv("XDG_CONFIG_HOME")

	if XDGConfigHome != "" {
		globalConfigDir = XDGConfigHome
	}

	stevedoreGlobalConfigDir := globalConfigDir + "/stevedore/"

	globalConfigFile := stevedoreGlobalConfigDir + "config.json"
	_, err := loadGlobalConfigFile(globalConfigFile, &config)

	if err != nil {
		fmt.Printf("Error attempting to read global configuration file: %s, continuing...\n", err)
	}

	configFile := ".stevedore.json"
	_, err = loadLocalConfigFile(configFile, &config)

	if err != nil {
		fmt.Printf("Error attempting to read configuration file: %s, continuing...\n", err)
	}

	// colorFlagExplicit tracks whether the user explicitly passed --color/-c or
	// --nocolor/-n on the command line, so those take precedence over the
	// NO_COLOR/CLICOLOR/CLICOLOR_FORCE environment variables below.
	colorFlagExplicit := false

	markFlags := func(f *flag.Flag) {
		switch {
		case f.Name == "debug":
			config.Debug = true
		case f.Name == "color" || f.Name == "c":
			config.Color = true
			config.Nocolor = false
			colorFlagExplicit = true
		case f.Name == "nocolor" || f.Name == "n":
			config.Nocolor = true
			config.Color = false
			colorFlagExplicit = true
		case f.Name == "ignorefile" || f.Name == "i":
			config.Ignorefile = ignorefile
		case f.Name == "excluded":
			config.Excluded = true
		case f.Name == "included":
			config.Included = true
		case f.Name == "invertcolors":
			config.Invertcolors = true
		case f.Name == "fullpath" || f.Name == "f":
			config.Fullpath = true
			config.Nofullpath = false
		case f.Name == "nofullpath":
			config.Nofullpath = true
			config.Fullpath = false
		case f.Name == "verbose" || f.Name == "v":
			config.Verbose = true
		default:
			fmt.Printf("No special handling of this command line flag %s.\n", f.Name)
		}
	}
	flag.Visit(markFlags)

	// Precedence for color settings: CLI flags > env vars > config file > defaults.
	// forceColor tracks CLICOLOR_FORCE so fatih/color's own TTY auto-detection
	// can be overridden below, since CLICOLOR_FORCE must apply even when not a TTY.
	var forceColor bool
	config.Color, config.Nocolor, forceColor = resolveColorConfig(
		colorFlagExplicit, config.Color, config.Nocolor,
		noColorEnv, cliColorEnv, cliColorForceEnv,
	)

	path := flag.Arg(0)

	if path == "" {
		path = "."
	}

	var ignoreLines string

	if stdin {
		scanner := bufio.NewScanner(os.Stdin)

		var lines []string
		for scanner.Scan() {
			// read line from stdin using newline as separator
			line := scanner.Text()

			// append the line to a slice
			lines = append(lines, line)
		}
		ignoreLines = strings.Join(lines, "\n")

		if config.Debug {
			fmt.Println("path: ", path)
			fmt.Println("ignore string from STDIN")
			fmt.Printf("ignorelines: \n%s\n\n", ignoreLines)
		}

	} else if config.Ignorefile == "" {
		config.Ignorefile = path + "/.dockerignore"

		if config.Debug {
			fmt.Println("path: ", path)
			fmt.Println("ignoreFile: ", config.Ignorefile)
		}
	}

	if config.Included && config.Excluded {
		config.Included = false
		config.Excluded = false
	} else if !config.Included && !config.Excluded {
		config.Included = true
		config.Excluded = true
	}

	if config.Debug {
		fmt.Println("CLI flags:")
		fmt.Println("\tcolor: ", color)
		fmt.Println("\tnocolor: ", nocolor)
		fmt.Println("\tignorefile: ", ignorefile)
		fmt.Println("\tdebug: ", debug)
		fmt.Println("\tverbose: ", verbose)
		fmt.Println("\texcluded: ", excluded)
		fmt.Println("\tincluded: ", included)
		fmt.Println("\tfullpath: ", fullpath)
		fmt.Println("\tnofullpath: ", nofullpath)
		fmt.Println("\ttail: ", flag.Args())
	}

	if config.Debug {
		fmt.Println("Environment Variables")
		fmt.Println("\tNO_COLOR: ", noColorEnv)
		fmt.Println("\tCLICOLOR: ", cliColorEnv)
		fmt.Println("\tCLICOLOR_FORCE: ", cliColorForceEnv)
	}

	var ignoreObject = ignore.CompileIgnoreLines([]string{}...)

	if stdin {
		ignoreObject = ignore.CompileIgnoreLines(ignoreLines)
	} else {

		ignoreObject, err = ignore.CompileIgnoreFile(config.Ignorefile)

		if err != nil {
			fmt.Printf("unable to read %s file\n", config.Ignorefile)
			return 1
		}
	}

	if config.Color {
		config.Nocolor = false
	}

	if config.Nocolor {
		config.Color = false
		config.Invertcolors = false
	}

	// An explicit --color/-c flag or CLICOLOR_FORCE must win over NO_COLOR and
	// fatih/color's own TTY auto-detection, to honor the documented "CLI flags
	// > env vars" precedence and CLICOLOR_FORCE's "regardless of piping"
	// contract. Otherwise restore fatih/color's originally detected value so
	// repeated invocations within the same process (e.g. tests) don't leak an
	// overridden state across each other.
	// fatih/color also re-checks the live NO_COLOR env var every time it
	// creates a Color (independent of the NoColor package var above), so it
	// must be cleared here too or a forced Set() would still be suppressed;
	// the deferred restore above puts it back before this call returns.
	overrideAutoDetect := forceColor || (colorFlagExplicit && config.Color)
	if overrideAutoDetect {
		colorize.NoColor = false
		_ = os.Unsetenv("NO_COLOR")
	} else {
		colorize.NoColor = defaultNoColor
	}

	if config.Invertcolors {
		ignoredColor, includedColor = includedColor, ignoredColor
	}

	if config.Nofullpath {
		config.Fullpath = false
	}

	ownIgnoreFile := ".stevedoreignore"
	ownIgnoreObject := ignore.CompileIgnoreLines([]string{}...)

	if _, err := os.Stat(ownIgnoreFile); !errors.Is(err, fs.ErrNotExist) {
		if config.Debug {
			fmt.Println("stevedore ignorefile found")
		}

		ownIgnoreObject, err = ignore.CompileIgnoreFile(ownIgnoreFile)

		if err != nil {
			fmt.Printf("unable to read %s file, ignoring - %v\n", ownIgnoreFile, err)
			ownIgnoreObject = ignore.CompileIgnoreLines([]string{}...)
		}

	} else if config.Debug {
		fmt.Println("No stevedore ignorefile found")
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if ownIgnoreObject.MatchesPath(path) && info.IsDir() {
			if config.Verbose {
				fmt.Printf("%s have been ignored by stevedore, no traversal\n", info.Name())
			}
			return filepath.SkipDir
		} else if ownIgnoreObject.MatchesPath(path) {
			if config.Verbose {
				fmt.Printf("%s have been ignored by stevedore\n", info.Name())
			}
			return nil
		}

		var entry string

		if config.Fullpath {
			entry = path
		} else {
			entry = info.Name()
		}

		if path == "." {
			if config.Verbose {
				fmt.Printf("%s is ignored, but traversed by stevedore by default\n", entry)
			}
			return nil
		}

		var aliasedPath = path

		if info.IsDir() {
			aliasedPath += "/"
		}

		if ignoreObject.MatchesPath(path) || ignoreObject.MatchesPath(aliasedPath) {
			if config.Excluded {
				if config.Color {
					colorize.Set(ignoredColor)
				}
				if config.Verbose {
					fmt.Printf("path %s ignored and is not included in Docker image\n", entry)
				} else {
					fmt.Printf("%s\n", entry)
				}
				colorize.Unset()
			}
		} else {
			if config.Included {
				if config.Color {
					colorize.Set(includedColor)
				}
				if config.Verbose {
					fmt.Printf("path %s not ignored and is included in Docker image\n", entry)
				} else {
					fmt.Printf("%s\n", entry)
				}
				colorize.Unset()
			}
		}

		if config.Debug {
			fmt.Printf("visited file or dir: %q\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return 2
	}

	return 0
}

// resolveColorConfig applies the NO_COLOR (https://no-color.org/) and
// CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/) environment
// variables to the color/nocolor settings, honoring precedence:
// CLI flags > env vars > config file > defaults.
//
// color and nocolor are the settings as resolved so far from defaults and
// config files. If colorFlagExplicit is true (the user passed --color/-c or
// --nocolor/-n on the command line), the env vars are ignored entirely and
// color/nocolor are returned unchanged. Otherwise, CLICOLOR_FORCE (if set to
// anything other than "0") takes precedence over NO_COLOR, which takes
// precedence over CLICOLOR=0.
func resolveColorConfig(colorFlagExplicit bool, color, nocolor bool, noColorEnv, cliColorEnv, cliColorForceEnv string) (resolvedColor, resolvedNocolor, forceColor bool) {
	resolvedColor, resolvedNocolor = color, nocolor

	if colorFlagExplicit {
		return
	}

	switch {
	case cliColorForceEnv != "" && cliColorForceEnv != "0":
		resolvedColor, resolvedNocolor, forceColor = true, false, true
	case noColorEnv != "":
		resolvedColor, resolvedNocolor = false, true
	case cliColorEnv == "0":
		resolvedColor, resolvedNocolor = false, true
	}

	return
}

func loadGlobalConfigFile(configFile string, config *Config) (rv bool, err error) {

	if _, err := os.Stat(configFile); !errors.Is(err, fs.ErrNotExist) {

		jsonData, err := os.ReadFile(configFile)

		if err != nil {
			return false, err
		}

		err = json.Unmarshal(jsonData, &config)
		if err != nil {
			return false, err
		}

		if config.Debug {
			fmt.Println("Config file:")
			fmt.Println("\tcolor: ", config.Color)
			fmt.Println("\tnocolor: ", config.Nocolor)
			fmt.Println("\tignorefile: ", config.Ignorefile)
			fmt.Println("\tdebug: ", config.Debug)
			fmt.Println("\tverbose: ", config.Verbose)
			fmt.Println("\texcluded: ", config.Excluded)
			fmt.Println("\tincluded: ", config.Included)
			fmt.Println("\tfullpath: ", config.Fullpath)
			fmt.Println("\tnofullpath: ", config.Nofullpath)
		}
	} else {
		return false, fmt.Errorf("Config file %s not found", configFile)
	}

	return true, nil
}

func loadLocalConfigFile(configFile string, config *Config) (rv bool, err error) {

	if _, err := os.Stat(configFile); !errors.Is(err, fs.ErrNotExist) {

		jsonData, err := os.ReadFile(configFile)

		if err != nil {
			return false, err
		}

		err = json.Unmarshal(jsonData, &config)
		if err != nil {
			return false, err
		}

		if config.Debug {
			fmt.Println("Config file:")
			fmt.Println("\tcolor: ", config.Color)
			fmt.Println("\tnocolor: ", config.Nocolor)
			fmt.Println("\tignorefile: ", config.Ignorefile)
			fmt.Println("\tdebug: ", config.Debug)
			fmt.Println("\tverbose: ", config.Verbose)
			fmt.Println("\texcluded: ", config.Excluded)
			fmt.Println("\tincluded: ", config.Included)
			fmt.Println("\tfullpath: ", config.Fullpath)
			fmt.Println("\tnofullpath: ", config.Nofullpath)
		}
	} else {
		return false, fmt.Errorf("Config file %s not found", configFile)
	}

	return true, nil
}
