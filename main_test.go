package main

import (
	"os"
	"testing"
)

func TestArguments(t *testing.T) {
	// We manipulate the Args to set them up for the testcases
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
			t.Errorf("Wrong exit code for args: %v, expected: %v, got: %v",
				tc.Args, tc.ExpectedExit, actualExit)
		}
	}
}

func TestFails(t *testing.T) {
	// We manipulate the Args to set them up for the testcases
	// After this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cases := []struct {
		Name         string
		Args         []string
		ExpectedExit int
	}{
		{"unreadable ignorefile", []string{"--verbose", "tests/unable_to_read_dockerignore"}, 1},
	}

	err := os.MkdirAll("tests/unable_to_read_dockerignore", 0755)
	check(err)

	createEmptyTestFile := func(name string) {
		d := []byte("")
		check(os.WriteFile(name, d, 0333))
	}

	createEmptyTestFile("tests/unable_to_read_dockerignore/.dockerignore")

	defer os.RemoveAll("tests")

	for _, tc := range cases {
		// we need a value to set Args[0] to cause flag begins parsing at Args[1]
		os.Args = append([]string{tc.Name}, tc.Args...)
		actualExit := realMain()
		if tc.ExpectedExit != actualExit {
			t.Errorf("Wrong exit code for args: %v, expected: %v, got: %v",
				tc.Args, tc.ExpectedExit, actualExit)
		}
	}
}

func TestConfig(t *testing.T) {
	// We manipulate the Args to set them up for the testcases
	// After this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cases := []struct {
		Name         string
		Args         []string
		ExpectedExit int
	}{
		{"basic config", []string{"--verbose", "tests/ok"}, 0},
	}

	err := os.MkdirAll("tests/ok", 0755)
	check(err)

	createEmptyTestFile := func(name string) {
		d := []byte("")
		check(os.WriteFile(name, d, 0644))
	}

	createEmptyTestFile("tests/ok/.dockerignore")

	defer os.RemoveAll("tests")

	for _, tc := range cases {
		// we need a value to set Args[0] to cause flag begins parsing at Args[1]
		os.Args = append([]string{tc.Name}, tc.Args...)
		actualExit := realMain()
		if tc.ExpectedExit != actualExit {
			t.Errorf("Wrong exit code for args: %v, expected: %v, got: %v",
				tc.Args, tc.ExpectedExit, actualExit)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
