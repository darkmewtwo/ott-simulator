package catalog

import (
	"fmt"

	"simulator/internal/platform/httpclient"
)

type MovieBaseResponse struct {
	ID                int      `json:"id"`
	Title             string   `json:"title"`
	Description       *string  `json:"description"`
	Filename          string   `json:"filename"`
	CreatedAt         string   `json:"created_at"`
	PosterFilename    *string  `json:"poster_filename"`
	TranscodingStatus string   `json:"transcoding_status"`
	ReleaseDate       *string  `json:"release_date"`
	Language          *string  `json:"language"`
	Genres            []string `json:"genres"`
	AgeRating         *string  `json:"age_rating"`
	Director          *string  `json:"director"`
	Cast              []string `json:"cast"`
	DurationSeconds   int      `json:"duration_seconds"`
}

type MovieResponse struct {
	MovieBaseResponse
	PosterURL *string `json:"poster_url"`
}

type MovieDetailsResponse struct {
	MovieBaseResponse
	StreamURL string  `json:"stream_url"`
	HLSURL    *string `json:"hls_url"`
}

func (c *Capability) ListMovies(
	client *httpclient.Client,
) ([]MovieResponse, error) {

	var response []MovieResponse

	err := client.Get("/movies", &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Capability) GetMovie(
	client *httpclient.Client,
	movieID int,
) (*MovieDetailsResponse, error) {

	var response MovieDetailsResponse

	path := fmt.Sprintf("/movies/%d", movieID)

	err := client.Get(path, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
