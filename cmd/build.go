package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/HardDie/goreleaser/internal/common"
	"github.com/HardDie/goreleaser/internal/logger"
	"github.com/HardDie/goreleaser/internal/utils"
	"github.com/HardDie/goreleaser/internal/windows"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Start to build binaries",
	Long:  `Start to build binaries`,
	Run: func(cmd *cobra.Command, args []string) {
		flagPath := cmd.Flag("path").Value.String()
		flagName := cmd.Flag("name").Value.String()
		flagImage := cmd.Flag("image").Value.String()
		flagVersion := cmd.Flag("version").Value.String()
		flagLicense := cmd.Flag("license").Value.String()

		// Valiate flags
		if !utils.IsFileExist(flagPath) {
			log.Fatalf("File %s not exist", flagPath)
		}
		if flagName == "" {
			log.Fatal("Name can't be empty")
		}

		// Validate required binaries
		err := common.ValidateBinaries()
		if err != nil {
			logger.Error.Fatal(err)
		}

		// Prepare folder for generated files
		if isExist, _ := utils.IsDirectoryExist("build_cache"); !isExist {
			err = utils.CreateDirectory("build_cache")
			if err != nil {
				logger.Error.Fatal(err)
			}
		}
		if isExist, _ := utils.IsDirectoryExist("release"); !isExist {
			err = utils.CreateDirectory("release")
			if err != nil {
				logger.Error.Fatal(err)
			}
		}

		err = windows.Build(flagName, flagImage, flagVersion, flagLicense, flagPath)
		if err != nil {
			logger.Error.Fatal(err)
		}

		// Cleanup
		err = utils.RemoveDirectory("build_cache")
		if err != nil {
			logger.Error.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("path", "p", "main.go", "Path fo file with entry point (main function)")
	buildCmd.Flags().StringP("name", "n", "", "Name of the result application")
	buildCmd.Flags().StringP("image", "i", "", "Path to the application image")
	buildCmd.Flags().StringP("version", "v", "v0.0.0", "Version of the result application")
	buildCmd.Flags().StringP("license", "l", "", "License of the result application")
}
