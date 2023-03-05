# stevedore

`stevedore` is a small command line tool taking it's name from the worker, working on the dock with loading cargo unto ships.

REF: [Wikipedia][WIKIPEDIA]

The tool reads a given directory and [Docker ignore file (`.dockerignore`)][DOCKERIGNORE] and outputs a report on what is to be included in a Docker image and what will be ignored.

Like so:

```bash
stevedore .
```

The above example

1. Locates the `.dockerignore` file
2. Reads the current directory (specified as `.`) recursively
3. Compares the located `.dockerignore` file with the contents of the specified directory
4. Outputs a report

```text
.
.dockerignore
.gitignore
README.md
TODO
go.mod
go.sum
main.go
stevedore
```

You can actually emit the path parameter, since `stevedore` defaults to current directory.

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
- `--stdin` / `-s` reads ignore file from STDIN
- `--fullpath` / `-f` emits full path of encountered files and directories

### Verbosity

If the verbose flag is set the output is altered and is more explanatory:

```bash
stevedore -verbose .
```

```text
path . is not ignored and is included in Docker image
path .dockerignore is not ignored and is included in Docker image
path .gitignore is not ignored and is included in Docker image
path README.md is not ignored and is included in Docker image
path TODO is not ignored and is included in Docker image
path go.mod is not ignored and is included in Docker image
path go.sum is not ignored and is included in Docker image
path main.go is ignored and is included in Docker image
path stevedore is not ignored and is included in Docker image
```

### Passing in a ignore file using either stdin and ignore file parameters

If you have a ignore file and you want to pass it to `stevedore` you can either use, the `--ignorefile parameter`:

`stevedore --ignorefile /path/to/my/ignorefile`

Or you can pass it in via STDIN:

`cat /path/to/my/ignorefile | stevedore --stdin`

These will render the same result.

## Configuration

If you find yourself constantly writing out the same command line parameters, you can create a local configuration file: `.stevedore.json`

You can specify the setting for all command line arguments, but with a JSON key/value structure:

```json
{
    "$schema": "stevedore-config.schema.json",
    "color": true,
    "debug": false,
    "excluded": false,
    "fullpath": true,
    "ignorefile": ".stevedoreignore",
    "included": false,
    "invertcolor": false,
    "verbose": false
}
```

Parameters not available for configuration:

- `--help`
- `--stdin`

## Return Values

- `0` indicates a successful run
- `1` ignore file was not found or could not be read
- `2` specified directory could not be read or only partially read

## The stevedore ignore

You can add an ignore file, named `.stevedoreignore` to your directory. It will tell `stevedore` what files and directories to ignore prior to making it's analysis.

Meaning that patterns in this files matched, will be _excluded_.

The `.stevedoreignore` file follows the general implementation pattern. and example could be:

```gitignore
.git
```

## Compatibility

- [Docker ignore][DOCKERIGNORE]: `.dockerignore` (main purpose)
- [Git ignore][GITIGNORE]: `.gitignore`
- [Yak ignore][YAKIGNORE]: `.yakignore`

## Incompatibility

- `stevedore` does not support following symbolic links in the traversal of directories

## Resources and References

- [Wikipedia article "stevedore"][WIKIPEDIA]
- [Docker ignore][DOCKERIGNORE]
- [Git ignore][GITIGNORE]
- [Go gitignore][GO-GITIGNORE]
- [Yak ignore][YAKIGNORE]

[WIKIPEDIA]: https://en.wikipedia.org/wiki/Stevedore
[GO-GITIGNORE]: https://pkg.go.dev/github.com/sabhiram/go-gitignore
[GITIGNORE]: https://git-scm.com/docs/gitignore
[DOCKERIGNORE]: https://docs.docker.com/engine/reference/builder/#dockerignore-file
[YAKIGNORE]: https://jonasbn.github.io/yak/
