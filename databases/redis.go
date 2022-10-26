package db

import (
	"bufio"
	"log"
	"os/exec"
)

func Create_RedisDB(password, port, image string) {
	cmd := exec.Command("docker", "run", "-d", "-e", "REDIS_PASSWORD="+password, "-p", port+":6379", image)

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
