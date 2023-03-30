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
	Stdin        bool
}

func (s *Config) fillDefaults() {

	s.Color = true
	s.Fullpath = true
}

// main function is a wrapper on the realMain function and emits OS exit code based on wrapped function
func main() {
	os.Exit(realMain())
}

func realMain() int {

	var config Config

	var flags Config

	var settings Config
	settings.fillDefaults()

	parseFlags(&flags, &settings)
	parseConfig(&config)

	markFlags := func(f *flag.Flag) {
		switch {
		case f.Name == "debug":
			config.Debug = true
		case f.Name == "color" || f.Name == "c":
			config.Color = true
			config.Nocolor = false
		case f.Name == "nocolor" || f.Name == "n":
			config.Nocolor = true
			config.Color = false
		case f.Name == "ignorefile" || f.Name == "i":
			config.Ignorefile = settings.Ignorefile
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

	/*
		colorEnv := resolveColorEnv()

		if colorEnv {
			config.Nocolor = false
		}

		if !colorEnv {
			config.Color = false
			config.Invertcolors = false
		}
	*/

	resolveSettings(&flags, &config, &settings)

	path := flag.Arg(0)

	if path == "" {
		path = "."
	}

	var ignoreLines string

	if settings.Stdin {
		scanner := bufio.NewScanner(os.Stdin)

		var lines []string
		for scanner.Scan() {
			// read line from stdin using newline as separator
			line := scanner.Text()

			// append the line to a slice
			lines = append(lines, line)
		}
		ignoreLines = strings.Join(lines, "\n")

		if settings.Debug {
			fmt.Println("path: ", path)
			fmt.Println("ignore string from STDIN")
			fmt.Printf("ignorelines: \n%s\n\n", ignoreLines)
		}

	} else if settings.Ignorefile == "" {
		settings.Ignorefile = path + "/.dockerignore"

		if settings.Debug {
			fmt.Println("path: ", path)
			fmt.Println("ignoreFile: ", config.Ignorefile)
		}
	}

	if settings.Included && settings.Excluded {
		settings.Included = false
		settings.Excluded = false
	} else if !settings.Included && !settings.Excluded {
		settings.Included = true
		settings.Excluded = true
	}

	if settings.Debug {
		fmt.Println("CLI flags:")
		fmt.Println("\tcolor: ", flags.Color)
		fmt.Println("\tnocolor: ", flags.Nocolor)
		fmt.Println("\tignorefile: ", flags.Ignorefile)
		fmt.Println("\tdebug: ", flags.Debug)
		fmt.Println("\tverbose: ", flags.Verbose)
		fmt.Println("\texcluded: ", flags.Excluded)
		fmt.Println("\tincluded: ", flags.Included)
		fmt.Println("\tfullpath: ", flags.Fullpath)
		fmt.Println("\tnofullpath: ", flags.Nofullpath)
		fmt.Println("\ttail: ", flag.Args())
	}

	if settings.Debug {
		fmt.Println("Environment Variables")
		fmt.Println("ENV COLOR: ", resolveColorEnv())
	}

	var err error
	var ignoreObject = ignore.CompileIgnoreLines([]string{}...)

	if settings.Stdin {
		ignoreObject = ignore.CompileIgnoreLines(ignoreLines)
	} else {

		ignoreObject, err = ignore.CompileIgnoreFile(settings.Ignorefile)

		if err != nil {
			fmt.Printf("unable to read %s file\n", settings.Ignorefile)
			return 1
		}
	}

	if settings.Color {
		settings.Nocolor = false
	}

	if settings.Nocolor {
		settings.Color = false
		settings.Invertcolors = false
	}

	if settings.Invertcolors {
		ignoredColor, includedColor = includedColor, ignoredColor
	}

	if settings.Nofullpath {
		settings.Fullpath = false
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

	} else if settings.Debug {
		fmt.Println("No stevedore ignorefile found")
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if ownIgnoreObject.MatchesPath(path) && info.IsDir() {
			if settings.Verbose {
				fmt.Printf("%s have been ignored by stevedore, no traversal\n", info.Name())
			}
			return filepath.SkipDir
		} else if ownIgnoreObject.MatchesPath(path) {
			if settings.Verbose {
				fmt.Printf("%s have been ignored by stevedore\n", info.Name())
			}
			return nil
		}

		var entry string

		if settings.Fullpath {
			entry = path
		} else {
			entry = info.Name()
		}

		if path == "." {
			if settings.Verbose {
				fmt.Printf("%s is ignored, but traversed by stevedore by default\n", entry)
			}
			return nil
		}

		if ignoreObject.MatchesPath(path) {
			if settings.Excluded {
				if settings.Color {
					colorize.Set(ignoredColor)
				}
				if settings.Verbose {
					fmt.Printf("path %s ignored and is not included in Docker image\n", entry)
				} else {
					fmt.Printf("%s\n", entry)
				}
				colorize.Unset()
			}
		} else {
			if settings.Included {
				if settings.Color {
					colorize.Set(includedColor)
				}
				if settings.Verbose {
					fmt.Printf("path %s not ignored and is included in Docker image\n", entry)
				} else {
					fmt.Printf("%s\n", entry)
				}
				colorize.Unset()
			}
		}

		if settings.Debug {
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

func loadGlobalConfigFile(configFile string, config *Config) (loaded bool, err error) {

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

func loadLocalConfigFile(configFile string, config *Config) (loaded bool, err error) {

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

func parseFlags(flags *Config, settings *Config) (err error) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.BoolVar(&flags.Verbose, "verbose", settings.Verbose, "emit verbose output")
	flag.BoolVar(&flags.Verbose, "v", settings.Verbose, "emit verbose output")

	flag.BoolVar(&flags.Debug, "debug", settings.Debug, "emit debug messages")

	flag.BoolVar(&flags.Stdin, "stdin", settings.Stdin, "read from ignore file from STDIN")
	flag.BoolVar(&flags.Stdin, "s", settings.Stdin, "read from ignore file from STDIN")

	flag.BoolVar(&flags.Color, "color", settings.Color, "enable colors")
	flag.BoolVar(&flags.Color, "c", settings.Color, "enable colors")

	flag.BoolVar(&flags.Nocolor, "nocolor", settings.Nocolor, "disable use of colors")
	flag.BoolVar(&flags.Nocolor, "n", settings.Nocolor, "disable use of colors")

	flag.StringVar(&flags.Ignorefile, "ignorefile", "", "a path to an specific ignore file")
	flag.StringVar(&flags.Ignorefile, "i", "", "a path to an specific ignore file")

	flag.BoolVar(&flags.Excluded, "excluded", settings.Excluded, "only output excluded files")
	flag.BoolVar(&flags.Excluded, "x", settings.Excluded, "only output excluded files")

	flag.BoolVar(&flags.Included, "included", settings.Included, "only output included files")

	flag.BoolVar(&flags.Invertcolors, "invertcolors", settings.Invertcolors, "inverts the used color")

	flag.BoolVar(&flags.Fullpath, "fullpath", settings.Fullpath, "emits files and directories with full path")
	flag.BoolVar(&flags.Fullpath, "f", settings.Fullpath, "emits files and directories with full path")

	flag.BoolVar(&flags.Nofullpath, "nofullpath", settings.Nofullpath, "emits files and directories without full path")

	flag.Parse()

	return nil
}

func parseConfig(config *Config) (err error) {

	configFile := ".stevedore.json"
	loaded, err := loadLocalConfigFile(configFile, config)

	if loaded {
		return nil
	}

	if err != nil {
		fmt.Printf("Error attempting to read configuration file: %s, continuing...\n", err)
	}

	globalConfigDir := os.Getenv("HOME") + "/.config"
	XDGConfigHome := os.Getenv("XDG_CONFIG_HOME")

	if XDGConfigHome != "" {
		globalConfigDir = XDGConfigHome
	}

	stevedoreGlobalConfigDir := globalConfigDir + "/stevedore/"

	globalConfigFile := stevedoreGlobalConfigDir + "config.json"
	loaded, err = loadGlobalConfigFile(globalConfigFile, config)

	if loaded {
		return nil
	}

	if err != nil {
		fmt.Printf("Error attempting to read global configuration file: %s, continuing...\n", err)
	}

	return nil
}

func resolveSettings(flags *Config, config *Config, settings *Config) (err error) {

	if flags.Verbose {
		settings.Verbose = true
	} else {
		settings.Verbose = config.Verbose
	}

	if flags.Debug {
		settings.Debug = true
	} else {
		settings.Debug = config.Debug
	}

	// we do not support STDIN via config file
	if flags.Stdin {
		settings.Stdin = true
	}

	if flags.Color {
		settings.Color = true
		settings.Nocolor = false
	} else {
		if !resolveColorEnv() {
			settings.Color = config.Color
		} else {
			if config.Color {
				settings.Color = true
			} else {
				settings.Color = false
			}
		}
	}

	if flags.Nocolor {
		settings.Nocolor = true
		settings.Color = false
	} else {
		settings.Nocolor = config.Nocolor
	}

	if flags.Ignorefile != "" {
		settings.Ignorefile = flags.Ignorefile
	} else {
		settings.Ignorefile = config.Ignorefile
	}

	if flags.Excluded {
		settings.Excluded = true
	} else {
		settings.Excluded = config.Excluded
	}

	if flags.Included {
		settings.Included = true
	} else {
		settings.Included = config.Included
	}

	if flags.Invertcolors {
		settings.Invertcolors = true
	} else {
		settings.Invertcolors = config.Invertcolors
	}

	if flags.Fullpath {
		settings.Fullpath = true
	} else {
		settings.Fullpath = config.Fullpath
	}

	if flags.Nofullpath {
		settings.Nofullpath = true
	} else {
		settings.Nofullpath = config.Nofullpath
	}

	return nil
}

func resolveColorEnv() bool {
	switch {
	case nocolor():
		return false
	case clicolor():
		return true
	default:
		return false
	}
}

func nocolor() bool {
	_, ok := os.LookupEnv("NO_COLOR")

	if ok {
		return true
	} else {
		return false
	}
}

func clicolor() bool {
	clicolor := os.Getenv("CLICOLOR")
	clicolorForce := os.Getenv("CLICOLOR_FORCE")

	// REF: https://stackoverflow.com/questions/43947363/detect-if-a-command-is-piped-or-not
	fi, _ := os.Stdin.Stat()

	if fi.Mode()&os.ModeCharDevice == 0 {
		return false
	}

	if clicolor == "0" || clicolor == "" {

		if clicolorForce == "0" || clicolorForce == "" {
			return false
		} else {
			return true
		}

	} else {
		return true
	}
}
