# Kubernetes Config Generator for Go Project

This is a simple tool to generate Kubernetes configuration files for a Go project. It generates the following files:
- deployment.yaml
- service.yaml
- secret.yaml
- build.sh
- deploy.sh

## Prerequisites
- Make sure to initialize git in your project directory.
- Make sure to have a `.env` file in your project directory.

## Usage
If you're using windows, you can download the binary from the release page.
Or if you want to build it yourself, follow these steps:

1. Clone this repository.
```bash
git clone https://github.com/bagasdisini/go-k8s-gen
cd go-k8s-gen
```
2. Optional: Modify the config in the `internal/gear` as needed.
3. Build the app.
```bash
go build 
```
4. Copy the binary to your Go project directory.
5. Run the binary.
```bash
./go-k8s-gen
```
6. The generated files will be in the `k8s` directory.