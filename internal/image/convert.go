package image

import (
	"fmt"
	"image"
	"os"
	"os/exec"

	"github.com/jackmordaunt/icns/v2"

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

func ConvertToDarwinIconsContainer(srcImg, dstFile string) error {
	// Open image
	imgFile, err := os.Open(srcImg)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at opening src image %q: %w", srcImg, err)
	}
	defer imgFile.Close()
	// Decode image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at decoding src image %q: %w", srcImg, err)
	}
	// Create MacOS image container
	dest, err := os.Create(dstFile)
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at creation image container %q: %w", dstFile, err)
	}
	defer dest.Close()
	// Build container
	if err = icns.Encode(dest, img); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("error at building image container: %w", err)
	}
	return nil
}
