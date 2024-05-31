package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aelpxy/dbctl/config"
	"github.com/aelpxy/dbctl/docker"
	"github.com/aelpxy/dbctl/utils"
	"github.com/briandowns/spinner"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// todo add support for memory and cpu args
var createCmd = &cobra.Command{
	Use:   "create <database-type>[:image-tag]",
	Short: "Create a new database",
	Long: `Create a new database using the provided options

Supported database types:
- postgres
- redis
- mysql
- mariadb
- mongo

You can optionally specify a custom image tag for the database.`,
	Example: ` 
dbctl create postgres
dbctl create redis:7.2.6-alpine
dbctl create mysql --password mypassword --port 3306 --name mydb`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		port, _ := cmd.Flags().GetInt("port")
		name, _ := cmd.Flags().GetString("name")

		createDatabase(args[0], password, port, name)
	},
}

func init() {
	createCmd.Flags().StringP("password", "P", "", "Specify a custom password for the database")
	createCmd.Flags().IntP("port", "p", 0, "Specify a custom port for the database (default: random port)")
	createCmd.Flags().StringP("name", "n", "", "Specify a custom name for the database (default: generated name)")

	rootCmd.AddCommand(createCmd)
}

func createDatabase(part string, password string, port int, name string) {
	dbType, imageVersion := utils.ParseDBTypeAndVersion(part)
	imageTag := getImageTag(dbType, imageVersion)

	err := docker.PullImage(imageTag)

	if err != nil {
		log.Fatalf("Error pulling image: %v", err)
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

	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = fmt.Sprintf(" Creating %s database... \n", name)
	spinner.Color("green")
	spinner.Start()

	envVars := getEnvVarsForDB(dbType, password)
	containerID, err := createContainerForDB(dbType, name, port, password, envVars)

	if err != nil {
		log.Fatalf("error creating container: %v", err)
	}

	spinner.Stop()

	printTable(dbType, imageTag, port, password, containerID)
	printConnectionString(dbType, password, port)
}

func getImageTag(dbType, imageVersion string) string {
	switch dbType {
	case "postgres":
		if imageVersion == "" {
			return config.PostgresImageTag
		}
		return fmt.Sprintf("postgres:%s", imageVersion)
	case "redis":
		if imageVersion == "" {
			return config.RedisImageTag
		}
		return fmt.Sprintf("redis:%s", imageVersion)
	case "mysql":
		if imageVersion == "" {
			return config.MySQLImageTag
		}
		return fmt.Sprintf("mysql:%s", imageVersion)
	case "mariadb":
		if imageVersion == "" {
			return config.MariaDBImageTag
		}
		return fmt.Sprintf("mariadb:%s", imageVersion)
	case "mongo":
		if imageVersion == "" {
			return config.MongoImageTag
		}
		return fmt.Sprintf("mongo:%s", imageVersion)
	default:
		log.Fatalf("Unsupported database type: %s", dbType)
		return ""
	}
}

func getEnvVarsForDB(dbType, password string) []string {
	switch dbType {
	case "postgres":
		return []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		}
	case "redis":
		return []string{
			fmt.Sprintf("REDIS_PASSWORD=%s", password),
		}
	case "mysql":
		return []string{
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password),
		}
	case "mariadb":
		return []string{
			fmt.Sprintf("MARIADB_ROOT_PASSWORD=%s", password),
		}
	case "mongo":
		return []string{
			fmt.Sprintf("MONGO_INITDB_ROOT_PASSWORD=%s", password),
		}
	default:
		return []string{}
	}
}

func createContainerForDB(dbType, name string, port int, password string, envVars []string) (string, error) {
	var envVarArgs []string

	envVarArgs = append(envVarArgs, envVars...)

	switch dbType {
	case "postgres":
		return docker.CreateContainer(getImageTag("postgres", ""), dbType, name, port, password, envVarArgs...)
	case "redis":
		return docker.CreateContainer(getImageTag("redis", ""), dbType, name, port, password, envVarArgs...)
	case "mysql":
		return docker.CreateContainer(getImageTag("mysql", ""), dbType, name, port, password, envVarArgs...)
	case "mariadb":
		return docker.CreateContainer(getImageTag("mariadb", ""), dbType, name, port, password, envVarArgs...)
	case "mongo":
		return docker.CreateContainer(getImageTag("mongo", ""), dbType, name, port, password, envVarArgs...)
	default:
		return "", fmt.Errorf("unsupported database type: %s", dbType)
	}
}

func printTable(dbType, imageTag string, port int, password, containerID string) {
	data := [][]string{
		{"Container ID", containerID},
		{"Database Type", dbType},
		{"Image Tag", imageTag},
		{"Port", strconv.Itoa(port)},
		{"Password", password},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
	fmt.Println()
}

func printConnectionString(dbType, password string, port int) {
	switch dbType {
	case "postgres":
		fmt.Printf("Connection String: postgres://postgres:%s@%s:%d/postgres\n", password, utils.GetIP(), port)
	case "redis":
		fmt.Printf("Connection String: redis://default:%s@%s:%d\n", password, utils.GetIP(), port)
	case "mysql":
		fmt.Printf("Connection String: mysql://root:%s@%s:%d/db\n", password, utils.GetIP(), port)
	case "mariadb":
		fmt.Printf("Connection String: mariadb://root:%s@%s:%d/db\n", password, utils.GetIP(), port)
	case "mongo":
		fmt.Printf("Connection String: mongodb://root:%s@%s:%d/db?authSource=admin\n", password, utils.GetIP(), port)
	}
}
