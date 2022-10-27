package db

import (
	"github.com/aelpxy/dbctl/docker"
)

func Create_PostgresDB(password, port, image string) {
	args := []string{"run", "-d", "-e", "POSTGRES_PASSWORD=" + password, "-p", port + ":5432", image}
	docker.Command_Exec("docker", args)
}
