# stevedore

`stevedore` is a small command line tool taking it's name from the worker working on the dock with loading cargo unto ships.

REF: [Wikipedia][WIKIPEDIA]

The tool reads a given directory and [Docker ignore file (`.dockerignore`)][DOCKERIGNORE] and outputs a report on what is to be included in a Docker image and what will be ignored.

Like so:

```bash
stevedore .
```

The above example

1. Locates the `.dockerignore` file
1. Reads the current directory (specified as `.`) recursively
1. Compares the located `.dockerignore` file with the contents of the specified directory
1. Outputs a report

```text
. included
.dockerignore included
.gitignore included
README.md included
TODO included
go.mod included
go.sum included
main.go ignored
stevedore included
```

## Parameters

Since this is just an analysis/reporting tool it can be fed with parameters to diverge from the default behaviour.

- `--help` / `-h` emits a brief help message
- `--ignorefile <path>` / `-i` points to alternative ignore file
- `--color` / `-c` emits output on color
- `--nocolor` / `-n` emits output suppressing use of colors
- `--verbose` / `-v` emits more verbose output
- `--debug` emits debug information
- `--included` emits only included files (non-ignored)
- `--excluded` emits only excluded files (ignored)
- `--invertcolors` inverts the used colors

### Verbose

If the verbose flag is set the output is altered and is more explanatory:

```bash
stevedore -verbose .
```

```text
path . not ignored and is included in Docker image
path .dockerignore not ignored and is included in Docker image
path .gitignore not ignored and is included in Docker image
path README.md not ignored and is included in Docker image
path TODO not ignored and is included in Docker image
path go.mod not ignored and is included in Docker image
path go.sum not ignored and is included in Docker image
path main.go ignored and is included in Docker image
path stevedore not ignored and is included in Docker image
```

## Return Values

- `0` indicates a successful run
- `1` ignore file was not found or could not be read
- `2` specified directory could not be read or only partially read

## The stevedore ignore

You can add an ignore file, named `.stevedoreignore` to your directory. It will tell `stevedore` what files and directories to ignore prior to making it's analysis.

Meaning that patterns in this files matched, will be excluded.

The `.stevedoreignore` file follows the general implementation pattern. and example could be:

```gitignore
.git
```

## Compatibility

- [Docker ignore][DOCKERIGNORE]: `.dockerignore` (main purpose)
- [Git ignore][GITIGNORE]: `.gitignore`
- Yak ignore: `.yakignore`

## Incompatibility

- `stevedore` does not support following symbolic links in the traversal of directories

## Resources and References

- [Wikipedia article "stevedore"][WIKIPEDIA]
- [Docker ignore][DOCKERIGNORE]
- [Git ignore][GITIGNORE]
- [Go gitignore][GO-GITIGNORE]

[WIKIPEDIA]: https://en.wikipedia.org/wiki/Stevedore
[GO-GITIGNORE]: https://pkg.go.dev/github.com/sabhiram/go-gitignore
[GITIGNORE]: https://git-scm.com/docs/gitignore
[DOCKERIGNORE]: https://docs.docker.com/engine/reference/builder/#dockerignore-file
