package main

import (
	"log"

	"github.com/aelpxy/dbctl/cmd"
	"github.com/aelpxy/dbctl/utils"
)

func main() {
	if !utils.IsDockerInstalled() {
		log.Fatalf("Docker is not installed please install docker and try again")
	}

	cmd.Execute()
}
