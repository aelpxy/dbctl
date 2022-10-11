package db

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Create_PostgresDB(password, port, image string) {
	cmd := exec.Command("docker", "run", "-d", "-e", "POSTGRES_PASSWORD="+password, "-p", port+":5432", image)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	cmd.Start()

	buf := bufio.NewReader(stdout)
	for {
		line, _, _ := buf.ReadLine()
		if line == nil {
			log.Println("Container created successfully")
			fmt.Printf("Connection String: postgres://postgres:%s@localhost:%s/postgres \n", password, port)
			os.Exit(0)
		}
		log.Println(string(line))
	}
}
