package generator

import (
	"math"
	"math/rand/v2"

	"simulator/internal/user"
)

type MentalStateGenerator struct {
	rng *rand.Rand
}

func NewMentalStateGenerator(rng *rand.Rand) *MentalStateGenerator {
	return &MentalStateGenerator{rng: rng}
}

func (g *MentalStateGenerator) Generate() user.MentalState {
	return user.MentalState{
		Energy:       g.randomValue(0.55, 0.90),
		Happiness:    g.randomValue(0.45, 0.75),
		Boredom:      g.randomValue(0.25, 0.60),
		Frustration:  0.0,
		Immersion:    0.0,
		Satisfaction: 0.5,
	}
}

func (g *MentalStateGenerator) randomValue(
	minimum float64,
	maximum float64,
) float64 {
	value := minimum + g.rng.Float64()*(maximum-minimum)
	return math.Round(value*100) / 100
}
