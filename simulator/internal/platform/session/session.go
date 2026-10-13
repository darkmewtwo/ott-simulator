package session

import (
	"log"
	"math/rand/v2"
	"simulator/internal/cognition"
	"simulator/internal/platform/capability/authentication"
	"simulator/internal/platform/capability/catalog"
	"simulator/internal/platform/httpclient"
	"simulator/internal/user"
)

type Session struct {
	User           *user.User
	RegisterFirst  bool
	HTTPClient     *httpclient.Client
	authentication *authentication.Capability
	catalog        *catalog.Capability
}

func NewSession(user *user.User, registerFirst bool) *Session {
	client := httpclient.New("http://api:8100")
	return &Session{
		User:           user,
		RegisterFirst:  registerFirst,
		HTTPClient:     client,
		authentication: authentication.New(),
		catalog:        catalog.New(),
	}
}

func (s *Session) Run() error {
	log.Println("SESSION: ", s.RegisterFirst, s.User.Identity.Username)
	log.Println("USER STATE", s.User.MentalState)
	if s.RegisterFirst {
		if _, err := s.authentication.Register(s.HTTPClient, authentication.RegisterRequest{
			Username: s.User.Identity.Username,
			Password: s.User.Identity.Password,
			Email:    s.User.Identity.Email,
		}); err != nil {
			return err
		}
	}
	if _, err := s.authentication.Login(s.HTTPClient, authentication.LoginRequest{
		Username: s.User.Identity.Username,
		Password: s.User.Identity.Password,
	}); err != nil {
		return err
	}

	defer func() {
		if err := s.authentication.Logout(s.HTTPClient); err != nil {
			log.Printf("logout failed: %v", err)
		}
	}()

	resp, err := s.catalog.ListMovies(s.HTTPClient)
	if err != nil {
		return err
	}
	log.Printf("LIST MOVIES: %+v", resp)
	movieIDs := make([]int, 0, len(resp))
	for _, movie := range resp {
		movieIDs = append(movieIDs, movie.ID)
	}
	var movieID int
	if len(movieIDs) > 0 {
		movieID = movieIDs[rand.IntN(len(movieIDs))]
	}
	mapper := NewCatalogActionMapper(s)
	for {
		scores := cognition.CalculateCatalogScores(cognition.CatalogContext{
			MentalState: s.User.MentalState,
			Personality: s.User.Personality,
		})
		best := scores[0]
		for _, score := range scores[1:] {
			if score.Score > best.Score {
				best = score
			}
		}
		if best.Action == cognition.ActionLeaveSession {
			break
		}
		switch best.Action {
		case cognition.ActionContinueBrowsing:
			movieIDs = mapper.Execute(best.Action)
			if len(movieIDs) > 0 {
				movieID = movieIDs[rand.IntN(len(movieIDs))]
			}
		case cognition.ActionSelectMovie:
			if len(movieIDs) == 0 {
				log.Printf("ending session: no movies available to select")
				return nil
			}
			mapper.Execute(best.Action, movieID)
		}
	}

	return nil
}
