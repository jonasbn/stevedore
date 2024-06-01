# Change log for stevedore

## 0.13.0 2024-06-01 maintenance release, update not required

- Merged pull request [#56](https://github.com/jonasbn/stevedore/pull/56) from @dependabot bumping fatih/color to version 1.17.0

## 0.12.0 2024-02-28 feature release, update recommended

- Merged pull request [#53](https://github.com/jonasbn/stevedore/pull/53) from @jonasbn. Since the path is compared as strings
  The usage of appending a `/` to the path as done in many ignore files.
  
  I had a look at the repository:
  
  - [github/gitignore](https://github.com/github/gitignore)
  
  `stevedore` now handles cases by checking directories both with and without slash appended to the directory name in the ignore file.

## 0.11.0 2023-11-13 maintenance release, update not required

- Merged pull request [#47](https://github.com/jonasbn/stevedore/pull/47) from @dependabot bumping fatih/color to version 1.16.0

- Merged pull request [#31](https://github.com/jonasbn/stevedore/pull/31) from @dependabot bumping fatih/color to version 1.15.0

## 0.10.0 2023-03-12 feature release, update not required

- Added support for global configuration file, now both: local and global configuration files are supported

## 0.9.0 2023-03-05 feature release, update not required

- Fix a _bug_, which was disturbing the tests
- Implemented proper precedence handling between CLI flags and configuration

## 0.8.0 2023-02-17 feature release, update not required

- `.` was included in the ignore pattern output, which does not really make sense, so it has been eliminated from the output. It is not super solution since it is no so flexible, but it will have to do for now, ref: [#19](https://github.com/jonasbn/stevedore/issues/19)

## 0.7.0 2023-02-07 feature release, update not required

- Improved output so both full path and flattened list is available
  - full path is the new default, can be explicitly requested using  `--fullpath`
  - flat structure, can be explicitly requested using `--nofullpath`

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
