package orchastrator

import (
	"fmt"
	"log"
	"simulator/internal/platform/session"
	"simulator/internal/population"
)

type Orchastrator struct {
	populationManager *population.Manager
}

func NewOrchastrator(seed1, seed2 uint64) (*Orchastrator, error) {
	populationManager, err := population.NewManager(seed1, seed2)

	if err != nil {
		return nil, err
	}

	return &Orchastrator{
		populationManager: populationManager,
	}, nil
}

func (o *Orchastrator) Run() error {
	// users, err := o.userManager.Users(10)
	// log.Println("here")
	// log.Println(users)
	// if err != nil {
	// 	return err
	// }
	// for i, user := range users {
	// 	fmt.Printf(
	// 		"[%02d] %-20s %-20s %-25s %-30s %s\n",
	// 		i+1,
	// 		user.Identity.FirstName,
	// 		user.Identity.LastName,
	// 		user.Identity.Username,
	// 		user.Identity.Password,
	// 		user.Identity.Email,
	// 	)
	// }
	if err := o.populationManager.Start(); err != nil {
		return err
	}
	for injection := range o.populationManager.Users() {
		user := injection.User
		fmt.Println("recieved USER:")
		fmt.Printf(
			"%-20s %-20s %-25s %-30s %s\n",
			user.Identity.FirstName,
			user.Identity.LastName,
			user.Identity.Username,
			user.Identity.Password,
			user.Identity.Email,
		)
		// time.Sleep(5 * time.Second)

		// o.populationManager.ReturnedUsers() <- user

		go func(injection *population.Injection) {
			session := session.NewSession(
				injection.User,
				injection.Type == population.NewUser,
			)

			if err := session.Run(); err != nil {
				log.Printf("session failed for %s: %v", injection.User.Identity.Username, err)
			}

			o.populationManager.ReturnedUsers() <- injection.User
		}(injection)
	}
	return nil
}
