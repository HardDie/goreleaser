package common

import (
	"fmt"

	"github.com/HardDie/goreleaser/internal/logger"
	"github.com/HardDie/goreleaser/internal/utils"
)

func ValidateBinaries() error {
	if !utils.IsBinaryExist("convert") {
		return fmt.Errorf("convert binary is not exist, please install")
	}
	logger.Debug.Println("convert: exist")

	if !utils.IsBinaryExist("zip") {
		return fmt.Errorf("zip binary is not exist, please install")
	}
	logger.Debug.Println("zip: exist")

	if !utils.IsBinaryExist("tar") {
		return fmt.Errorf("tar binary is not exist, please install")
	}
	logger.Debug.Println("tar: exist")

	return nil
}
