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
  backup      Backup a database
  completion  Generate the autocompletion script for the specified shell
  connect     Open a local port that forwards to a database
  create      Create a new database
  delete      Stop and delete one or more databases
  help        Help about any command
  http        Start the API and serve it
  inspect     Inspect a database
  logs        Stream live logs of a database
  ls          List all databases
  restart     Restart one or more databases
  restore     Restore a database from a backup file
  shell       Open the database client, or a shell with --sh
  start       Start one or more stopped databases
  stop        Stop one or more databases without deleting them
  url         Print the connection string and credentials of a database
  version     Prints the current dbctl version

Flags:
  -h, --help   help for dbctl

Use "dbctl [command] --help" for more information about a command.
```

### Examples

```sh
dbctl create postgres                         # default image, random port, name and password; waits until ready
dbctl create redis:latest                     # custom image tag
dbctl create mysql -P mypassword -p 3306 -n mydb
dbctl create pgvector -o json                 # machine-readable output
dbctl ls -o json
dbctl url <name>                      # print the connection string again
dbctl shell <name>                    # open psql, redis-cli, mongosh, ...
dbctl shell <name> --sh               # plain /bin/sh instead
dbctl connect <name>                          # forward 127.0.0.1:<port> to the database until Ctrl-C
dbctl shell                                   # no name: pick a database from a list
dbctl stop <name>                     # keep the data, free the memory
dbctl start <name>
dbctl backup <name> -o backup.sql
dbctl restore <name> backup.sql
dbctl delete <name>                           # asks you to type the name to confirm
dbctl delete <name> --yes --force=false       # no prompt, keep the data volume
```

Commands that take a database open an interactive picker when you leave the name out (in a terminal), and accept the short name shown by `dbctl ls` (e.g. `misty-river-bold-pine`) or a container ID prefix. `ls`, `inspect`, `create` and `url` accept `-o json`. Progress output goes to stderr, so stdout stays pipeable, and colors are disabled when output is piped or `NO_COLOR` is set.

### Shell completion

Completion suggests database names and types. Enable it for your shell, for example:

```sh
dbctl completion bash > /etc/bash_completion.d/dbctl     # bash
dbctl completion zsh > "${fpath[1]}/_dbctl"              # zsh
dbctl completion fish > ~/.config/fish/completions/dbctl.fish
```

### Supported databases

| Type          | Default image                                         | Client        | Backup / restore |
| ------------- | ----------------------------------------------------- | ------------- | ---------------- |
| `postgres`    | `postgres:18-alpine`                                  | `psql`        | yes              |
| `pgvector`    | `pgvector/pgvector:pg18`                              | `psql`        | yes              |
| `mysql`       | `mysql:9`                                             | `mysql`       | yes              |
| `mariadb`     | `mariadb:11.8`                                        | `mariadb`     | yes              |
| `mongo`       | `mongo:8.0`                                           | `mongosh`     | yes              |
| `redis`       | `redis:8.4-alpine`                                    | `redis-cli`   |                  |
| `valkey`      | `valkey/valkey:9.1-alpine`                            | `valkey-cli`  |                  |
| `keydb`       | `eqalpha/keydb:latest`                                | `keydb-cli`   |                  |
| `dragonfly`   | `docker.dragonflydb.io/dragonflydb/dragonfly:v1.37.0` |               |                  |
| `meilisearch` | `getmeili/meilisearch:v1.37`                          |               |                  |
| `couchdb`     | `couchdb:3.5`                                         |               |                  |
| `clickhouse`  | `clickhouse/clickhouse-server:26.3`                   | `clickhouse-client` |            |

### Database templates

Each database is a JSON template in [`internal/database/templates`](./internal/database/templates), embedded into the binary. To add a database, drop in a new file:

```json
{
  "name": "postgres",
  "repository": "postgres",
  "tag": "18-alpine",
  "port": 5432,
  "data_dir": "/var/lib/postgresql",
  "env": ["POSTGRES_PASSWORD={{.Password}}", "POSTGRES_USER=postgres", "POSTGRES_DB=postgres"],
  "url": "postgres://postgres:{{urlquery .Password}}@{{.Host}}/postgres",
  "backup": {
    "command": ["pg_dump", "-U", "{{.Env.POSTGRES_USER}}", "{{.Env.POSTGRES_DB}}"],
    "extension": "sql"
  }
}
```

Optional fields:

- `env` and `command`: the container environment and command
- `healthcheck`: a readiness probe; `create` waits for it to pass and `ls` shows the health
- `client`: what `dbctl shell` opens (falls back to `/bin/sh`)
- `backup` (`command`, `extension`, optional `env`): writes a dump to stdout
- `restore` (`command`, optional `env`): reads a dump from stdin
- `unsupported_arch`: a list of `GOARCH` values the image does not run on
- `required_cpu_features`: `/proc/cpuinfo` flags the image needs, per `GOARCH` (e.g. ClickHouse and MongoDB need ARMv8.2-A on arm64, so a Raspberry Pi 4 is refused up front) Strings are Go [`text/template`](https://pkg.go.dev/text/template)s with these fields:

- `{{.Password}}`: the database password, in `env`, `command`, `healthcheck`, `client` and `url`
- `{{.Host}}`: the published `host:port`, in `url`
- `{{.Env.NAME}}`: the container's environment variables, in `client`, `backup` and `restore`; a missing variable is an error

### HTTP API

`dbctl http localhost:5000` serves a read-only JSON API:

- `GET /healthcheck`
- `GET /databases`
- `GET /databases/{id}`

## Building

Make sure Go (>= 1.27) is installed, then clone the repository:

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

First of all, make sure Go (>= 1.27) and Docker are installed on your system.

To start developing `dbctl`, clone the repository:

```sh
git clone git@github.com:aelpxy/dbctl.git
```

Create a new branch following the conventional naming schema (not required but preferred - `feat/, fix/, refactor/, chore/`).

The entrypoint lives in `cmd/dbctl` and everything else under `internal/`. Before committing, make sure the code is formatted, lint-free and tested:

```sh
make fmt   # golangci-lint fmt
make lint  # golangci-lint run ./...
make test  # go test -race ./...
```

Make your changes and then commit the message.

Finally, create a PR.

## License

This repository is licensed under the terms of the [MIT](./LICENSE)
license, as specified in the [LICENSE](./LICENSE)
file.
