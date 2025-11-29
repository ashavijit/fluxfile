package docker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ashavijit/fluxfile/internal/logger"
)

type Docker struct {
	image  string
	logger *logger.Logger
}

func New(image string) *Docker {
	if image == "" {
		image = "alpine:latest"
	}
	return &Docker{
		image:  image,
		logger: logger.New(),
	}
}

func (d *Docker) RunCommand(command string, env map[string]string, workdir string) error {
	args := []string{"run", "--rm"}

	if workdir != "" {
		args = append(args, "-v", fmt.Sprintf("%s:/workspace", workdir))
		args = append(args, "-w", "/workspace")
	}

	for k, v := range env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, d.image, "sh", "-c", command)

	cmd := exec.Command("docker", args...)

	d.logger.Info(fmt.Sprintf("Running in Docker: %s", d.image))
	d.logger.Command(strings.Join(cmd.Args, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return fmt.Errorf("docker command failed: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func (d *Docker) IsAvailable() bool {
	cmd := exec.Command("docker", "version")
	return cmd.Run() == nil
}

func (d *Docker) PullImage() error {
	d.logger.Info(fmt.Sprintf("Pulling Docker image: %s", d.image))
	cmd := exec.Command("docker", "pull", d.image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return fmt.Errorf("failed to pull image: %w", err)
	}
	return nil
}
