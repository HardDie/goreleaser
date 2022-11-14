package windows

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/HardDie/goreleaser/internal/image"
	"github.com/HardDie/goreleaser/internal/logger"
	"github.com/HardDie/goreleaser/internal/utils"

	"github.com/HardDie/goversioninfo"
	goversioninfoCmd "github.com/HardDie/goversioninfo/cmd"
)

func Build(name, imagePath, version, license, path, company string) error {
	// Convert image to windows icon
	err := image.ConvertToWindowsIcon(imagePath, "build_cache/win_icon.ico")
	if err != nil {
		return err
	}

	newWorkDirectory := filepath.Dir(path)

	currentDirBackup, err := os.Getwd()
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error get work directory: %w", err)
	}

	err = os.Chdir(newWorkDirectory)
	if err != nil {
		return fmt.Errorf("error change directory to %s: %w", newWorkDirectory, err)
	}

	// Parse version
	major, minor, patch, build := utils.VersionToInt(version)

	// Generate file
	goversioninfoCmd.Cmd(goversioninfoCmd.Arguments{
		FlagOut:     utils.Allocate("resource.syso"),
		FlagPackage: utils.Allocate("main"),
		FlagIcon:    utils.Allocate("../../build_cache/win_icon.ico"),
		Flag64:      utils.Allocate(true),

		FlagVerMajor: &major,
		FlagVerMinor: &minor,
		FlagVerPatch: &patch,
		FlagVerBuild: &build,

		FlagProductVerMajor: &major,
		FlagProductVerMinor: &minor,
		FlagProductVerPatch: &patch,
		FlagProductVerBuild: &build,

		FlagExample:          utils.Allocate(false),
		FlagGo:               utils.Allocate(""),
		FlagPlatformSpecific: utils.Allocate(false),
		FlagManifest:         utils.Allocate(""),
		FlagSkipVersion:      utils.Allocate(true),

		FlagComment:        utils.Allocate(""),
		FlagCompany:        &company,
		FlagDescription:    utils.Allocate(""),
		FlagFileVersion:    &version,
		FlagInternalName:   &name,
		FlagCopyright:      &license,
		FlagTrademark:      utils.Allocate(""),
		FlagOriginalName:   &name,
		FlagPrivateBuild:   utils.Allocate(""),
		FlagProductName:    &name,
		FlagProductVersion: &version,
		FlagSpecialBuild:   utils.Allocate(""),

		FlagTranslation: utils.Allocate(int(goversioninfo.LngUSEnglish)),
		FlagCharset:     utils.Allocate(int(goversioninfo.CsUnicode)),

		Flagarm: utils.Allocate(true),
	})

	arches := []string{"amd64", "386", "arm64"}
	for _, arch := range arches {
		// Compile app
		cmd := exec.Command("go", "build", "-a", "-o", name+".exe", ".")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
		cmd.Env = append(cmd.Env, "GOARCH="+arch)
		cmd.Env = append(cmd.Env, "GOOS=windows")
		logger.Debug.Println("Execute:", cmd.String())
		err = cmd.Run()
		if err != nil {
			logger.Error.Println(err)
			return fmt.Errorf("error building application: %w", err)
		}

		// Create zip archive
		cmd = exec.Command("zip", "../../release/"+name+".windows-"+arch+".zip", name+".exe")
		logger.Debug.Println("Execute:", cmd.String())
		err = cmd.Run()
		if err != nil {
			logger.Error.Println(err)
			return fmt.Errorf("error creation archive: %w", err)
		}

		// Remove binary file
		err = utils.RemoveFile(name + ".exe")
		if err != nil {
			return err
		}
	}

	// Cleanup before exit
	err = utils.RemoveFile("resource.syso")
	if err != nil {
		return err
	}

	// Return to the root folder
	err = os.Chdir(currentDirBackup)
	if err != nil {
		return fmt.Errorf("error return from directory: %w", err)
	}

	return nil
}
