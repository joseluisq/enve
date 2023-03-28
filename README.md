# Enve ![devel](https://github.com/joseluisq/enve/workflows/devel/badge.svg) [![PkgGoDev](https://pkg.go.dev/badge/github.com/joseluisq/enve)](https://pkg.go.dev/github.com/joseluisq/enve)

> Run a program in a modified environment providing a `.env` file.

**Enve** is a cross-platform tool which can load environment variables from a [`.env` file](https://www.ibm.com/docs/en/aix/7.2?topic=files-env-file) and execute a given command.
It also has the ability to output environment variables in `text`, `json` or `xml` format.

It can be considered as a counterpart of [GNU env](https://www.gnu.org/software/coreutils/manual/html_node/env-invocation.html) command.

## Install

- **Platforms supported:** `linux`, `darwin`, `windows`, `freebsd`, `openbsd`
- **Architectures supported:** `amd64`, `386`, `arm`, `arm64`, `ppc64le`

```sh
curl -sSL \
   "https://github.com/joseluisq/enve/releases/download/v1.4.0/enve_v1.4.0_linux_amd64.tar.gz" \
| sudo tar zxf - -C /usr/local/bin/ enve
```

Using Go:

```sh
go install github.com/joseluisq/enve@latest
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
$ enve 1.4.0
Run a program in a modified environment using .env files

USAGE:
   enve [OPTIONS] COMMAND

OPTIONS:
   -f --file      Load environment variables from a file path (optional) [default: .env]
   -o --output    Output environment variables using text, json or xml format [default: text]
   -h --help      Prints help information
   -v --version   Prints version information
```

## Contributions

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in current work by you, as defined in the Apache-2.0 license, shall be dual licensed as described below, without any additional terms or conditions.

Feel free to send some [Pull request](https://github.com/joseluisq/enve/pulls) or [issue](https://github.com/joseluisq/enve/issues).

## License

This work is primarily distributed under the terms of both the [MIT license](LICENSE-MIT) and the [Apache License (Version 2.0)](LICENSE-APACHE).

© 2020-present [Jose Quintana](https://git.io/joseluisq)
