# dbctl

A CLI built to help you easily manage containerized databases.

## Installation

Before doing anything, make sure you have Docker installed with proper permissions and that you are able to perform operations without using `sudo`.

To install `dbctl` on your system, follow the guide based on your operating system.

### Using a script (Linux/macOS)

First, ensure that `curl` and `tar` are already installed on your OS. Then, execute the following command:

```sh
curl -s https://raw.githubusercontent.com/aelpxy/dbctl/main/scripts/install.sh | bash
```

## Usage

```sh
❯ dbctl help

A command-line tool designed to simplify the management of databases, including creating, deleting, and other operations.

Usage:
  dbctl [flags]
  dbctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create a new database
  delete      Stop and delete a database
  help        Help about any command
  inspect     Inspect a running database
  list        List all running databases
  logs        Stream live logs of a database
  shell       Connect to a running database
  version     Prints the current dbctl version

Flags:
  -h, --help   help for dbctl

Use "dbctl [command] --help" for more information about a command.

```

## Building

Make sure Go (>= 1.20) is installed, then clone the repository:

```sh
git clone git@github.com:aelpxy/dbctl.git
```

There's a `Makefile` to make the build process easier:

```sh
make build # build for your system
make build-all # build for common systems (darwin, windows, linux - arm64/amd64)
```

## Contributing

Pull requests (PRs) are welcome. I recommend maintaining a consistent style of code. When making a PR, please ensure it is detailed enough for me to understand, and the code is self-explanatory, clearly indicating what it does. You are welcome to suggest new features and report any bugs.

## Developing

First of all, make sure Go (>= 1.20) and Docker are installed on your system.

To start developing `dbctl`, clone the repository:

```sh
git clone git@github.com:aelpxy/dbctl.git
```

Create a new branch following the conventional naming schema (not required but preferred - `feat/, fix/, refactor/, chore/`).

Make your changes and then commit the message.

Finally, create a PR.

## License

This repository is licensed under the terms of the [MIT](./LICENSE)
license, as specified in the [LICENSE](./LICENSE)
file.
