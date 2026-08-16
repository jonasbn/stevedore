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

func TestResolveColorConfig(t *testing.T) {
	cases := []struct {
		Name              string
		ColorFlagExplicit bool
		Color             bool
		Nocolor           bool
		NoColorEnv        string
		CliColorEnv       string
		CliColorForceEnv  string
		ExpectedColor     bool
		ExpectedNocolor   bool
		ExpectedForce     bool
	}{
		{
			Name:            "no env vars set, defaults preserved",
			Color:           true,
			ExpectedColor:   true,
			ExpectedNocolor: false,
		},
		{
			Name:            "NO_COLOR set disables color",
			Color:           true,
			NoColorEnv:      "1",
			ExpectedColor:   false,
			ExpectedNocolor: true,
		},
		{
			Name:            "NO_COLOR set to empty string does not disable color",
			Color:           true,
			NoColorEnv:      "",
			ExpectedColor:   true,
			ExpectedNocolor: false,
		},
		{
			Name:            "NO_COLOR set to any non-empty value disables color regardless of value",
			Color:           true,
			NoColorEnv:      "0",
			ExpectedColor:   false,
			ExpectedNocolor: true,
		},
		{
			Name:            "CLICOLOR=0 disables color",
			Color:           true,
			CliColorEnv:     "0",
			ExpectedColor:   false,
			ExpectedNocolor: true,
		},
		{
			Name:            "CLICOLOR set to non-zero value has no effect",
			Color:           true,
			CliColorEnv:     "1",
			ExpectedColor:   true,
			ExpectedNocolor: false,
		},
		{
			Name:             "CLICOLOR_FORCE forces color on and requests TTY override",
			Color:            false,
			Nocolor:          true,
			CliColorForceEnv: "1",
			ExpectedColor:    true,
			ExpectedNocolor:  false,
			ExpectedForce:    true,
		},
		{
			Name:             "CLICOLOR_FORCE=0 has no effect",
			Color:            true,
			CliColorForceEnv: "0",
			ExpectedColor:    true,
			ExpectedNocolor:  false,
		},
		{
			Name:             "CLICOLOR_FORCE takes precedence over NO_COLOR",
			Color:            true,
			NoColorEnv:       "1",
			CliColorForceEnv: "1",
			ExpectedColor:    true,
			ExpectedNocolor:  false,
			ExpectedForce:    true,
		},
		{
			Name:              "explicit CLI flag ignores NO_COLOR entirely",
			ColorFlagExplicit: true,
			Color:             true,
			Nocolor:           false,
			NoColorEnv:        "1",
			ExpectedColor:     true,
			ExpectedNocolor:   false,
		},
		{
			Name:              "explicit CLI flag ignores CLICOLOR_FORCE entirely",
			ColorFlagExplicit: true,
			Color:             false,
			Nocolor:           true,
			CliColorForceEnv:  "1",
			ExpectedColor:     false,
			ExpectedNocolor:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			gotColor, gotNocolor, gotForce := resolveColorConfig(
				tc.ColorFlagExplicit, tc.Color, tc.Nocolor,
				tc.NoColorEnv, tc.CliColorEnv, tc.CliColorForceEnv,
			)

			if gotColor != tc.ExpectedColor {
				t.Errorf("color: expected %v, got %v", tc.ExpectedColor, gotColor)
			}
			if gotNocolor != tc.ExpectedNocolor {
				t.Errorf("nocolor: expected %v, got %v", tc.ExpectedNocolor, gotNocolor)
			}
			if gotForce != tc.ExpectedForce {
				t.Errorf("forceColor: expected %v, got %v", tc.ExpectedForce, gotForce)
			}
		})
	}
}

func TestColorEnvArguments(t *testing.T) {
	// We manipulate the Args to set them up for the testcases
	// After this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cases := []struct {
		Name             string
		Args             []string
		NoColorEnv       string
		CliColorEnv      string
		CliColorForceEnv string
		ExpectedExit     int
	}{
		{"NO_COLOR set", []string{"."}, "1", "", "", 0},
		{"CLICOLOR=0 set", []string{"."}, "", "0", "", 0},
		{"CLICOLOR_FORCE=1 set", []string{"."}, "", "", "1", 0},
		{"NO_COLOR set with explicit --color flag", []string{"--color", "."}, "1", "", "", 0},
	}

	for _, tc := range cases {
		t.Setenv("NO_COLOR", tc.NoColorEnv)
		t.Setenv("CLICOLOR", tc.CliColorEnv)
		t.Setenv("CLICOLOR_FORCE", tc.CliColorForceEnv)

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
