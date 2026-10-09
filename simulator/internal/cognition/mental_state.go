package cognition

import (
	"math"

	"simulator/internal/user"
)

func ApplyMentalStateCost(
	state *user.MentalState,
	cost MentalStateCost,
) {
	state.Energy = applyRate(state.Energy, cost.Energy)
	state.Happiness = applyRate(state.Happiness, cost.Happiness)
	state.Boredom = applyRate(state.Boredom, cost.Boredom)
	state.Frustration = applyRate(state.Frustration, cost.Frustration)
	state.Immersion = applyRate(state.Immersion, cost.Immersion)
	state.Satisfaction = applyRate(state.Satisfaction, cost.Satisfaction)
}

func applyRate(value, rate float64) float64 {
	value = clamp01(value)

	if rate > 0 {
		return clamp01(value + rate*(1-value))
	}

	return clamp01(value + rate*value)
}

func clamp01(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}
