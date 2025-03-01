import Package from "../package.json";

export const CmdVersion = Package.version;
export const CmdName = "dbctl";
export const CmdShortDescription =
  "A CLI tool for managing containerized databases";
export const CmdLongDescription =
  "A command-line tool designed to simplify the management of databases, including creating, deleting, and other operations.";

export const SupportedDatabases: string[] = [
  "redis",
  "postgres",
  "mysql",
  "mariadb",
  "keydb",
  "mongo",
  "meilisearch",
  "couchdb",
  "clickhouse",
];

export const Defaults = {
  DockerContainerPrefix: "dbctl.",
  DockerNetworkName: "dbctl.network",
  DockerVolumeName: "dbctl.volume.",
};
