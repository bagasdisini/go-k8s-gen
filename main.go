package main

import (
	"fmt"
	"go-k8s-gen/internal"
	"math/rand/v2"
	"os"
)

func main() {
	_ = os.MkdirAll(internal.Directory, 0755)

	targetPort := rand.IntN(internal.TargetPortMax-internal.TargetPortMin) + internal.TargetPortMin

	err := internal.CreateBuild()
	if err != nil {
		fmt.Printf("failed to create build: %v", err)
		return
	}

	err = internal.CreateDeploy()
	if err != nil {
		fmt.Printf("failed to create deploy: %v", err)
		return
	}

	err = internal.CreateService(targetPort)
	if err != nil {
		fmt.Printf("failed to create service: %v", err)
		return
	}

	err = internal.CreateSecret(targetPort)
	if err != nil {
		fmt.Printf("failed to create secret: %v", err)
		return
	}
}
