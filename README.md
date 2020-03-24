# fenv [![Build Status](https://travis-ci.com/joseluisq/fenv.svg?branch=master)](https://travis-ci.com/joseluisq/fenv)

> Run a program in a modified environment using `.env` files.

**fenv** just sets all environment variables of one `.env` file and run a command.

## Install

```sh
go get -u github.com/joseluisq/fenv
```

Release binaries also available on [joseluisq/fenv/releases](https://github.com/joseluisq/fenv/releases)

## Usage

A `.env` file is loaded by default from current working directory.

```sh
fenv test.sh
```

Or a custom `.env` file can be loaded using `--file` (`-f`) flag.

```sh
fenv -f dev.env test.sh
```

## Options

```
$ fenv -h

NAME:
   fenv - run a program in a modified environment using .env files

USAGE:
   main [global options] command [command options] [arguments...]

DESCRIPTION:
   Set all environment variables of one .env file and run `command`.

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value  read in a file of environment variables (default: ".env")
   --help, -h              show help (default: false)
```

## Contributions

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in current work by you, as defined in the Apache-2.0 license, shall be dual licensed as described below, without any additional terms or conditions.

Feel free to send some [Pull request](https://github.com/joseluisq/fenv/pulls) or [issue](https://github.com/joseluisq/fenv/issues).

## License

This work is primarily distributed under the terms of both the [MIT license](LICENSE-MIT) and the [Apache License (Version 2.0)](LICENSE-APACHE).

© 2020 [Jose Quintana](https://git.io/joseluisq)
