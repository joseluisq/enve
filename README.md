# Enve ![devel](https://github.com/joseluisq/enve/workflows/devel/badge.svg) [![codecov](https://codecov.io/gh/joseluisq/enve/graph/badge.svg?token=U77DXS42C6)](https://codecov.io/gh/joseluisq/enve) [![Go Report Card](https://goreportcard.com/badge/github.com/joseluisq/enve)](https://goreportcard.com/report/github.com/joseluisq/enve)

> Run a program in a modified environment providing an optional `.env` file or variables from `stdin`.

**Enve** is a cross-platform tool that can load environment variables from a [`.env` file](https://www.ibm.com/docs/en/aix/7.2?topic=files-env-file) or from standard input (stdin) and run a command with those variables set in the environment.

It also allows you to output environment variables in `text`, `json` or `xml` format as well as to overwrite existing ones with values from a custom `.env` file or `stdin`.

Enve can be considered as a counterpart of [GNU env](https://www.gnu.org/software/coreutils/manual/html_node/env-invocation.html) command.

## Install

- **Platforms supported:** `linux`, `darwin`, `windows`, `freebsd`, `openbsd`
- **Architectures supported:** `amd64`, `386`, `arm`, `arm64`, `ppc64le`

```sh
curl -sSL \
   "https://github.com/joseluisq/enve/releases/download/v1.5.1/enve_v1.5.1_linux_amd64.tar.gz" \
   | sudo tar zxf - -C /usr/local/bin/ enve
```

Using Go:

```sh
go install github.com/joseluisq/enve@latest
```

Pre-compiled binaries also available on [joseluisq/enve/releases](https://github.com/joseluisq/enve/releases)

## Usage

By default, **enve** will print all environment variables like `env` command. 

```sh
enve
# Or its equivalent
enve --output text
```

### Executing commands

By default, an optional `.env` file can be loaded from the current working directory.

```sh
enve test.sh
```

## Options

#### `-f, --file`

Loads environment variables from a specific file path.
By default, `enve` will look for a file named `.env` in the current directory.

```sh
# Use a .env file (default)
enve test.sh
# Or specify a custom one
enve --file dev.env test.sh
```

#### `-o, --output`

Outputs all environment variables in a specified format.

```sh
# Print environment variables
enve -o text
enve -o json
enve -o xml

# Or export them to a file
enve -o text > config.txt
enve -o xml > config.xml
enve -o json > config.json
```

#### `-w, --overwrite`

Overwrites existing environment variables with values from the `.env` file or stdin.

```sh
# Overwrite via .env
export API_URL="http://localhost:3000"
enve --overwrite -f .env ./tests.sh

# Or via stdin (which ignores .env file if present)
echo -e "API_URL=http://127.0.0.1:4000" | enve --stdin -w -o text
```

#### `-c, --chdir`

Changes the current working directory before executing the command.

```sh
# Change working directory of a script
enve --chdir /opt/my-app ./test.sh
```

#### `-n, --new-environment`

Starts a new environment containing only variables from either a `.env` file or stdin.

```sh
# Isolated the environment using only variables from .env
enve --new-environment -f devel.env ./test.sh

# Isolate the environment using only variables from stdin
echo -e "APP_HOST=localhost\nAPP_PORT=8080" | enve --stdin -n test.sh
```

#### `-s, --stdin`

Reads environment variables from the standard input (stdin) instead of a file.
When using `--stdin`, the `.env` file is ignored.

```sh
# Pipe environment variables from stdin and run a script
cat development.env | enve --stdin ./my_script.sh
echo -e "APP_HOST=127.0.0.1" | enve -s test.sh
```

#### `-i, --ignore-environment`

Starts with an empty environment skipping any existing environment variables.

```sh
# Run a script in a clean environment
enve --ignore-environment my_script.sh

echo -e "APP_HOST=127.0.0.1" | enve -i --stdin -o json
# {"environment":[]}
```

#### `-z, --no-file`

Prevents `enve` from loading any `.env` file, printing or running a command only with the existing environment.

```sh
# Run a command without loading the default .env file
enve --no-file my_app

# Behaves like the standard 'env' command, printing the current environment
enve -z
```

#### `-h, --help`

```
Run a program in a modified environment providing an optional .env file or variables from stdin

USAGE:
   enve [OPTIONS] COMMAND

OPTIONS:
   -f --file                 Load environment variables from a file path (optional) [default: .env]
   -o --output               Output environment variables using text, json or xml format [default: text]
   -w --overwrite            Overwrite environment variables if already set [default: false]
   -c --chdir                Change currrent working directory
   -n --new-environment      Start a new environment with only variables from the .env file or stdin [default: false]
   -i --ignore-environment   Starts with an empty environment, ignoring any existing environment variables [default: false]
   -z --no-file              Do not load a .env file [default: false]
   -s --stdin                Read only environment variables from stdin and ignore the .env file [default: false]
   -h --help                 Prints help information
   -v --version              Prints version information
```

## Contributions

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in current work by you, as defined in the Apache-2.0 license, shall be dual licensed as described below, without any additional terms or conditions.

Feel free to send some [Pull request](https://github.com/joseluisq/enve/pulls) or file an [issue](https://github.com/joseluisq/enve/issues).

## License

This work is primarily distributed under the terms of both the [MIT license](LICENSE-MIT) and the [Apache License (Version 2.0)](LICENSE-APACHE).

© 2020-present [Jose Quintana](https://joseluisq.net)
