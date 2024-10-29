package main

import (
	"github.com/aelpxy/dbctl/cmd"
)

func main() {
	// if !utils.IsDockerInstalled() {
	// 	log.Fatalf("Docker is not installed please install docker and try again")
	// }

	cmd.Execute()
}
