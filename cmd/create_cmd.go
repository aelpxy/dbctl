package cmd

import (
	"fmt"
	"log"
	"os"
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

var createCmd = &cobra.Command{
	Use:   "create <database-type>[:image-tag]",
	Short: "Create a new database",
	Long: `Create a new database using the provided options.

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
	parts := strings.Split(part, ":")
	dbType := parts[0]
	imageVersion := ""

	if len(parts) > 1 {
		imageVersion = parts[1]
	}

	var imageTag string

	if password == "" {
		password = utils.GeneratePassword(16)
	}

	if port == 0 {
		port = utils.GetAvailablePort()
	}

	if name == "" {
		name = utils.GenerateName()
	}

	switch dbType {
	case "postgres":
		imageTag = config.PostgresImageTag
		if imageVersion != "" {
			imageTag = fmt.Sprintf("postgres:%s", imageVersion)
		}

		err := docker.PullImage(imageTag)

		if err != nil {
			log.Fatalf("Error pulling image: %v", err)
		}

		spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		spinner.Suffix = fmt.Sprintf(" Creating %s database... \n", name)
		spinner.Color("green")
		spinner.Start()

		containerID, err := docker.CreateContainer(imageTag, dbType, name, port, password, []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		})

		if err != nil {
			log.Fatalf("error creating container: %v", err)
		}

		spinner.Stop()

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

		// adds a space
		fmt.Println()

		fmt.Printf("Connection String: postgres://postgres:%s@%s:%s/postgres\n", password, utils.GetIP(), strconv.Itoa(port))
		// these ones are TODO
	case "redis":
		imageTag = config.RedisImageTag
		if imageVersion != "" {
			imageTag = fmt.Sprintf("redis:%s", imageVersion)
		}

		err := docker.PullImage(imageTag)

		if err != nil {
			log.Fatalf("Error pulling image: %v", err)
		}
	case "mysql":
		imageTag = config.MySQLImageTag
		if imageVersion != "" {
			imageTag = fmt.Sprintf("mysql:%s", imageVersion)
		}

		err := docker.PullImage(imageTag)

		if err != nil {
			log.Fatalf("Error pulling image: %v", err)
		}
	case "mariadb":
		imageTag = config.MariaDBImageTag
		if imageVersion != "" {
			imageTag = fmt.Sprintf("mariadb:%s", imageVersion)
		}

		err := docker.PullImage(imageTag)

		if err != nil {
			log.Fatalf("Error pulling image: %v", err)
		}
	case "mongo":
		imageTag = config.MongoImageTag
		if imageVersion != "" {
			imageTag = fmt.Sprintf("mongo:%s", imageVersion)
		}

		err := docker.PullImage(imageTag)

		if err != nil {
			log.Fatalf("Error pulling image: %v", err)
		}
	default:
		log.Fatalf("Unsupported database type: %s", dbType)
	}

}
