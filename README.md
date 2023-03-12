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

Precedence for configuration of parameters are:

- Global configuration file
- Local configuration file
- Command line parameters

Use the global configuration file for the configuration you prefer for all you projects and invocations.

Add a local configuration file, where you want to continuously override the global configuration for that particular directory and for all your invocations.

See Configuration section for details on configuration.

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

If you find yourself constantly writing out the same command line parameters, you have several options for for using a configuration file:

1. You can configure per project/repository, by having a file named `.stevedore.json`
2. You can specify a file in `$HOME/.config/stevedore/config.json`

You can in either file specify the setting for all command line arguments, with a JSON key/value structure:

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

Precedence for the configuration files are:

- Global configuration file
- Local configuration file
- Command line parameters

Use the global configuration file for the configuration you prefer for all you projects and invocations.

Add a local configuration file, where you want to continuously override the global configuration for that particular directory and for all your invocations.

See also Configuration for more details.

See Parameters section for details on parameters.

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

## Environment

`stevedore` support locating a configuration file in:

- `$HOME/.config/stevedore`
- Named: `config.json`

The directory can be specified using the environment variable:

`$XDG_CONFIG_HOME`, the default is: `$HOME/.config`. If the environment variable is not set, the default is evaluated.

Do note `stevedore` does not support: `$XDG_CONFIG_DIRS`.

See Configuration section for details on configuration.

## Compatibility

- [Docker ignore][DOCKERIGNORE]: `.dockerignore` (main purpose)
- [Git ignore][GITIGNORE]: `.gitignore`
- [Yak ignore][YAKIGNORE]: `.yakignore`

## Incompatibility

`stevedore` does not support:

- following symbolic links in the traversal of directories, this limitation is imposed by the limitation from the library used for the implementation: [`path/filepath` documentation for `WalkDir` function](https://pkg.go.dev/path/filepath#WalkDir)
- `$XDG_CONFIG_DIRS` which are part of the "XDG Base Directory Specification" are not supported at this time

See Configuration section for details on configuration.

## Resources and References

- [Wikipedia article "stevedore"][WIKIPEDIA]
- [Docker ignore][DOCKERIGNORE]
- [Git ignore][GITIGNORE]
- [Go gitignore][GO-GITIGNORE]
- [Yak ignore][YAKIGNORE]
- [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
- [`path/filepath` documentation for `WalkDir` function](https://pkg.go.dev/path/filepath#WalkDir)
- [Background image by photographer Josh Young](https://unsplash.com/photos/Huv8EWe2Vo8)

## License and Copyright

- stevedore command line utility is copyright by jonasbn under a MIT license
- The background image used on the stevedore [website](https://jonasbn.github.io/stevedore/) is copyright by [Josh Young](https://unsplash.com/@joshalexyoung) and is used under the [Unsplash license](https://unsplash.com/license)

[WIKIPEDIA]: https://en.wikipedia.org/wiki/Stevedore
[GO-GITIGNORE]: https://pkg.go.dev/github.com/sabhiram/go-gitignore
[GITIGNORE]: https://git-scm.com/docs/gitignore
[DOCKERIGNORE]: https://docs.docker.com/engine/reference/builder/#dockerignore-file
[YAKIGNORE]: https://jonasbn.github.io/yak/
