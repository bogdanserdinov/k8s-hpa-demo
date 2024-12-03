//go:build mage

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

type service struct {
	name                 string
	version              string
	containerName        string
	additionalMigrations []service
	skipMigrations       bool
}

func (s *service) nameWithVersion() string {
	if s.version == "" {
		return s.name
	}
	return fmt.Sprintf("%v/%v", s.name, s.version)
}

var (
	services = []*service{
		{name: "addition", version: "v1", containerName: "addition-service"},
		{name: "subtraction", version: "v1", containerName: "subtraction-service"},
		{name: "multiplication", version: "v1", containerName: "multiplication-service"},
		{name: "division", version: "v1", containerName: "division-service"},
		{name: "gateway", version: "v1", containerName: "gateway-service"},
	}
)

func Services() error {
	// Clean up the build directory
	sh.Run("rm", "-rf", "build")

	if err := sh.Run("mkdir", "-p", "build"); err != nil {
		return errors.Wrap(err, "failed to create build directory")
	}

	for _, s := range services {
		servicePath := fmt.Sprintf("./services/%v/cmd", s.nameWithVersion())
		if exists, _ := dirExists(servicePath); !exists {
			log.Printf("Skipping %v, not found", servicePath)
			continue
		}

		log.Printf("Building service %v", servicePath)
		buildPath := "build/" + s.name

		err := sh.Run("go", "build", "-o", buildPath, servicePath)
		if err != nil && mg.ExitStatus(err) != 1 {
			return errors.Wrap(err, "failed to build service")
		}
	}

	return nil
}

func Up() error {
	return sh.Run("docker-compose", "up", "-d")
}

func Down() error {
	return sh.Run("docker-compose", "down", "-v", "--rmi=local", "--remove-orphans")
}

type Docker mg.Namespace

func (Docker) BuildImage() error {
	return sh.Run(
		"docker",
		"build",
		"--build-arg", "CACHE_DATE=$(date +%Y-%m-%d:%H:%M:%S)",
		"-t",
		"infra-example",
		".",
	)
}

type Proto mg.Namespace

// Generate the protobuf files for all services,
// Use: mage proto:buf.
func (Proto) Buf() error {
	// Clean up the gen directory
	if err := sh.Run("rm", "-rf", "gen"); err != nil {
		return err
	}

	// Generate the protobuf files for all registered services.
	for _, service := range services {
		path := fmt.Sprintf("proto/%v", service.nameWithVersion())
		// Skip if the service doesn't have a proto file.
		if !doesDirectoryHasFiles(path, ".proto") {
			fmt.Printf("-- Skipping %v, no protobuf files found\n", path)
			continue
		}

		// Generate the protobuf files for the service.
		if err := sh.Run("buf", "generate",
			"--path", fmt.Sprintf("proto/%v", service.nameWithVersion()),
		); err != nil {
			fmt.Printf("❌ Failed to generate protobuf for %v\n", service.name)
		} else {
			fmt.Printf("✅ Generated protobuf for %v\n", service.name)
		}

		os.MkdirAll("docs", 0750)

		protofiles, _ := getProtoFiles(path)
		for _, protofile := range protofiles {
			// Generate documentation for the service.
			err := sh.Run(
				"protoc",
				"-I.",
				"-Ithird_party/proto",
				"-I./proto",
				"--openapiv2_out=./docs",
				"--openapiv2_opt=logtostderr=true",
				fmt.Sprintf("proto/%v/%s", service.nameWithVersion(), protofile),
			)
			if err != nil && mg.ExitStatus(err) != 1 {
				fmt.Printf("❌ Failed to generate documentation for %v\n", service.name)
			} else {
				fmt.Printf("✅ Generated documentation for %v\n", service.name)
			}
		}
	}

	return nil
}

// check if the directory exists
// returns true if it does, false if it doesn't
// or an error if something went wrong
func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Check if the directory not empty
// returns true if it is, false if it is empty or an error if something went wrong
func doesDirectoryHasFiles(path string, ext string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}

	if len(files) == 0 {
		return false
	}

	// check if the directory contains any files with given extension
	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("%s is a directory\n", file.Name())
			continue
		}
		if strings.HasSuffix(file.Name(), ext) {
			return true
		}
	}

	return false
}

// Get proto files from the directory
// returns a slice of proto files or an error if something went wrong
func getProtoFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't read directory")
	}

	var protoFiles []string
	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("%s is a directory\n", file.Name())
			continue
		}
		if strings.HasSuffix(file.Name(), ".proto") {
			protoFiles = append(protoFiles, file.Name())
		}
	}

	return protoFiles, nil
}
