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

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create HLS directory: %w",
			err,
		)
	}

	for _, variant := range variants {

		if err := s.generateVariant(
			inputPath,
			outputDir,
			variant,
		); err != nil {
			return "", err
		}
	}
	err = s.generateMasterPlaylist(outputDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(
		strconv.FormatInt(movieID, 10),
		"master.m3u8",
	), nil
}

func (s *Service) generateVariant(
	inputPath string,
	outputRoot string,
	variant Variant,
) error {

	variantDir := filepath.Join(
		outputRoot,
		variant.Name,
	)

	err := os.MkdirAll(
		variantDir,
		0755,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create variant directory: %w",
			err,
		)
	}

	playlistPath := filepath.Join(
		variantDir,
		"index.m3u8",
	)

	cmd := exec.Command(
		s.ffmpegPath,

		"-i", inputPath,

		"-vf",
		fmt.Sprintf(
			"scale=%d:%d",
			variant.Width,
			variant.Height,
		),

		"-c:v", "libx264",
		"-c:a", "aac",

		"-b:v", variant.VideoRate,
		"-b:a", variant.AudioRate,

		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",

		playlistPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"%s variant failed: %w\n\n%s",
			variant.Name,
			err,
			string(output),
		)
	}

	return nil
}

func (s *Service) generateMasterPlaylist(
	outputDir string,
) error {
	master := filepath.Join(
		outputDir,
		"master.m3u8",
	)

	file, err := os.Create(master)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Fprintf(file, "#EXTM3U\n\n")
	fmt.Fprintf(file, "#EXT-X-VERSION:3\n\n")
	for _, variant := range variants {

		fmt.Fprintf(
			file,
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n",
			variant.Bandwidth,
			variant.Width,
			variant.Height,
		)

		fmt.Fprintf(
			file,
			"%s/index.m3u8\n",
			variant.Name,
		)
	}
	return nil
}
