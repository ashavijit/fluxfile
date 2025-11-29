package remote

import (
	"fmt"
	"os"
	"strings"

	"github.com/ashavijit/fluxfile/internal/logger"
	"golang.org/x/crypto/ssh"
)

type Remote struct {
	host   string
	user   string
	logger *logger.Logger
}

func New(connectionString string) (*Remote, error) {
	parts := strings.Split(connectionString, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid connection string format, expected user@host")
	}

	return &Remote{
		user:   parts[0],
		host:   parts[1],
		logger: logger.New(),
	}, nil
}

func (r *Remote) RunCommand(command string, env map[string]string) error {
	config := &ssh.ClientConfig{
		User: r.user,
		Auth: []ssh.AuthMethod{
			r.getSSHAuth(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", r.host+":22", config)
	if err != nil {
		return fmt.Errorf("failed to connect to remote host: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var envVars []string
	for k, v := range env {
		envVars = append(envVars, fmt.Sprintf("export %s=%s", k, v))
	}

	fullCommand := command
	if len(envVars) > 0 {
		fullCommand = strings.Join(envVars, "; ") + "; " + command
	}

	r.logger.Info(fmt.Sprintf("Executing on %s@%s", r.user, r.host))
	r.logger.Command(fullCommand)

	output, err := session.CombinedOutput(fullCommand)
	fmt.Println(string(output))

	if err != nil {
		return fmt.Errorf("remote command failed: %w", err)
	}

	return nil
}

func (r *Remote) getSSHAuth() ssh.AuthMethod {
	keyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(signer)
}

func (r *Remote) CopyFile(localPath, remotePath string) error {
	return fmt.Errorf("file copy not implemented yet")
}
