package generator

import (
	"math"
	"math/rand/v2"

	"simulator/internal/user"
)

type PersonalityGenerator struct {
	rng *rand.Rand
}

func newPersonalityGenerator(rng *rand.Rand) *PersonalityGenerator {
	return &PersonalityGenerator{rng: rng}
}

func (g *PersonalityGenerator) generate() user.Personality {
	return user.Personality{
		Patience:      g.randomTrait(0.60, 0.20),
		Curiosity:     g.randomTrait(0.50, 0.25),
		AttentionSpan: g.randomTrait(0.60, 0.20),
		Exploration:   g.randomTrait(0.50, 0.25),
		BingeWatching: g.randomTrait(0.50, 0.30),
		Consistency:   g.randomTrait(0.60, 0.20),
		Impulsiveness: g.randomTrait(0.40, 0.20),
	}
}

// Generates a value around the given mean.
// variation = maximum deviation from the mean.
func (g *PersonalityGenerator) randomTrait(mean, variation float64) float64 {
	value := mean + (g.rng.Float64()*2-1)*variation

	value = math.Max(0.0, math.Min(1.0, value))

	return math.Round(value*100) / 100
}
