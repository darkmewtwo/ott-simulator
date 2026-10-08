package cognition

import "simulator/internal/user"

type CatalogContext struct {
	Personality user.Personality
	MentalState user.MentalState
}

type ActionScore struct {
	Action CatalogAction
	Score  float64
}

type CatalogScorer func(CatalogContext) ActionScore

func CalculateCatalogScores(
	context CatalogContext,
	scorers ...CatalogScorer,
) []ActionScore {
	if len(scorers) == 0 {
		scorers = []CatalogScorer{
			ScoreSelectMovie,
			ScoreContinueBrowsing,
			ScoreLeaveSession,
		}
	}

	scores := make([]ActionScore, 0, len(scorers))

	for _, scorer := range scorers {
		scores = append(scores, scorer(context))
	}

	return scores
}

func ScoreSelectMovie(context CatalogContext) ActionScore {
	p := context.Personality
	m := context.MentalState

	positive := 0.35*p.Impulsiveness +
		0.25*m.Satisfaction +
		0.20*m.Happiness +
		0.20*(m.Boredom*p.Impulsiveness)

	negative := 0.40*m.Frustration +
		0.35*(1-m.Energy) +
		0.25*(1-p.AttentionSpan)

	return ActionScore{
		Action: ActionSelectMovie,
		Score:  positive - negative,
	}
}

func ScoreContinueBrowsing(context CatalogContext) ActionScore {
	p := context.Personality
	m := context.MentalState

	positive := 0.30*p.Curiosity +
		0.25*p.Exploration +
		0.25*p.Patience +
		0.20*m.Energy

	negative := 0.40*m.Frustration +
		0.25*m.Boredom +
		0.20*(1-p.AttentionSpan) +
		0.15*(1-m.Satisfaction)

	return ActionScore{
		Action: ActionContinueBrowsing,
		Score:  positive - negative,
	}
}

func ScoreLeaveSession(context CatalogContext) ActionScore {
	p := context.Personality
	m := context.MentalState

	positive := 0.35*m.Frustration +
		0.25*(1-m.Energy) +
		0.20*(1-m.Satisfaction) +
		0.20*(m.Boredom*m.Frustration)

	negative := 0.35*p.Patience +
		0.30*p.Curiosity +
		0.20*m.Energy +
		0.15*m.Happiness

	return ActionScore{
		Action: ActionLeaveSession,
		Score:  positive - negative,
	}
}
