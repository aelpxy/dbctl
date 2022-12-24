package docker

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Backup_Database(id, backup_name string) {
	args := []string{"exec", "-it", id, "pg_dump", "-U", "postgres", "--column-inserts", "--data-only", "postgres", ">", backup_name, ".sql"}
	Command_Exec("docker", args)
}

func Delete_Container(id string) {
	args := []string{"rm", id, "-f"}
	Command_Exec("docker", args)
}

func Pull_Image(name string) {
	args := []string{"pull", name}
	Command_Exec("docker", args)
}

func Purge_Image(name string) {
	args := []string{"image", "rm", name}
	Command_Exec("docker", args)
}

func Create_Network(name string) {
	args := []string{"network", "create", name}
	Command_Exec("docker", args)
}

func List_Containers() {
	args := []string{"ps"}
	Command_Exec("docker", args)
}

func Command_Exec(bin string, arg []string) {
	cmd := exec.Command(bin, arg...)
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
