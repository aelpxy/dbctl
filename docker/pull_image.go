package docker

import (
	"fmt"
	"io"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/image"
)

func PullImage(imageName string) error {
	dockerClient, err := DockerClient()
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}

	_, _, err = dockerClient.ImageInspectWithRaw(Ctx, imageName)
	if err == nil {
		fmt.Printf("Image %s already exists, skipping pull\n", imageName)
		return nil
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Pulling %s...", imageName)
	s.Color("green")
	s.Start()

	defer s.Stop()

	out, err := dockerClient.ImagePull(Ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("error pulling image %s: %w", imageName, err)
	}

	defer out.Close()

	// the pull only completes once the progress stream has been fully read
	_, err = io.Copy(io.Discard, out)
	if err != nil {
		return fmt.Errorf("error copying image pull output: %w", err)
	}

	s.Stop()
	fmt.Printf("Image %s pulled successfully\n", imageName)

	return nil
}
