package config

var SupportedDatabases = []string{"redis", "postgres", "mysql", "mariadb", "mongo", "meilisearch", "keydb", "couchdb", "clickhouse"}

const (
	RedisImageTag       string = "redis:7.4.1-alpine3.20"
	PostgresImageTag    string = "postgres:17.0-alpine3.19"
	MySQLImageTag       string = "mysql:9.0.1"
	MariaDBImageTag     string = "mariadb:11.2.5-jammy"
	MongoImageTag       string = "mongo:8.0.1-noble"
	MeiliSearchImageTag string = "getmeili/meilisearch:v1.10.3"
	KeyDBImageTag       string = "eqalpha/keydb:latest"
	CouchDbTag          string = "couchdb:3.4.2"
	ClickHouseTag       string = "clickhouse/clickhouse-server:24.3.12-alpine"
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
const Version = "1.2.0"
