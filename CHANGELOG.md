# Change log for stevedore

## 0.6.1 2022-11-22 Bug fix release, update not required

- Fixed bug in consumption of file from STDIN, the contents would only be partial

## 0.6.0 2022-11-21 Feature release, update not required

- Added support for reading ignore file via STDIN, using new parameter:
  - `--stdin`

## 0.5.0 2022-11-15 Feature release, update not required

- Implementation of basic support of `.stevedorignore` file

## 0.4.0 2022-09-17 Feature release, update not required

- Added new command line option `--invertcolor` which invert the used colors

## 0.3.1 2022-09-17 Minor bug fix release, update recommended

- Had shuffled the logic around, confusing the terms when updated the main body of code

## 0.3.0 2022-09-17 Feature release, update not required

- Added two new command line options
  - `--excluded` which only outputs what is excluded
  - `--included` which only outputs what is included

  The first need more work and I need to settle on a useful path output, so this is WIP

## 0.2.0 2022-09-11 Feature release, update not required

- Got the short forms for the command line in place

## 0.1.0 2022-09-10 Feature release, update not required

- First working version, supporting:
  - reporting on ignored and non-ignored file system components handled by a given Docker ignore file
  - Specification of path to analyze
  - Specification of alternative Docker ignore file
  - Colorized output
  - Basic support for `NO_COLOR` environment variable
  - Verbose output
  - Debug output
