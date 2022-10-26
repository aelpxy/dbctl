package main

import (
	"fmt"
	"os"

	db "github.com/aelpxy/dbctl/databases"
	"github.com/spf13/cobra"
)

const (
	REDIS_IMAGE    string = "redis:alpine"
	POSTGRES_IMAGE string = "postgres:alpine"
)

func main() {
	var container_port string
	var container_password string

	dbCmd := &cobra.Command{
		Use:       "deploy [redis, postgresql]",
		Short:     "Deploy a database [redis, postgresql]",
		Long:      `Deploy database in a docker container`,
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"redis", "postgresql"},
		Run: func(cmd *cobra.Command, args []string) {
			// value := strings.Join(args, " ")

			if container_password == "" {
				fmt.Println("Container password cannot be null.")
				os.Exit(0)
			}

			switch args[0] {
			case "postgresql":
				db.Create_PostgresDB(container_password, container_port, POSTGRES_IMAGE)
				fmt.Printf("Connection String: postgres://postgres:%s@localhost:%s/postgres \n", container_password, container_port)
			case "redis":
				db.Create_RedisDB(container_password, container_port, REDIS_IMAGE)
				fmt.Printf("Connection String: redis://default:%s@localhost:%s \n", container_password, container_port)
			default:
				fmt.Println("Valid options are redis & postgresql.")
			}
		},
	}

	dbCmd.Flags().StringVarP(&container_port, "port", "p", "", "Port to expose database on.")
	dbCmd.Flags().StringVarP(&container_password, "password", "w", "", "Password to set on database.")

	rootCmd := &cobra.Command{Use: "dbctl"}
	rootCmd.AddCommand(dbCmd)
	rootCmd.Execute()
}
