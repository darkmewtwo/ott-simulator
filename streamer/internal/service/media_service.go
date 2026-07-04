package service

import (
	"os"
	"path/filepath"
	"strconv"
)

type MediaService struct {
	movieDir  string
	posterDir string
	hlsDir    string
}

func NewMediaService(movieDir string, posterDir string, hlsDir string) *MediaService {
	return &MediaService{
		movieDir:  movieDir,
		posterDir: posterDir,
		hlsDir:    hlsDir,
	}
}

func (s *MediaService) MoviePath(filename string) (string, error) {
	fullPath := filepath.Join(s.movieDir, filename)
	if _, err := os.Stat(fullPath); err != nil {
		return "", err
	}
	return fullPath, nil
}

func (s *MediaService) PosterPath(filename string) (string, error) {
	fullPath := filepath.Join(s.posterDir, filename)
	if _, err := os.Stat(fullPath); err != nil {
		return "", err
	}
	return fullPath, nil
}

func (s *MediaService) HLSPath(movieID int64, fileName string) (string, error) {
	movieDir := filepath.Join(
		s.hlsDir,
		strconv.FormatInt(int64(movieID), 10),
	)

	fullPath := filepath.Join(
		movieDir,
		fileName,
	)

	if _, err := os.Stat(fullPath); err != nil {
		return "", err
	}

	return fullPath, nil
}
