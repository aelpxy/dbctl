package docker

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Command_Exec(bin, arg string) {
	cmd := exec.Command(bin, arg)
	r, _ := cmd.StdoutPipe()

	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}

		done <- struct{}{}
	}()

	cmd.Start()

	<-done

	cmd.Wait()
}
