package generator

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"simulator/internal/user"
	"strings"
)

type IdentityGenerator struct {
	firstNames   []string
	lastNames    []string
	emailDomains []string

	// Could be removed in the future when we introduce db
	// for now it is added because it is stateless
	generatedUsernames map[string]struct{}

	rng *rand.Rand
}

func loadLines(path string) ([]string, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func newIdentityGenerator(rng *rand.Rand) (*IdentityGenerator, error) {
	g := &IdentityGenerator{
		generatedUsernames: make(map[string]struct{}),
		rng:                rng,
	}

	var err error

	if g.firstNames, err = loadLines("internal/population/simulation_data/first_names.txt"); err != nil {
		return nil, err
	}

	if g.lastNames, err = loadLines("internal/population/simulation_data/last_names.txt"); err != nil {
		return nil, err
	}

	if g.emailDomains, err = loadLines("internal/population/simulation_data/email_domains.txt"); err != nil {
		return nil, err
	}

	return g, nil
}

func (g *IdentityGenerator) generateIdentity() user.Identity {
	first := g.randomFirstName()
	last := g.randomLastName()

	username := g.generateUsername(first, last)

	return user.Identity{
		FirstName: first,
		LastName:  last,
		Username:  username,
		Email:     g.generateEmail(username),
		Password:  g.generatePassword(first, last),
	}
}

func (g *IdentityGenerator) randomFirstName() string {
	return g.firstNames[g.rng.IntN(len(g.firstNames))]
}

func (g *IdentityGenerator) randomLastName() string {
	return g.lastNames[g.rng.IntN(len(g.lastNames))]
}

func (g *IdentityGenerator) randomEmailDomain() string {
	return g.emailDomains[g.rng.IntN(len(g.emailDomains))]
}

func (g *IdentityGenerator) generateUsername(first, last string) string {
	first = strings.ToLower(first)
	last = strings.ToLower(last)

	patterns := []string{
		first + last,
		first + "." + last,
		first + "_" + last,
		string(first[0]) + last,
		last + first,
	}

	username := patterns[g.rng.IntN(len(patterns))]

	if _, exists := g.generatedUsernames[username]; exists {
		username = fmt.Sprintf("%s%d", username, g.rng.IntN(10000))
	}

	g.generatedUsernames[username] = struct{}{}

	return username
}

func (g *IdentityGenerator) generateEmail(username string) string {
	return fmt.Sprintf("%s@%s", username, g.randomEmailDomain())
}

func (g *IdentityGenerator) generatePassword(first, last string) string {
	patterns := []string{
		fmt.Sprintf("%s@123", first),
		fmt.Sprintf("%s@2025", first),
		fmt.Sprintf("%s%s@1", first, last),
		fmt.Sprintf("%s#%d", last, g.rng.IntN(1000)),
		fmt.Sprintf("%s%d!", first, g.rng.IntN(10000)),
	}

	return patterns[g.rng.IntN(len(patterns))]
}
