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
  dbctl [command]

Available Commands:
  backup      Backup a database (postgres, mysql, mariadb, mongo)
  completion  Generate the autocompletion script for the specified shell
  create      Create a new database
  delete      Stop and delete one or more databases
  help        Help about any command
  http        Start the API and serve it
  inspect     Inspect a running database
  logs        Stream live logs of a database
  ls          List all running databases
  shell       Connect to a running database container
  version     Prints the current dbctl version

Flags:
  -h, --help   help for dbctl

Use "dbctl [command] --help" for more information about a command.
```

### Examples

```sh
dbctl create postgres                         # default image tag, random port, name and password
dbctl create redis:latest                     # custom image tag
dbctl create mysql -P mypassword -p 3306 -n mydb
dbctl ls
dbctl logs <container-id> --tail 100
dbctl backup <container-id> -o backup.sql
dbctl delete <container-id> --force=false     # keep the data volume
```

### Supported databases

| Type          | Default image                        |
| ------------- | ------------------------------------ |
| `postgres`    | `postgres:18-alpine`                 |
| `redis`       | `redis:8.4-alpine`                   |
| `mysql`       | `mysql:9`                            |
| `mariadb`     | `mariadb:11.8`                       |
| `mongo`       | `mongo:8.0`                          |
| `meilisearch` | `getmeili/meilisearch:v1.37`         |
| `keydb`       | `eqalpha/keydb:latest`               |
| `couchdb`     | `couchdb:3.5`                        |
| `clickhouse`  | `clickhouse/clickhouse-server:26.3`  |

### HTTP API

`dbctl http localhost:5000` serves a read-only JSON API:

- `GET /healthcheck`
- `GET /databases`
- `GET /databases/{id}`

## Building

Make sure Go (>= 1.23) is installed, then clone the repository:

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

First of all, make sure Go (>= 1.23) and Docker are installed on your system.

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
