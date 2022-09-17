package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	ignore "github.com/sabhiram/go-gitignore"
)

var (
	ignoredColor  = color.FgGreen
	includedColor = color.FgHiRed
)

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "emit verbose output")
	flag.BoolVar(&verbose, "v", false, "emit verbose output")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "emit debug messages")

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

	nocolorEnv := os.Getenv("NO_COLOR")

	flag.Parse()

	path := flag.Arg(0)

	if ignoreFile == "" {
		ignoreFile = path + "/.dockerignore"

		if debug {
			fmt.Println("path:", path)
			fmt.Println("ignoreFile:", ignoreFile)
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
		fmt.Println("tail: ", flag.Args())
		fmt.Println("ENV: ", nocolorEnv)
	}

	ignoreObject, err := ignore.CompileIgnoreFile(ignoreFile)

	if err != nil {
		log.Fatalf("unable to read .dockerignore file")
		os.Exit(1)
	}

	if nocolorOutput || nocolorEnv != "" || nocolorEnv == "1" {
		colorOutput = false
		colorOutputInverted = false
	}

	if colorOutputInverted {
		tmpColor := ignoredColor
		ignoredColor = includedColor
		includedColor = tmpColor
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if ignoreObject.MatchesPath(info.Name()) {
			if excluded {
				if colorOutput {
					color.Set(ignoredColor)
				}
				if verbose {
					fmt.Printf("path %s ignored and is not included in Docker image\n", info.Name())
				} else {
					fmt.Printf("%s\n", info.Name())
				}
				color.Unset()
			}
		} else {
			if included {
				if colorOutput {
					color.Set(includedColor)
				}
				if verbose {
					fmt.Printf("path %s not ignored and is included in Docker image\n", info.Name())
				} else {
					fmt.Printf("%s\n", info.Name())
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
		os.Exit(2)
	}

	os.Exit(0)
}
