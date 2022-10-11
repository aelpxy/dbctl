package main

import (
	"flag"
	"fmt"
	"os"

	util "github.com/aelpxy/dbctl/utils"
)

const (
	POSTGRES_DEFAULT_IMAGE string = "postgres:14.5-alpine"
	REDIS_DEFAULT_IMAGE    string = "redis:7.0.5-alpine"
)

func main() {
	containerName := flag.String("name", "default", "Name to use for the database container (Don't use default name).")
	dbType := flag.String("db", "postgres", "Type of database to deploy.")
	dbPort := flag.String("port", "5432", "The port on which the database will listen.")

	flag.Parse()

	if *containerName == "default" {
		fmt.Println("Need help? Use the -h flag for more information.")
		os.Exit(0)
	}

	if *dbType == "postgres" {
		fmt.Printf("Container Name: %s \n", *containerName)
		fmt.Printf("Seletcted Database: %s \n", *dbType)
		fmt.Printf("Seletcted Port: %s \n", *dbPort)

		util.Pull_Image(POSTGRES_DEFAULT_IMAGE)
	}

	if *dbType == "redis" {
		fmt.Printf("Container Name: %s \n", *containerName)
		fmt.Printf("Seletcted Database: %s \n", *dbType)
		fmt.Printf("Seletcted Port: %s \n", *dbPort)
		util.Pull_Image(REDIS_DEFAULT_IMAGE)
	}
}
