package main

import (
	"flag"
	"fmt"

	util "github.com/aelpxy/dbctl/utils"
)

const (
	POSTGRES_DEFAULT_IMAGE string = "postgres:14.5-alpine"
	REDIS_DEFAULT_IMAGE    string = "redis:7.0.5-alpine"
)

func main() {
	containerName := flag.String("name", "postgres-db", "Name to use for the database container.")
	dbType := flag.String("db", "postgres", "Type of database to deploy.")
	dbPort := flag.String("port", "5432", "The port on which the database will listen.")

	if *dbType == "postgres" {
		fmt.Printf("Container Name: %s", *containerName)
		fmt.Printf("Seletcted Database: %s", *dbType)
		fmt.Printf("Seletcted Port: %s", *dbPort)

		util.Pull_Image(POSTGRES_DEFAULT_IMAGE)
	}

	if *dbType == "redis" {
		util.Pull_Image(REDIS_DEFAULT_IMAGE)
	}

	flag.Parse()
}
