package config

var SupportedDatabases = []string{
	"redis",
	"postgres",
	"mysql",
	"mariadb",
	"mongo",
	"meilisearch",
	"keydb",
	"couchdb",
	"clickhouse",
}

const (
	RedisImageTag       string = "redis:8.4-alpine"
	PostgresImageTag    string = "postgres:18-alpine"
	MySQLImageTag       string = "mysql:9"
	MariaDBImageTag     string = "mariadb:11.8"
	MongoImageTag       string = "mongo:8.0"
	MeiliSearchImageTag string = "getmeili/meilisearch:v1.37"
	KeyDBImageTag       string = "eqalpha/keydb:latest"
	CouchDbTag          string = "couchdb:3.5"
	ClickHouseTag       string = "clickhouse/clickhouse-server:26.3"
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
	DockerTypeLabel       string = "dbctl.type"
)

const DNSResolverAddress = "9.9.9.9:80"
const Version = "1.2.0"
