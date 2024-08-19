package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	buildFileName  = "build.sh"
	deployFileName = "deploy.sh"

	deploymentFileName = "go-deployment.yaml"
	serviceFileName    = "go-service.yaml"
	secretFileName     = "go-secret.yaml"

	k8sTarget = "192.168.1.2:5000"
	appHost   = "192.168.1.2"

	Directory     = "k8s"
	TargetPortMin = 30200
	TargetPortMax = 30300
)

func getProjectName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}
	return "", fmt.Errorf("project name not found in go.mod")
}

func getCurrentGitBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git branch name: %v", err)
	}

	return strings.TrimSpace(string(out)), nil
}

func getCurrentGitUsername() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git branch name: %v", err)
	}

	return strings.ToLower(strings.TrimSpace(string(out))), nil
}

func getEnvAndValue() (map[string]interface{}, error) {
	data, err := os.ReadFile(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %v", err)
	}

	envMap := make(map[string]interface{})
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				envMap[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
			}
		}
	}
	return envMap, nil
}
