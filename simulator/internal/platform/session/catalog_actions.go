package session

import "simulator/internal/cognition"

type CatalogActionHandler func()

type CatalogActionMapper struct {
	handlers map[cognition.CatalogAction]CatalogActionHandler
}

func NewCatalogActionMapper(s *Session) *CatalogActionMapper {
	return &CatalogActionMapper{
		handlers: map[cognition.CatalogAction]CatalogActionHandler{
			cognition.ActionSelectMovie:      s.selectMovie,
			cognition.ActionContinueBrowsing: s.continueBrowsing,
		},
	}
}

func (m *CatalogActionMapper) Register(
	action cognition.CatalogAction,
	handler CatalogActionHandler,
) {
	m.handlers[action] = handler
}

func (m *CatalogActionMapper) Execute(action cognition.CatalogAction) {
	if handler := m.handlers[action]; handler != nil {
		handler()
	}
}

func (s *Session) selectMovie() {
	cognition.ApplyMentalStateCost(
		&s.User.MentalState,
		cognition.CostForAction(cognition.ActionSelectMovie),
	)
}

func (s *Session) continueBrowsing() {
	cognition.ApplyMentalStateCost(
		&s.User.MentalState,
		cognition.CostForAction(cognition.ActionContinueBrowsing),
	)
}
