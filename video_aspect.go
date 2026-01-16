package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe failed: %w", err)
	}

	var parsed ffprobeOutput
	if err := json.Unmarshal(out.Bytes(), &parsed); err != nil {
		return "", fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	if len(parsed.Streams) == 0 {
		return "", fmt.Errorf("no streams found")
	}

	w := parsed.Streams[0].Width
	h := parsed.Streams[0].Height
	if w == 0 || h == 0 {
		return "", fmt.Errorf("invalid dimensions: %dx%d", w, h)
	}

	ratio := float64(w) / float64(h)

	// Tolerance-based classification
	if math.Abs(ratio-(16.0/9.0)) < 0.05 {
		return "16:9", nil
	}
	if math.Abs(ratio-(9.0/16.0)) < 0.05 {
		return "9:16", nil
	}
	return "other", nil
}
