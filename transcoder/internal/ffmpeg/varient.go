package ffmpeg

type Variant struct {
	Name      string
	Width     int
	Height    int
	VideoRate string
	AudioRate string
	Bandwidth int
}

var variants = []Variant{
	{
		Name:      "1080p",
		Width:     1920,
		Height:    1080,
		VideoRate: "5000k",
		AudioRate: "192k",
		Bandwidth: 5192000,
	},
	{
		Name:      "720p",
		Width:     1280,
		Height:    720,
		VideoRate: "2800k",
		AudioRate: "128k",
		Bandwidth: 2928000,
	},
	{
		Name:      "480p",
		Width:     854,
		Height:    480,
		VideoRate: "1400k",
		AudioRate: "128k",
		Bandwidth: 1528000,
	},
	{
		Name:      "360p",
		Width:     640,
		Height:    360,
		VideoRate: "800k",
		AudioRate: "96k",
		Bandwidth: 896000,
	},
}
