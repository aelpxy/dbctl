package main

import (
	"flag"
	"fmt"

	util "github.com/aelpxy/dbctl/utils"
)

func main() {
	containerName := flag.String("name", "postgres-db", "Name to use for the database container.")
	dbType := flag.String("db", "postgres", "Type of database to deploy.")
	dbPort := flag.String("port", "5432", "The port on which the database will listen.")

	if *dbType == "postgres" {
		fmt.Printf("Port: %s | Type: %s, Name: %s \n", *dbPort, *dbType, *containerName)
		util.Pull_Image("postgres:14.5-alpine")
	}

	flag.Parse()
}
