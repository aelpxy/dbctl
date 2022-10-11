package main

import (
	"flag"
	"fmt"
	"os"

	db "github.com/aelpxy/dbctl/databases"
)

const (
	REDIS_IMAGE    string = "redis:alpine"
	POSTGRES_IMAGE string = "postgres:alpine"
)

func main() {
	dbPassword := flag.String("password", "password", "Password to use for that database.")
	dbType := flag.String("db", "none", "Type of database to deploy.")
	dbPort := flag.String("p", "5432", "The port on which the database will listen.")

	flag.Parse()

	if *dbType == "none" {
		fmt.Println("Need help? Use the -h flag for more information.")
		os.Exit(0)
	}

	if *dbType == "postgres" {

		fmt.Println("Deploying Database...")
		fmt.Printf("Using %s Image from DockerHub \n", POSTGRES_IMAGE)
		db.Create_PostgresDB(*dbPassword, *dbPort, POSTGRES_IMAGE)
	}

	if *dbType == "redis" {

		fmt.Println("Deploying Database...")
		fmt.Printf("Using %s Image from DockerHub \n", REDIS_IMAGE)
		db.Create_RedisDB(*dbPassword, *dbPort, REDIS_IMAGE)
	}
}
