package session

import (
	"fmt"

	"simulator/internal/cognition"
	"simulator/internal/platform/capability/catalog"
)

type CatalogActionResult struct {
	SelectedMovie *catalog.MovieDetailsResponse
	EndSession    bool
}

type CatalogActionHandler func(
	candidate catalog.MovieResponse,
) (CatalogActionResult, error)

type CatalogActionMapper struct {
	handlers map[cognition.CatalogAction]CatalogActionHandler
}

func NewCatalogActionMapper(s *Session) *CatalogActionMapper {
	mapper := &CatalogActionMapper{
		handlers: make(
			map[cognition.CatalogAction]CatalogActionHandler,
		),
	}

	mapper.Register(
		cognition.ActionSelectMovie,
		s.selectMovie,
	)

	mapper.Register(
		cognition.ActionContinueBrowsing,
		s.continueBrowsing,
	)

	mapper.Register(
		cognition.ActionLeaveSession,
		s.leaveSession,
	)

	return mapper
}

func (m *CatalogActionMapper) Register(
	action cognition.CatalogAction,
	handler CatalogActionHandler,
) {
	m.handlers[action] = handler
}

func (m *CatalogActionMapper) Execute(
	action cognition.CatalogAction,
	candidate catalog.MovieResponse,
) (CatalogActionResult, error) {
	handler, exists := m.handlers[action]
	if !exists || handler == nil {
		return CatalogActionResult{}, fmt.Errorf(
			"unsupported catalog action: %q",
			action,
		)
	}

	return handler(candidate)
}

func (s *Session) selectMovie(
	candidate catalog.MovieResponse,
) (CatalogActionResult, error) {
	movie, err := s.catalog.GetMovie(
		s.HTTPClient,
		candidate.ID,
	)
	if err != nil {
		return CatalogActionResult{}, err
	}

	cognition.ApplyMentalStateCost(
		&s.User.MentalState,
		cognition.CostForAction(cognition.ActionSelectMovie),
	)

	return CatalogActionResult{
		SelectedMovie: movie,
	}, nil
}

func (s *Session) continueBrowsing(
	_ catalog.MovieResponse,
) (CatalogActionResult, error) {
	cognition.ApplyMentalStateCost(
		&s.User.MentalState,
		cognition.CostForAction(
			cognition.ActionContinueBrowsing,
		),
	)

	return CatalogActionResult{}, nil
}

func (s *Session) leaveSession(
	_ catalog.MovieResponse,
) (CatalogActionResult, error) {
	return CatalogActionResult{
		EndSession: true,
	}, nil
}
