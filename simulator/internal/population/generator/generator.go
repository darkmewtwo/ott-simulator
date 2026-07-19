package generator

import (
	"math/rand/v2"
	"simulator/internal/user"
)

type UserGenerator struct {
	IdentityGenerator *IdentityGenerator
}

func NewUserGenerator(rng *rand.Rand) (*UserGenerator, error) {
	identityGenerator, err := newIdentityGenerator(rng)

	if err != nil {
		return nil, err
	}

	return &UserGenerator{
		IdentityGenerator: identityGenerator,
	}, nil
}

func (g *UserGenerator) Generate() user.User {
	userIdentity := g.IdentityGenerator.generateIdentity()
	return user.User{Identity: userIdentity}
}
