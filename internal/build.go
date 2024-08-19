package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateBuild() error {
	projectName, err := getProjectName()
	if err != nil {
		return fmt.Errorf("failed to get project name: %v", err)
	}

	branchName, err := getCurrentGitBranch()
	if err != nil {
		return fmt.Errorf("failed to get git branch name: %v", err)
	}

	content := fmt.Sprintf(`#!/bin/bash

docker build -t %s/%s:%s .
docker image push %s/%s:%s
`, k8sTarget, projectName, branchName, k8sTarget, projectName, branchName)

	fileName := filepath.Join(Directory, buildFileName)

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}(file)

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}
