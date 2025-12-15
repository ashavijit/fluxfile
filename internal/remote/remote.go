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
	// Try multiple key locations
	keyPaths := []string{
		os.Getenv("HOME") + "/.ssh/id_rsa",
		os.Getenv("HOME") + "/.ssh/id_ed25519",
		os.Getenv("USERPROFILE") + "/.ssh/id_rsa", // Windows
		os.Getenv("USERPROFILE") + "/.ssh/id_ed25519",
	}

	for _, keyPath := range keyPaths {
		key, err := os.ReadFile(keyPath)
		if err != nil {
			continue
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			continue
		}

		return ssh.PublicKeys(signer)
	}

	return nil
}

// CopyFile copies a local file to the remote server using SCP
func (r *Remote) CopyFile(localPath, remotePath string) error {
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

	// Read local file
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %w", err)
	}

	// Get file info for permissions
	info, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("failed to stat local file: %w", err)
	}

	// Set up stdin pipe
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Start scp command on remote
	go func() {
		defer stdin.Close()
		fmt.Fprintf(stdin, "C%04o %d %s\n", info.Mode().Perm(), len(content), remotePath)
		stdin.Write(content)
		fmt.Fprint(stdin, "\x00")
	}()

	r.logger.Info(fmt.Sprintf("Copying %s to %s@%s:%s", localPath, r.user, r.host, remotePath))

	if err := session.Run("scp -t " + remotePath); err != nil {
		return fmt.Errorf("scp failed: %w", err)
	}

	return nil
}
