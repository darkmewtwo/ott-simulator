package ffmpeg

import (
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *Service) GetDuration(
	fileName string,
) (int, error) {

	fullPath := filepath.Join(
		s.moviesDir,
		fileName,
	)

	cmd := exec.Command(
		s.ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		fullPath,
	)
	output, err := cmd.Output()

	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	durationStr := strings.TrimSpace(string(output))

	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid duration: %w", err)
	}

	return int(math.Round(duration)), nil
}
