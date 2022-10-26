package db

import (
	"bufio"
	"log"
	"os/exec"
)

func Create_PostgresDB(password, port, image string) {
	cmd := exec.Command("docker", "run", "-d", "-e", "POSTGRES_PASSWORD="+password, "-p", port+":5432", image)
	r, _ := cmd.StdoutPipe()

	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
		}

		done <- struct{}{}
	}()

	cmd.Start()

	<-done

	cmd.Wait()
}
