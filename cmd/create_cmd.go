package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/aelpxy/dbctl/config"
	"github.com/aelpxy/dbctl/docker"
	"github.com/aelpxy/dbctl/utils"
	"github.com/briandowns/spinner"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// TODO: add support for memory and cpu args
var createCmd = &cobra.Command{
	Use:   "create <database-type>[:image-tag]",
	Short: "Create a new database",
	Long:  "Create a new database.\n\nSupported databases: " + strings.Join(config.SupportedDatabases, ", "),
	Example: `  dbctl create postgres
  dbctl create redis:latest
  dbctl create mysql --password mypassword --port 3306 --name mydb`,
	Aliases:   []string{"mk"},
	Args:      cobra.ExactArgs(1),
	ValidArgs: config.SupportedDatabases,
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		port, _ := cmd.Flags().GetInt("port")
		name, _ := cmd.Flags().GetString("name")

		createDatabase(args[0], password, port, name)
	},
}

func init() {
	createCmd.Flags().StringP("password", "P", "", "Specify a custom password for the database.")
	createCmd.Flags().IntP("port", "p", 0, "Specify a custom port for the database (default: random port).")
	createCmd.Flags().StringP("name", "n", "", "Specify a custom name for the database (default: generated name).")

	rootCmd.AddCommand(createCmd)
}

func createDatabase(part string, password string, port int, name string) {
	dbType, imageVersion := utils.ParseDBTypeAndVersion(part)

	imageTag, err := getImageTag(dbType, imageVersion)
	if err != nil {
		log.Fatal(err)
	}

	if dbType == "mongo" && runtime.GOARCH == "arm" {
		log.Fatal("MongoDB does not support 32-bit arm systems")
	}

	if err := docker.PullImage(imageTag); err != nil {
		log.Fatalf("error pulling image: %v", err)
	}

	if password == "" {
		password = utils.GeneratePassword(16)
	}

	if port == 0 {
		port = utils.GetAvailablePort()
	}

	if name == "" {
		name = utils.GenerateName()
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Creating %s database...", name)
	s.Color("green")
	s.Start()

	containerID, err := docker.CreateContainer(imageTag, dbType, name, port, password, getEnvVarsForDB(dbType, password)...)

	s.Stop()

	if err != nil {
		log.Fatalf("error creating container: %v", err)
	}

	printTable(dbType, name, imageTag, port, password, containerID)
	printConnectionString(dbType, password, port)
}

func getImageTag(dbType, imageVersion string) (string, error) {
	var repository, defaultTag string

	switch dbType {
	case "postgres":
		repository, defaultTag = "postgres", config.PostgresImageTag
	case "redis":
		repository, defaultTag = "redis", config.RedisImageTag
	case "mysql":
		repository, defaultTag = "mysql", config.MySQLImageTag
	case "mariadb":
		repository, defaultTag = "mariadb", config.MariaDBImageTag
	case "mongo":
		repository, defaultTag = "mongo", config.MongoImageTag
	case "meilisearch":
		repository, defaultTag = "getmeili/meilisearch", config.MeiliSearchImageTag
	case "keydb":
		repository, defaultTag = "eqalpha/keydb", config.KeyDBImageTag
	case "couchdb":
		repository, defaultTag = "couchdb", config.CouchDbTag
	case "clickhouse":
		repository, defaultTag = "clickhouse/clickhouse-server", config.ClickHouseTag
	default:
		return "", fmt.Errorf("unsupported database type: %s (supported: %s)", dbType, strings.Join(config.SupportedDatabases, ", "))
	}

	if imageVersion == "" {
		return defaultTag, nil
	}

	return fmt.Sprintf("%s:%s", repository, imageVersion), nil
}

func getEnvVarsForDB(dbType, password string) []string {
	switch dbType {
	case "postgres":
		return []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=postgres",
		}
	case "mysql":
		return []string{
			"MYSQL_ROOT_PASSWORD=" + password,
			"MYSQL_DATABASE=db",
		}
	case "mariadb":
		return []string{
			"MARIADB_ROOT_PASSWORD=" + password,
			"MARIADB_DATABASE=db",
		}
	case "mongo":
		return []string{
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=" + password,
			"MONGO_INITDB_DATABASE=db",
		}
	case "couchdb":
		return []string{
			"COUCHDB_USER=root",
			"COUCHDB_PASSWORD=" + password,
			"COUCHDB_SECRET=" + password,
		}
	case "clickhouse":
		return []string{
			"CLICKHOUSE_DB=db",
			"CLICKHOUSE_USER=root",
			"CLICKHOUSE_PASSWORD=" + password,
			"CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1",
		}
	default:
		return []string{}
	}
}

func printTable(dbType, name, imageTag string, port int, password, containerID string) {
	data := [][]string{
		{"Container ID", containerID[:12]},
		{"Name", config.DockerContainerPrefix + name},
		{"Database Type", dbType},
		{"Image Tag", imageTag},
		{"Port", strconv.Itoa(port)},
		{"Password", password},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(data)
	table.Render()

	fmt.Println()
}

func printConnectionString(dbType, password string, port int) {
	ip := utils.GetIP()

	switch dbType {
	case "postgres":
		fmt.Printf("Connection String: postgres://postgres:%s@%s:%d/postgres\n", password, ip, port)
	case "redis", "keydb":
		fmt.Printf("Connection String: redis://default:%s@%s:%d\n", password, ip, port)
	case "mysql":
		fmt.Printf("Connection String: mysql://root:%s@%s:%d/db\n", password, ip, port)
	case "mariadb":
		fmt.Printf("Connection String: mariadb://root:%s@%s:%d/db\n", password, ip, port)
	case "mongo":
		fmt.Printf("Connection String: mongodb://root:%s@%s:%d/db?authSource=admin\n", password, ip, port)
	case "meilisearch":
		fmt.Printf("URL: http://%s:%d\nMaster Key: %s\n", ip, port, password)
	case "couchdb":
		fmt.Printf("Connection String: http://root:%s@%s:%d\n", password, ip, port)
	case "clickhouse":
		fmt.Printf("Connection String: clickhouse://root:%s@%s:%d/db\n", password, ip, port)
	}
}
