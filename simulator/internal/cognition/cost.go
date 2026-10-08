package cognition

type MentalStateCost struct {
	Energy       float64
	Happiness    float64
	Boredom      float64
	Frustration  float64
	Immersion    float64
	Satisfaction float64
}

func CostForAction(action CatalogAction) MentalStateCost {
	switch action {
	case ActionSelectMovie:
		return MentalStateCost{
			Energy:       -0.01,
			Happiness:    0.05,
			Boredom:      -0.20,
			Frustration:  -0.10,
			Satisfaction: 0.10,
		}

	case ActionContinueBrowsing:
		return MentalStateCost{
			Energy:       -0.04,
			Happiness:    -0.01,
			Boredom:      0.06,
			Frustration:  0.04,
			Satisfaction: -0.03,
		}

	default:
		return MentalStateCost{}
	}
}
