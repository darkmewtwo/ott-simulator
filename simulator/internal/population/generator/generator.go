package generator

import (
	"math/rand/v2"

	"simulator/internal/user"
)

type UserGenerator struct {
	IdentityGenerator    *IdentityGenerator
	PersonalityGenerator *PersonalityGenerator
}

func NewUserGenerator(rng *rand.Rand) (*UserGenerator, error) {
	identityGenerator, err := newIdentityGenerator(rng)
	if err != nil {
		return nil, err
	}
	personalityGenerator := newPersonalityGenerator(rng)

	return &UserGenerator{
		IdentityGenerator:    identityGenerator,
		PersonalityGenerator: personalityGenerator,
	}, nil
}

func (g *UserGenerator) Generate() user.User {
	userIdentity := g.IdentityGenerator.generateIdentity()
	userPersonality := g.PersonalityGenerator.generate()
	return user.User{Identity: userIdentity, Personality: userPersonality}
}
