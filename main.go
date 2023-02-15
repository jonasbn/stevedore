package main

import (
	"bufio"
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

// main function is a wrapper on the realMain function and emits OS exit code based on wrapped function
func main() {
	os.Exit(realMain())
}

func realMain() int {

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "emit verbose output")
	flag.BoolVar(&verbose, "v", false, "emit verbose output")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "emit debug messages")

	var stdin bool
	flag.BoolVar(&stdin, "stdin", false, "read from ignore file from STDIN")
	flag.BoolVar(&stdin, "s", false, "read from ignore file from STDIN")

	var colorOutput bool
	flag.BoolVar(&colorOutput, "color", true, "enable colors")
	flag.BoolVar(&colorOutput, "c", true, "enable colors")

	var nocolorOutput bool
	flag.BoolVar(&nocolorOutput, "nocolor", false, "disable use of colors")
	flag.BoolVar(&nocolorOutput, "n", false, "disable use of colors")

	var ignoreFile string
	flag.StringVar(&ignoreFile, "ignorefile", "", "a path to an specific ignore file")
	flag.StringVar(&ignoreFile, "i", "", "a path to an specific ignore file")

	var excluded bool
	flag.BoolVar(&excluded, "excluded", false, "only output excluded files")
	flag.BoolVar(&excluded, "x", false, "only output excluded files")

	var included bool
	flag.BoolVar(&included, "included", false, "only output included files")

	var colorOutputInverted bool
	flag.BoolVar(&colorOutputInverted, "invertcolors", false, "inverts the used color")

	var fullPath bool
	flag.BoolVar(&fullPath, "fullpath", true, "emits files and directories with full path")
	flag.BoolVar(&fullPath, "f", true, "emits files and directories with full path")

	var noFullPath bool
	flag.BoolVar(&noFullPath, "nofullpath", false, "emits files and directories without full path")

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

	} else if ignoreFile == "" {
		ignoreFile = path + "/.dockerignore"

		if debug {
			fmt.Println("path: ", path)
			fmt.Println("ignoreFile: ", ignoreFile)
		}
	}

	if excluded {
		included = false
	}

	if included {
		excluded = false
	}

	if !included && !excluded {
		included = true
		excluded = true
	}

	if debug {
		fmt.Println("color: ", colorOutput)
		fmt.Println("nocolor: ", nocolorOutput)
		fmt.Println("ignoreFile: ", ignoreFile)
		fmt.Println("debug: ", debug)
		fmt.Println("verbose: ", verbose)
		fmt.Println("excluded: ", excluded)
		fmt.Println("included: ", included)
		fmt.Println("fullpath: ", fullPath)
		fmt.Println("tail: ", flag.Args())
		fmt.Println("ENV: ", nocolorEnv)
	}

	var ignoreObject = ignore.CompileIgnoreLines([]string{}...)

	if stdin {
		ignoreObject = ignore.CompileIgnoreLines(ignoreLines)
	} else {
		var err error
		ignoreObject, err = ignore.CompileIgnoreFile(ignoreFile)

		if err != nil {
			log.Fatalf("unable to read %s file", ignoreFile)
			return 1
		}
	}

	if nocolorOutput || nocolorEnv != "" || nocolorEnv == "1" {
		colorOutput = false
		colorOutputInverted = false
	}

	if colorOutputInverted {
		ignoredColor, includedColor = includedColor, ignoredColor
	}

	if noFullPath {
		fullPath = false
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
			if verbose {
				fmt.Printf("%s have been ignored by stevedore, no traversal\n", info.Name())
			}
			return filepath.SkipDir
		} else if ownIgnoreObject.MatchesPath(path) {
			if verbose {
				fmt.Printf("%s have been ignored by stevedore\n", info.Name())
			}
			return nil
		}

		var entry string

		if fullPath {
			entry = path
		} else {
			entry = info.Name()
		}

		if path == "." {
			if verbose {
				fmt.Printf("%s is ignored, but traversed by stevedore by default\n", entry)
			}
			return nil
		}

		if ignoreObject.MatchesPath(path) {
			if excluded {
				if colorOutput {
					color.Set(ignoredColor)
				}
				if verbose {
					fmt.Printf("path %s ignored and is not included in Docker image\n", entry)
				} else {
					fmt.Printf("%s\n", entry)
				}
				color.Unset()
			}
		} else {
			if included {
				if colorOutput {
					color.Set(includedColor)
				}
				if verbose {
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
