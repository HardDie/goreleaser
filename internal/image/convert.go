package image

import (
	"fmt"
	"os/exec"

	"github.com/HardDie/goreleaser/internal/logger"
)

func ConvertToWindowsIcon(srcImg, dstImg string) error {
	cmd := exec.Command("convert")
	cmd.Args = append(cmd.Args, srcImg)
	cmd.Args = append(cmd.Args, "-filter", "point")
	cmd.Args = append(cmd.Args, "-resize", "256x256")
	cmd.Args = append(cmd.Args, "-define", "icon:auto-resize=256,128,96,64,48,32,16")
	cmd.Args = append(cmd.Args, dstImg)
	logger.Debug.Println("Execute:", cmd.String())
	err := cmd.Run()
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at convert image process: %w", err)
	}
	return nil
}
