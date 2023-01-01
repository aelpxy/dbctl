package main

import (
	"fmt"
	"log"
	"net"
	"os"

	db "github.com/aelpxy/dbctl/databases"
	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

const (
	REDIS_IMAGE    string = "redis:alpine"
	POSTGRES_IMAGE string = "postgres:alpine"
)

var (
	container_port     string
	container_password string
)

func main() {
	dbCmd := &cobra.Command{
		Use:       "deploy [redis, postgresql]",
		Short:     "Deploy a database [redis, postgresql]",
		Long:      "Deploy database in a docker container",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"redis", "postgresql"},
		Run: func(cmd *cobra.Command, args []string) {
			if container_password == "" {
				fmt.Println("Database password was not provided")
				os.Exit(0)
			}

			if container_port == "" {
				fmt.Println("Database port was not provided")
				os.Exit(0)
			}

			switch args[0] {
			case "postgresql":
				db.Create_PostgresDB(container_password, container_port, POSTGRES_IMAGE)
				fmt.Printf("Connection String: postgres://postgres:%s@%s:%s/postgres \n", container_password, GetIP(), container_port)
			case "redis":
				db.Create_RedisDB(container_password, container_port, REDIS_IMAGE)
				fmt.Printf("Connection String: redis://default:%s@%s:%s \n", container_password, GetIP(), container_port)
			default:
				fmt.Println("Valid options are redis & postgresql.")
			}
		},
	}

	dbCmd.Flags().StringVarP(&container_port, "port", "p", "", "Port to expose database on.")
	dbCmd.Flags().StringVarP(&container_password, "password", "w", "", "Password to set on database.")

	dbDeleteCmd := &cobra.Command{
		Use:   "delete [container id]",
		Short: "Delete a docker container",
		Long:  "Delete a docker container using its id",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] == "" {
				fmt.Println("id cannot be null")
				os.Exit(0)
			}
			docker.Delete_Container(args[0])
		},
	}

	dbBackupCmd := &cobra.Command{
		Use:   "backup [container id] [backup file name]",
		Short: "Backup PostgreSQL container",
		Long:  "Backup PostgreSQL container using its id",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] == "" {
				fmt.Println("id cannot be null")
				os.Exit(0)
			}

			if args[1] == "" {
				fmt.Println("name cannot be null")
				os.Exit(0)
			}

			docker.Backup_Database(args[0], args[1])
		},
	}

	dbListCmd := &cobra.Command{
		Use:   "list",
		Short: "List docker containers",
		Long:  "List all the docker containers",
		// Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			docker.List_Containers()
		},
	}

	rootCmd := &cobra.Command{Use: "dbctl"}
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(dbDeleteCmd)
	rootCmd.AddCommand(dbBackupCmd)
	rootCmd.AddCommand(dbListCmd)

	rootCmd.Execute()
}

func GetIP() net.IP {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
