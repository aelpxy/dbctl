package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	db "github.com/aelpxy/dbctl/databases"
)

const (
	POSTGRES_DEFAULT_IMAGE string = "postgres:14.5-alpine"
	REDIS_DEFAULT_IMAGE    string = "redis:7.0.5-alpine"
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
		log.Println("Creating Postgres DB container...")
		db.Create_PostgresDB(*dbPassword, *dbPort, POSTGRES_DEFAULT_IMAGE)
	}
}
