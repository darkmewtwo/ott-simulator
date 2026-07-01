package ffmpeg

type Service struct {
	moviesDir   string
	hlsDir      string
	ffprobePath string
	ffmpegPath  string
}

func New(moviesDir string, hlsDir string, ffprobePath string, ffmpegPath string) *Service {
	return &Service{
		moviesDir:   moviesDir,
		hlsDir:      hlsDir,
		ffprobePath: ffprobePath,
		ffmpegPath:  ffmpegPath,
	}
}
