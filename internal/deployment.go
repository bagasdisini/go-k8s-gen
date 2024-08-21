package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Deployment struct {
	APIVersion string             `yaml:"apiVersion"`
	Kind       string             `yaml:"kind"`
	MetaData   MetaDataDeployment `yaml:"metadata"`
	Spec       DeploymentSpec     `yaml:"spec"`
}

type MetaDataDeployment struct {
	Name      string `yaml:"name"`
	NameSpace string `yaml:"namespace"`
	Labels    App    `yaml:"labels"`
}

type App struct {
	App string `yaml:"app"`
}

type DeploymentSpec struct {
	Selector Selector    `yaml:"selector"`
	Template PodTemplate `yaml:"template"`
}

type Selector struct {
	MatchLabels App `yaml:"matchLabels"`
}

type PodTemplate struct {
	Metadata MetaDataPod `yaml:"metadata"`
	Spec     PodSpec     `yaml:"spec"`
}

type MetaDataPod struct {
	Labels App `yaml:"labels"`
}

type PodSpec struct {
	Containers []Container `yaml:"containers"`
}

type Container struct {
	Name            string           `yaml:"name"`
	Image           string           `yaml:"image"`
	ImagePullPolicy string           `yaml:"imagePullPolicy"`
	Ports           []PortDeployment `yaml:"ports"`
	Env             []EnvVar         `yaml:"env"`
}

type PortDeployment struct {
	ContainerPort int `yaml:"containerPort"`
}

type EnvVar struct {
	Name      string        `yaml:"name"`
	ValueFrom *ValueFromRef `yaml:"valueFrom"`
}

type ValueFromRef struct {
	SecretKeyRef SecretKeyRef `yaml:"secretKeyRef"`
}

type SecretKeyRef struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
}

func CreateDeployment() error {
	d := Deployment{
		APIVersion: "v1",
		Kind:       "Deployment",
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

	d.MetaData.Labels.App = projectName + "-" + branchName
	d.MetaData.Name = projectName + "-" + branchName
	d.MetaData.NameSpace = userName

	d.Spec.Selector.MatchLabels.App = projectName + "-" + branchName
	d.Spec.Template.Metadata.Labels.App = projectName + "-" + branchName

	container := Container{
		Name:            projectName + "-" + branchName,
		Image:           fmt.Sprintf("%s/%s:%s", k8sTarget, projectName, branchName),
		ImagePullPolicy: "Always",
		Ports: []PortDeployment{
			{
				ContainerPort: 5000,
			},
		},
	}

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)

	files, err := filepath.Glob(filepath.Join("./k8s", "*.yaml"))
	if err != nil {
		log.Fatalf("Failed to list YAML files: %v", err)
	}

	var secret *Secret
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Failed to read file %s: %v\n", file, err)
			continue
		}

		var genericYaml map[string]interface{}
		err = yaml.Unmarshal(data, &genericYaml)
		if err != nil {
			fmt.Printf("Failed to parse YAML in file %s: %v\n", file, err)
			continue
		}

		_, apiVersionOk := genericYaml["apiVersion"].(string)
		kind, kindOk := genericYaml["kind"].(string)

		if !apiVersionOk || !kindOk {
			fmt.Printf("File %s does not have valid apiVersion or kind: skipping\n", file)
			continue
		}

		if kind == "Secret" {
			var sec Secret
			err := yaml.Unmarshal(data, &sec)
			if err != nil {
				fmt.Printf("Failed to parse Secret YAML in file %s: %v\n", file, err)
				continue
			}
			secret = &sec
		}
	}

	if secret == nil {
		return nil
	}

	var newEnvVars []EnvVar

	for key := range secret.Data {
		newEnvVars = append(newEnvVars, EnvVar{
			Name: strings.ToUpper(key),
			ValueFrom: &ValueFromRef{
				SecretKeyRef: SecretKeyRef{
					Name: secret.MetaData.Name,
					Key:  key,
				},
			},
		})
	}

	sort.Slice(newEnvVars, func(i, j int) bool {
		return newEnvVars[i].Name < newEnvVars[j].Name
	})

	for i := range d.Spec.Template.Spec.Containers {
		d.Spec.Template.Spec.Containers[i].Env = newEnvVars
	}

	yamlData, err := yaml.Marshal(&d)
	if err != nil {
		fmt.Printf("Failed to marshal updated Deployment YAML: %v\n", err)
		return err
	}

	fileName := filepath.Join(Directory, deploymentFileName)

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
