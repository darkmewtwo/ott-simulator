package generator

import (
	"bufio"
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

func NewIdentityGenerator(seed int64) (*IdentityGenerator, error) {
	g := &IdentityGenerator{
		generatedUsernames: make(map[string]struct{}),
		rng:                rand.New(rand.NewSource(seed)),
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

func (g *IdentityGenerator) GenerateIdentity() user.Identity {
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

func (g *IdentityGenerator) randomFirstName() string
func (g *IdentityGenerator) randomLastName() string

func (g *IdentityGenerator) generateUsername(first, last string) string
func (g *IdentityGenerator) generateEmail(username string) string
func (g *IdentityGenerator) generatePassword(first, last string) string
