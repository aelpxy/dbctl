package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	db "github.com/aelpxy/dbctl/databases"
)

const (
	REDIS_IMAGE    string = "bitnami/redis:latest"
	POSTGRES_IMAGE string = "bitnami/postgresql:latest"
	MYSQL_IMAGE    string = "bitnami/mysql:latest"
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
		fmt.Printf("Using %s Image from DockerHub", MYSQL_IMAGE)
		db.Create_PostgresDB(*dbPassword, *dbPort, POSTGRES_IMAGE)
	}

	if *dbType == "mysql" {
		fmt.Println("MySQL")
		fmt.Printf("Using %s Image from DockerHub", MYSQL_IMAGE)
	}

	if *dbType == "redis" {
		fmt.Println("Redis")
		fmt.Printf("Using %s Image from DockerHub", REDIS_IMAGE)
	}
}
