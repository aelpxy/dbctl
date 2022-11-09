# dbctl

## Building

Just run this after cloning the repo. (Make sure you have Go & Docker installed.)

```sh
go build
```

## Usage

```sh
  dbctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a docker container
  deploy      Deploy a database [redis, postgresql]
  help        Help about any command

Flags:
  -h, --help   help for dbctl

Use "dbctl [command] --help" for more information about a command.
```

## License
### [MIT](./LICENSE)
