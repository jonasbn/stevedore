package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	ignore "github.com/sabhiram/go-gitignore"
)

var (
	ignoredColor  = color.FgGreen
	includedColor = color.FgHiRed
)

type config struct {
	Verbose      bool
	Debug        bool
	Color        bool
	Nocolor      bool
	Ignorefile   string
	Excluded     bool
	Included     bool
	Invertcolors bool
	Fullpath     bool
	Nofullpath   bool
}

// main function is a wrapper on the realMain function and emits OS exit code based on wrapped function
func main() {
	os.Exit(realMain())
}

func realMain() int {

	var config config
	configFile := ".stevedore.json"

	if _, err := os.Stat(configFile); !errors.Is(err, fs.ErrNotExist) {
		if config.Debug {
			fmt.Println("stevedore configuration file found")
		}
		jsonData, err := os.ReadFile(configFile)

		if err != nil {
			fmt.Printf("unable to read %s file, ignoring - %v\n", configFile, err)
		}

		err = json.Unmarshal(jsonData, &config)
		if err != nil {
			fmt.Println("error unmarshalling JSON configuration:", err)
		}

	} else if config.Debug {
		fmt.Println("No stevedore configuration file found")
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.BoolVar(&config.Verbose, "verbose", false, "emit verbose output")
	flag.BoolVar(&config.Verbose, "v", false, "emit verbose output")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "emit debug messages")

	var stdin bool
	flag.BoolVar(&stdin, "stdin", false, "read from ignore file from STDIN")
	flag.BoolVar(&stdin, "s", false, "read from ignore file from STDIN")

	flag.BoolVar(&config.Color, "color", true, "enable colors")
	flag.BoolVar(&config.Color, "c", true, "enable colors")

	flag.BoolVar(&config.Nocolor, "nocolor", false, "disable use of colors")
	flag.BoolVar(&config.Nocolor, "n", false, "disable use of colors")

	flag.StringVar(&config.Ignorefile, "ignorefile", "", "a path to an specific ignore file")
	flag.StringVar(&config.Ignorefile, "i", "", "a path to an specific ignore file")

	flag.BoolVar(&config.Excluded, "excluded", false, "only output excluded files")
	flag.BoolVar(&config.Excluded, "x", false, "only output excluded files")

	flag.BoolVar(&config.Included, "included", false, "only output included files")

	flag.BoolVar(&config.Invertcolors, "invertcolors", false, "inverts the used color")

	flag.BoolVar(&config.Fullpath, "fullpath", true, "emits files and directories with full path")
	flag.BoolVar(&config.Fullpath, "f", true, "emits files and directories with full path")

	flag.BoolVar(&config.Nofullpath, "nofullpath", false, "emits files and directories without full path")

	nocolorEnv := os.Getenv("NO_COLOR")

	flag.Parse()

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

		if debug {
			fmt.Println("path: ", path)
			fmt.Println("ignore string from STDIN")
			fmt.Printf("ignorelines: \n%s\n\n", ignoreLines)
		}

	} else if config.Ignorefile == "" {
		config.Ignorefile = path + "/.dockerignore"

		if debug {
			fmt.Println("path: ", path)
			fmt.Println("ignoreFile: ", config.Ignorefile)
		}
	}

	if config.Excluded {
		config.Included = false
	}

	if config.Included {
		config.Excluded = false
	}

	if !config.Included && !config.Excluded {
		config.Included = true
		config.Excluded = true
	}

	if debug {
		fmt.Println("color: ", config.Color)
		fmt.Println("nocolor: ", config.Nocolor)
		fmt.Println("ignoreFile: ", config.Ignorefile)
		fmt.Println("debug: ", debug)
		fmt.Println("verbose: ", config.Verbose)
		fmt.Println("excluded: ", config.Excluded)
		fmt.Println("included: ", config.Included)
		fmt.Println("fullpath: ", config.Fullpath)
		fmt.Println("nofullpath: ", config.Nofullpath)
		fmt.Println("tail: ", flag.Args())
		fmt.Println("ENV: ", nocolorEnv)
	}

	var ignoreObject = ignore.CompileIgnoreLines([]string{}...)

	if stdin {
		ignoreObject = ignore.CompileIgnoreLines(ignoreLines)
	} else {
		var err error
		ignoreObject, err = ignore.CompileIgnoreFile(config.Ignorefile)

		if err != nil {
			log.Fatalf("unable to read %s file", config.Ignorefile)
			return 1
		}
	}

	if config.Nocolor || nocolorEnv != "" || nocolorEnv == "1" {
		config.Color = false
		config.Invertcolors = false
	}

	if config.Invertcolors {
		ignoredColor, includedColor = includedColor, ignoredColor
	}

	if config.Nofullpath {
		config.Fullpath = false
	}

	var err error

	ownIgnoreFile := ".stevedoreignore"
	ownIgnoreObject := ignore.CompileIgnoreLines([]string{}...)

	if _, err := os.Stat(ownIgnoreFile); !errors.Is(err, fs.ErrNotExist) {
		if debug {
			fmt.Println("stevedore ignorefile found")
		}

		ownIgnoreObject, err = ignore.CompileIgnoreFile(ownIgnoreFile)

		if err != nil {
			fmt.Printf("unable to read %s file, ignoring - %v\n", ownIgnoreFile, err)
			ownIgnoreObject = ignore.CompileIgnoreLines([]string{}...)
		}

	} else if debug {
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

		if ignoreObject.MatchesPath(path) {
			if config.Excluded {
				if config.Color {
					color.Set(ignoredColor)
				}
				if config.Verbose {
					fmt.Printf("path %s ignored and is not included in Docker image\n", entry)
				} else {
					fmt.Printf("%s\n", entry)
				}
				color.Unset()
			}
		} else {
			if config.Included {
				if config.Color {
					color.Set(includedColor)
				}
				if config.Verbose {
					fmt.Printf("path %s not ignored and is included in Docker image\n", entry)
				} else {
					fmt.Printf("%s\n", entry)
				}
				color.Unset()
			}
		}

		if debug {
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
