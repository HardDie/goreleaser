package darwin

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/HardDie/goreleaser/internal/image"
	"github.com/HardDie/goreleaser/internal/logger"
	"github.com/HardDie/goreleaser/internal/utils"
)

const (
	infoTmpl = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>CFBundleExecutable</key>
        <string>{{.Name}}</string>
        <key>CFBundleIconFile</key>
        <string>icon.icns</string>
        <key>CFBundleIdentifier</key>
        <string>{{.Company}}</string>
        <key>NSHighResolutionCapable</key>
        <true/>
        <key>LSUIElement</key>
        <true/>
</dict>
</plist>`
)

type Info struct {
	Name    string
	Company string
}

func Build(name, imagePath, version, license, path, company string) error {
	err := image.ConvertToDarwinIconsContainer(imagePath, filepath.Join("build_cache", "icon.icns"))
	if err != nil {
		return err
	}

	newWorkDirectory, _ := filepath.Abs(filepath.Dir(path))
	entryPointFile := filepath.Base(path)

	log.Println("work:", newWorkDirectory)

	currentDirBackup, err := os.Getwd()
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error get work directory: %w", err)
	}
	buildCacheDirectory := filepath.Join(currentDirBackup, "build_cache")

	arches := []string{"amd64", "arm64"}
	for _, arch := range arches {
		err = os.Chdir(newWorkDirectory)
		if err != nil {
			return fmt.Errorf("error change directory to %s: %w", newWorkDirectory, err)
		}

		// Compile app
		cmd := exec.Command("go", "build", "-a", "-o", filepath.Join(buildCacheDirectory, name), entryPointFile)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
		cmd.Env = append(cmd.Env, "GOARCH="+arch)
		cmd.Env = append(cmd.Env, "GOOS=darwin")
		logger.Debug.Println("Execute:", cmd.String())
		err = cmd.Run()
		if err != nil {
			logger.Error.Println(err)
			return fmt.Errorf("error building application: %w", err)
		}

		err = os.Chdir(buildCacheDirectory)
		if err != nil {
			return fmt.Errorf("error change directory to %s: %w", buildCacheDirectory, err)
		}

		// Prepare folders
		err = CreateDirectoryHierarchy(name)
		if err != nil {
			return err
		}

		// Move binary into app folder
		err = os.Rename(name, filepath.Join(name+".app", "Contents", "MacOS", name))
		if err != nil {
			logger.Error.Println(err)
			return fmt.Errorf("error moving binary file into MacOS directory: %w", err)
		}

		// Copy icons container into app folder
		err = utils.CopyFile("icon.icns", filepath.Join(name+".app", "Contents", "Resources", "icon.icns"))
		if err != nil {
			return err
		}

		// Create Info.plist file
		err = CreateInfoPlistFile(name, company)
		if err != nil {
			return err
		}

		// Create tar archive
		cmd = exec.Command("tar", "-czf", "../release/"+name+".darwin-"+arch+".tar.gz", name+".app")
		logger.Debug.Println("Execute:", cmd.String())
		err = cmd.Run()
		if err != nil {
			logger.Error.Println(err)
			return fmt.Errorf("error creation archive: %w", err)
		}

		// Cleanup
		err = utils.RemoveDirectory(name + ".app")
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

func CreateDirectoryHierarchy(name string) error {
	err := utils.CreateDirectory(name + ".app")
	if err != nil {
		return err
	}

	err = utils.CreateDirectory(filepath.Join(name+".app", "Contents"))
	if err != nil {
		return err
	}

	err = utils.CreateDirectory(filepath.Join(name+".app", "Contents", "MacOS"))
	if err != nil {
		return err
	}

	err = utils.CreateDirectory(filepath.Join(name+".app", "Contents", "Resources"))
	if err != nil {
		return err
	}
	return nil
}

func CreateInfoPlistFile(name, company string) error {
	// Generate Info.plist file
	tmpl, err := template.New("info").Parse(infoTmpl)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at render Info.plist file: %w", err)
	}
	file, err := os.Create(filepath.Join(name+".app", "Contents", "Info.plist"))
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at creating Info.plist file: %w", err)
	}
	defer file.Close()
	err = tmpl.Execute(file, &Info{
		Name:    name,
		Company: company,
	})
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at executing Info.plist template: %w", err)
	}
	return nil
}
