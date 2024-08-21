package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Service struct {
	APIVersion string          `yaml:"apiVersion"`
	Kind       string          `yaml:"kind"`
	MetaData   MetaDataService `yaml:"metadata"`
	Spec       SpecService     `yaml:"spec"`
}

type MetaDataService struct {
	Labels    App    `yaml:"labels"`
	Name      string `yaml:"name"`
	NameSpace string `yaml:"namespace"`
}

type SpecService struct {
	Ports    []PortService `yaml:"ports"`
	Selector App           `yaml:"selector"`
	Type     string        `yaml:"type"`
}

type PortService struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	NodePort int    `yaml:"nodePort"`
}

func CreateService(portTarget int) error {
	serviceDoc := Service{
		APIVersion: "v1",
		Kind:       "Service",
	}

	projectName, err := getProjectName()
	if err != nil {
		return fmt.Errorf("failed to get project name: %v", err)
	}

	branchName, err := getCurrentGitBranch()
	if err != nil {
		return fmt.Errorf("failed to get git branch name: %v", err)
	}

	userName, err := getCurrentGitUsername()
	if err != nil {
		return fmt.Errorf("failed to get git username: %v", err)
	}

	serviceDoc.MetaData.Labels.App = projectName + "-" + branchName
	serviceDoc.MetaData.Name = projectName + "-" + branchName
	serviceDoc.MetaData.NameSpace = userName

	port := PortService{
		Port:     5000,
		Protocol: "TCP",
		NodePort: portTarget,
	}

	serviceDoc.Spec.Ports = append(serviceDoc.Spec.Ports, port)

	serviceDoc.Spec.Selector.App = projectName + "-" + branchName
	serviceDoc.Spec.Type = "NodePort"

	yamlData, err := yaml.Marshal(&serviceDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal service to YAML: %v", err)
	}

	fileName := filepath.Join(Directory, serviceFileName)

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

	_, err = file.Write(yamlData)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}
