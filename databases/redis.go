package db

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Create_RedisDB(password, port, image string) {
	cmd := exec.Command("docker", "run", "-d", "-e", "REDIS_PASSWORD="+password, "-p", port+":6379", image)

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
			fmt.Printf("Connection String: redis://default:%s@localhost:%s/db \n", password, port)
			os.Exit(0)
		}
		log.Println(string(line))
	}
}
