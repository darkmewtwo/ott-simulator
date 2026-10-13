package session

import (
	"log"
	"simulator/internal/cognition"
)

type CatalogActionHandler func(movieIDs ...int) []int

type CatalogActionMapper struct {
	handlers map[cognition.CatalogAction]CatalogActionHandler
}

func NewCatalogActionMapper(s *Session) *CatalogActionMapper {
	return &CatalogActionMapper{
		handlers: map[cognition.CatalogAction]CatalogActionHandler{
			cognition.ActionSelectMovie: func(movieIDs ...int) []int {
				if len(movieIDs) != 1 {
					log.Printf("SELECT_MOVIE requires exactly one movie ID")
					return nil
				}
				s.selectMovie(movieIDs[0])
				return nil
			},
			cognition.ActionContinueBrowsing: func(_ ...int) []int {
				return s.continueBrowsing()
			},
		},
	}
}

func (m *CatalogActionMapper) Register(
	action cognition.CatalogAction,
	handler CatalogActionHandler,
) {
	m.handlers[action] = handler
}

func (m *CatalogActionMapper) Execute(action cognition.CatalogAction, movieIDs ...int) []int {
	if handler := m.handlers[action]; handler != nil {
		return handler(movieIDs...)
	}
	return nil
}

func (s *Session) selectMovie(movieID int) {
	resp, err := s.catalog.GetMovie(s.HTTPClient, movieID)
	log.Printf("GET MOVIE: %+v, error: %v", resp, err)
	cognition.ApplyMentalStateCost(
		&s.User.MentalState,
		cognition.CostForAction(cognition.ActionSelectMovie),
	)
}

func (s *Session) continueBrowsing() []int {
	resp, err := s.catalog.ListMovies(s.HTTPClient)
	log.Printf("LIST MOVIES: %+v, error: %v", resp, err)
	cognition.ApplyMentalStateCost(
		&s.User.MentalState,
		cognition.CostForAction(cognition.ActionContinueBrowsing),
	)
	movieIDs := make([]int, 0, len(resp))
	for _, movie := range resp {
		movieIDs = append(movieIDs, movie.ID)
	}
	return movieIDs
}
