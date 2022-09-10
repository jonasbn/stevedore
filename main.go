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

	helpPtr := flag.Bool("help", false, "emit help message")
	debugPtr := flag.Bool("debug", false, "emit debug messages")
	colorPtr := flag.Bool("color", true, "enable colors")
	nocolorPtr := flag.Bool("nocolor", false, "disable use of colors")
	verbosePtr := flag.Bool("verbose", false, "emit verbose output")
	verbosePtr = flag.Bool("v", false, "emit verbose output")

	ignoreFilePtr := flag.String("ignorefile", "", "a path to an specific ignore file")

	nocolorEnv := os.Getenv("NO_COLOR")

	flag.Parse()

	path := flag.Arg(0)

	if *ignoreFilePtr == "" {
		*ignoreFilePtr = path + "/.dockerignore"

		if *debugPtr {
			fmt.Println("path:", path)
			fmt.Println("ignoreFilePtr:", *ignoreFilePtr)
		}
	}

	if *debugPtr {
		fmt.Println("helpPtr:", *helpPtr)
		fmt.Println("colorPtr:", *colorPtr)
		fmt.Println("nocolorPtr:", *nocolorPtr)
		fmt.Println("ignoreFilePtr:", *ignoreFilePtr)
		fmt.Println("verbosePtr:", *verbosePtr)
		fmt.Println("tail:", flag.Args())
		fmt.Println("ENV:", nocolorEnv)
	}

	ignoreObject, err := ignore.CompileIgnoreFile(*ignoreFilePtr)

	if err != nil {
		log.Fatalf("unable to read .dockerignore file")
		os.Exit(1)
	}

	if *nocolorPtr || nocolorEnv != "" || nocolorEnv == "1" {
		*colorPtr = false
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if ignoreObject.MatchesPath(info.Name()) {
			if *colorPtr {
				color.Set(ignored)
			}
			if *verbosePtr {
				fmt.Printf("path %s ignored and is included in Docker image\n", info.Name())
			} else {
				fmt.Printf("%s ignored\n", info.Name())
			}
			color.Unset()
		} else {
			if *colorPtr {
				color.Set(included)
			}
			if *verbosePtr {
				fmt.Printf("path %s not ignored and is included in Docker image\n", info.Name())
			} else {
				fmt.Printf("%s included\n", info.Name())
			}
			color.Unset()
		}
		if *debugPtr {
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
