# enve [![Build Status](https://travis-ci.com/joseluisq/enve.svg?branch=master)](https://travis-ci.com/joseluisq/enve)

> Run a program in a modified environment providing a .env file.

**enve** loads all environment variables of a `.env` file and run a command. It can be considered as a counterpart of [GNU env](https://www.gnu.org/software/coreutils/manual/html_node/env-invocation.html) command.

## Install

```sh
go get -u github.com/joseluisq/enve
```

Release binaries also available on [joseluisq/enve/releases](https://github.com/joseluisq/enve/releases)

## Usage

A `.env` file is loaded by default from current working directory.

```sh
enve test.sh
```

Or a custom `.env` file can be loaded using `--file` (`-f`) flag.

```sh
enve -f dev.env test.sh
```

## Options

```
$ enve -h

NAME:
   enve - run a program in a modified environment using .env files

USAGE:
   enve [global options] command [command options] [arguments...]

DESCRIPTION:
   Set all environment variables of one .env file and run `command`.

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value    read in a file of environment variables (default: ".env")
   --output value, -o value  output environment variables in specific format (default: "text")
   --version, -v             shows the current version (default: false)
   --help, -h                show help (default: false)
```

## Contributions

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in current work by you, as defined in the Apache-2.0 license, shall be dual licensed as described below, without any additional terms or conditions.

Feel free to send some [Pull request](https://github.com/joseluisq/enve/pulls) or [issue](https://github.com/joseluisq/enve/issues).

## License

This work is primarily distributed under the terms of both the [MIT license](LICENSE-MIT) and the [Apache License (Version 2.0)](LICENSE-APACHE).

© 2020-present [Jose Quintana](https://git.io/joseluisq)
