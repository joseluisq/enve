# enve [![Build Status](https://travis-ci.com/joseluisq/enve.svg?branch=master)](https://travis-ci.com/joseluisq/enve) [![PkgGoDev](https://pkg.go.dev/badge/github.com/joseluisq/enve)](https://pkg.go.dev/github.com/joseluisq/enve)

> Run a program in a modified environment providing a .env file.

**enve** loads all environment variables from a `.env` file and executes a command. It also has the ability to print environment variables in `text`, `json` or `xml` formats.

It can be considered as a counterpart of [GNU env](https://www.gnu.org/software/coreutils/manual/html_node/env-invocation.html) command.

## Install

- **Platforms supported:** `linux`, `darwin`, `freebsd`, `openbsd`
- **Architectures supported:** `amd64`, `386`, `arm`, `arm64`, `ppc64le`

```sh
curl -sSL \
   "https://github.com/joseluisq/enve/releases/download/v1.1.0/enve_v1.1.0_linux_amd64.tar.gz" \
| sudo tar zxf - -C /usr/local/bin/ enve
```

Using Go:

```sh
go get -u github.com/joseluisq/enve
```

Release binaries also available on [joseluisq/enve/releases](https://github.com/joseluisq/enve/releases)

## Usage

By default **enve** will print all environment variables like `env` command. 

```sh
enve
# Or its equivalent
enve --output text
```

### Executing commands

By default, an optional `.env` file can be loaded from current working directory.

```sh
enve test.sh
```

However it's possible to specify a custom `.env` file using the `--file` or `-f` flags.

```sh
enve --file dev.env test.sh
```

### Printing environment variables

**enve** supports `text`, `json` and `xml` formats.

```sh
enve --output text # or just `enve`
enve --output json
enve --output xml
```

## Options

```
$ enve -h

NAME:
   enve - run a program in a modified environment using .env files

USAGE:
   enve [global options] command [command options] [arguments...]

DESCRIPTION:
   Set all environment variables of one .env file and run a `command`.

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value    load environment variables from a file path (default: ".env")
   --output value, -o value  output environment variables using text, json or xml format (default: "text")
   --version, -v             shows the current version (default: false)
   --help, -h                show help (default: false)
```

## Contributions

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in current work by you, as defined in the Apache-2.0 license, shall be dual licensed as described below, without any additional terms or conditions.

Feel free to send some [Pull request](https://github.com/joseluisq/enve/pulls) or [issue](https://github.com/joseluisq/enve/issues).

## License

This work is primarily distributed under the terms of both the [MIT license](LICENSE-MIT) and the [Apache License (Version 2.0)](LICENSE-APACHE).

© 2020-present [Jose Quintana](https://git.io/joseluisq)
