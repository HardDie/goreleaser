package linux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/HardDie/goreleaser/internal/logger"
	"github.com/HardDie/goreleaser/internal/utils"
)

func Build(name, imagePath, version, license, path, ldflags string) error {
	newWorkDirectory := filepath.Dir(path)
	entryPointFile := filepath.Base(path)

	currentDirBackup, err := os.Getwd()
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error get work directory: %w", err)
	}

	err = os.Chdir(newWorkDirectory)
	if err != nil {
		return fmt.Errorf("error change directory to %s: %w", newWorkDirectory, err)
	}

	arches := []string{"amd64", "386", "arm64"}
	for _, arch := range arches {
		// Compile app
		cmd := exec.Command("go", "build", "-a", "-o", name, entryPointFile)
		if ldflags != "" {
			cmd.Args = append(cmd.Args, "-ldflags", ldflags)
		}
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
		cmd.Env = append(cmd.Env, "GOARCH="+arch)
		cmd.Env = append(cmd.Env, "GOOS=linux")
		logger.Debug.Println("Execute:", cmd.String())
		err = cmd.Run()
		if err != nil {
			logger.Error.Println(err)
			return fmt.Errorf("error building application: %w", err)
		}

		// Create tar archive
		cmd = exec.Command("tar", "-czf", "../../release/"+name+".linux-"+arch+".tar.gz", name)
		logger.Debug.Println("Execute:", cmd.String())
		err = cmd.Run()
		if err != nil {
			logger.Error.Println(err)
			return fmt.Errorf("error creation archive: %w", err)
		}

		// Remove binary file
		err = utils.RemoveFile(name)
		if err != nil {
			return err
		}
	}

	// Return to the root folder
	err = os.Chdir(currentDirBackup)
	if err != nil {
		return fmt.Errorf("error return from directory: %w", err)
	}

	return nil
}
