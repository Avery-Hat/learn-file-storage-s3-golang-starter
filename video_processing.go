package main

import (
	"fmt"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {
	outPath := filePath + ".processing"

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", filePath,
		"-c", "copy",
		"-movflags", "faststart",
		"-f", "mp4",
		outPath,
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg faststart failed: %w", err)
	}

	return outPath, nil
}
