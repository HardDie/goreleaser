package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/HardDie/goreleaser/internal/logger"
)

const (
	DirPerm = 0755
)

func IsBinaryExist(name string) bool {
	path, err := exec.LookPath(name)
	if err != nil {
		logger.Error.Println(err.Error())
		return false
	}
	if path == "" {
		return false
	}
	return true
}
func IsDirectoryExist(path string) (isExist bool, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// folder not exist
			return false, nil
		}

		// other error
		logger.Error.Println(err)
		return false, fmt.Errorf("error get stats of file: %w", err)
	}

	// check if it is a folder
	if !stat.IsDir() {
		return false, fmt.Errorf("there should be a folder, but it's file")
	}

	// folder exists
	return true, nil
}
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
func DataToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at creating file: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Error.Println("error at closing file:", err)
		}
	}()
	logger.Debug.Println("Created file:", filename)
	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at creating file: %w", err)
	}
	return nil
}

func CreateDirectory(path string) error {
	err := os.Mkdir(path, DirPerm)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at creatind directory: %w", err)
	}
	return nil
}
func RemoveDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error deleting directory %s: %w", path, err)
	}
	return nil
}
func RemoveFile(file string) error {
	err := os.Remove(file)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error deleting file %s: %w", file, err)
	}
	return nil
}

func CopyFile(srcFile, dstFile string) error {
	sourceFileStat, err := os.Stat(srcFile)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at get stats of source file %s for copy: %w", srcFile, err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", srcFile)
	}

	source, err := os.Open(srcFile)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at open source file %s for copy: %w", srcFile, err)
	}
	defer source.Close()

	destination, err := os.Create(dstFile)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at creation destination file %s for copy: %w", dstFile, err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at copy data from %s to %s: %w", srcFile, dstFile, err)
	}
	return nil
}
