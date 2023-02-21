package main

import (
	"os"
	"testing"
)

func TestArguments(T *testing.T) {
	// We manipuate the Args to set them up for the testcases
	// After this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cases := []struct {
		Name         string
		Args         []string
		ExpectedExit int
	}{
		{"no arguments", []string{""}, 0},
		{"single path argument", []string{"."}, 0},
		{"color argument", []string{"--color"}, 0},
		{"debug argument", []string{"--debug"}, 0},
		{"excluded argument", []string{"--excluded"}, 0},
		{"fullpath argument", []string{"--fullpath"}, 0},
		{"included argument", []string{"--included"}, 0},
		{"invertcolors argument", []string{"--invertcolors"}, 0},
		{"nocolor argument", []string{"--nocolor"}, 0},
		{"nofillpath argument", []string{"--nofullpath"}, 0},
		{"verbose argument", []string{"--verbose"}, 0},
	}

	for _, tc := range cases {
		// we need a value to set Args[0] to cause flag begins parsing at Args[1]
		os.Args = append([]string{tc.Name}, tc.Args...)
		actualExit := realMain()
		if tc.ExpectedExit != actualExit {
			T.Errorf("Wrong exit code for args: %v, expected: %v, got: %v",
				tc.Args, tc.ExpectedExit, actualExit)
		}
	}
}
