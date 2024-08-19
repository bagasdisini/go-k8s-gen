package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDeploy() error {
	userName, err := getCurrentGitUsername()
	if err != nil {
		return fmt.Errorf("failed to get git username: %v", err)
	}

	content := fmt.Sprintf(`#!/bin/bash

kubectl delete -f %s -n %s
kubectl delete -f %s -n %s
kubectl delete -f %s -n %s

kubectl apply -f %s -n %s
kubectl apply -f %s -n %s
kubectl apply -f %s -n %s
`, secretFileName, userName, serviceFileName, userName, deploymentFileName, userName, secretFileName, userName, serviceFileName, userName, deploymentFileName, userName)

	fileName := filepath.Join(Directory, deployFileName)

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
