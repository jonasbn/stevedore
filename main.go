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
	ignored  = color.FgGreen
	included = color.FgHiRed
)

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "emit verbose output")
	flag.BoolVar(&verbose, "v", false, "emit verbose output")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "emit debug messages")
	flag.BoolVar(&debug, "d", false, "emit debug messages")

	var colorOutput bool
	flag.BoolVar(&colorOutput, "color", true, "enable colors")
	flag.BoolVar(&colorOutput, "c", true, "enable colors")

	var nocolorOutput bool
	flag.BoolVar(&nocolorOutput, "nocolor", false, "disable use of colors")
	flag.BoolVar(&nocolorOutput, "n", false, "disable use of colors")

	var ignoreFile string
	flag.StringVar(&ignoreFile, "ignorefile", "", "a path to an specific ignore file")
	flag.StringVar(&ignoreFile, "i", "", "a path to an specific ignore file")

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

	if debug {
		fmt.Println("color:", colorOutput)
		fmt.Println("nocolor:", nocolorOutput)
		fmt.Println("ignoreFile:", ignoreFile)
		fmt.Println("verbose:", verbose)
		fmt.Println("tail:", flag.Args())
		fmt.Println("ENV:", nocolorEnv)
	}

	ignoreObject, err := ignore.CompileIgnoreFile(ignoreFile)

	if err != nil {
		log.Fatalf("unable to read .dockerignore file")
		os.Exit(1)
	}

	if nocolorOutput || nocolorEnv != "" || nocolorEnv == "1" {
		colorOutput = false
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if ignoreObject.MatchesPath(info.Name()) {
			if colorOutput {
				color.Set(ignored)
			}
			if verbose {
				fmt.Printf("path %s ignored and is included in Docker image\n", info.Name())
			} else {
				fmt.Printf("%s ignored\n", info.Name())
			}
			color.Unset()
		} else {
			if colorOutput {
				color.Set(included)
			}
			if verbose {
				fmt.Printf("path %s not ignored and is included in Docker image\n", info.Name())
			} else {
				fmt.Printf("%s included\n", info.Name())
			}
			color.Unset()
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
