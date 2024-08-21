package internal

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
)

type Secret struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	MetaData   MetaDataSecret         `yaml:"metadata"`
	Data       map[string]interface{} `yaml:"data"`
}

type MetaDataSecret struct {
	CreationTimestamp *string `yaml:"creationTimestamp"`
	Name              string  `yaml:"name"`
	NameSpace         string  `yaml:"namespace"`
}

func CreateSecret(portTarget int) error {
	secretDoc := Secret{
		APIVersion: "v1",
		Kind:       "Secret",
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

	envValue, err := getEnvAndValue()
	if err != nil {
		return fmt.Errorf("failed to get env value: %v", err)
	}

	app := appHost + ":" + strconv.Itoa(portTarget)

	secretDoc.MetaData.CreationTimestamp = nil
	secretDoc.MetaData.Name = projectName + "-" + branchName
	secretDoc.MetaData.NameSpace = userName

	secretDoc.Data = envValue

	secretDoc.Data["app_host"] = "0.0.0.0"
	secretDoc.Data["app_port"] = 5000

	if _, ok := secretDoc.Data["swagger_host"]; ok {
		secretDoc.Data["swagger_host"] = app
	}

	if _, ok := secretDoc.Data["cors_allow_origins"]; ok {
		secretDoc.Data["cors_allow_origins"] = secretDoc.Data["cors_allow_origins"].(string) + ",http://" + app
	}

	for key, value := range secretDoc.Data {
		var encoded string
		switch v := value.(type) {
		case int:
			encoded = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(v)))
		case bool:
			encoded = base64.StdEncoding.EncodeToString([]byte(strconv.FormatBool(v)))
		case string:
			encoded = base64.StdEncoding.EncodeToString([]byte(v))
		default:
			encoded = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", v)))
		}
		secretDoc.Data[key] = encoded
	}

	yamlData, err := yaml.Marshal(&secretDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal service to YAML: %v", err)
	}

	fileName := filepath.Join(Directory, secretFileName)

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
