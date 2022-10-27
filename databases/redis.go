package db

import (
	"github.com/aelpxy/dbctl/docker"
)

func Create_RedisDB(password, port, image string) {
	args := []string{"run", "-d", "-e", "REDIS_PASSWORD=" + password, "-p", port + ":6379", image}
	docker.Command_Exec("docker", args)
}
