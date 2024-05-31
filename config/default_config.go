package config

var SupportedDatabases = []string{"redis", "postgres", "mysql", "mariadb", "mongo"}

const (
	RedisImageTag    string = "redis:7.2.5-alpine"
	PostgresImageTag string = "postgres:16.3-alpine3.20"
	MySQLImageTag    string = "mysql:8.4.0"
	MariaDBImageTag  string = "mariadb:11.3.2-jammy"
	MongoImageTag    string = "mongo:8.0.0-rc5-jammy"
	VoidImageTag     string = "zotehq/void:alpine" // for future
)

const (
	CmdName             = "dbctl"
	CmdShortDescription = "A CLI tool for managing containerized databases"
	CmdLongDescription  = "A command-line tool designed to simplify the management of databases, including creating, deleting, and other operations."
)

const (
	DockerContainerPrefix string = "dbctl."
	DockerNetworkName     string = "dbctl.network"
	DockerVolumeName      string = "dbctl.volume."
)

const DNSResolverAddress = "9.9.9.9:80"
const Version = "1.1.0"
