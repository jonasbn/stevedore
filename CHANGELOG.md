# Change log for stevedore

## 2022-09-17 Feature release, update not required

- Added two new command line options
  - `--excluded` which only outputs what is excluded
  - `--included` which only outputs what is included

  The first need more work and I need to settle on a useful path output, so this is WIP

## 2022-09-11 Feature release, update not required

- Got the short forms for the command line in place

## 2022-09-10 Feature release, update not required

- First working version, supporting:
  - reporting on ignored and non-ignored file system components handled by a given Docker ignore file
  - Specification of path to analyze
  - Specification of alternative Docker ignore file
  - Colorized output
  - Basic support for `NO_COLOR` environment variable
  - Verbose output
  - Debug output
