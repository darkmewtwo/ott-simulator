package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func (s *Service) GenerateHlsPlaylist(
	movieID int64,
	fileName string,
) (string, error) {
	inputPath := filepath.Join(
		s.moviesDir,
		fileName,
	)

	outputDir := filepath.Join(
		s.hlsDir,
		strconv.FormatInt(movieID, 10),
	)

	playlistPath := filepath.Join(
		outputDir,
		"index.m3u8",
	)

	relativePlaylistPath := filepath.Join(
		strconv.FormatInt(movieID, 10),
		"index.m3u8",
	)

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create HLS directory: %w",
			err,
		)
	}

	cmd := exec.Command(
		s.ffmpegPath,
		"-i", inputPath,
		"-codec:v", "libx264",
		"-codec:a", "aac",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		playlistPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"ffmpeg failed: %w\n\n\n%s",
			err,
			string(output),
		)
	}

	return relativePlaylistPath, nil
}
