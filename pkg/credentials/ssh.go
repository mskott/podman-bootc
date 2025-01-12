package credentials

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/containers/podman-bootc/pkg/config"
)

// Generatekeys creates an RSA set of keys
func Generatekeys(outputDir string) (string, error) {
	sshIdentity := filepath.Join(outputDir, config.SshKeyFile)
	_ = os.Remove(sshIdentity)
	_ = os.Remove(sshIdentity + ".pub")

	// we use RSA here so it works on FIPS mode
	args := []string{"-N", "", "-t", "rsa", "-f", sshIdentity}
	cmd := exec.Command("ssh-keygen", args...)
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("ssh key generation: redirecting stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("ssh key generation: executing ssh-keygen: %w", err)
	}

	waitErr := cmd.Wait()
	if waitErr == nil {
		return sshIdentity, nil
	}

	errMsg, err := io.ReadAll(stdErr)
	if err != nil {
		return "", fmt.Errorf("ssh key generation, unable to read from stderr: %w", waitErr)
	}

	return "", fmt.Errorf("failed to generate ssh keys: %s: %w", string(errMsg), waitErr)
}
