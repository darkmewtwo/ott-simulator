package population

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand/v2"
	"os"
	"simulator/internal/population/generator"
	"simulator/internal/user"
	"sync"
	"time"
)

const usersFile = "internal/population/simulation_data/registered_users.json"

type Manager struct {
	mu            sync.RWMutex
	userGenerator *generator.UserGenerator
	// Complete registered population
	registeredUsers []*user.User

	// Users available for this cycle
	availableUsers []*user.User

	// Stream consumed by the orchestrator
	injectUserCh chan *user.User

	// Stream consumed by Manager for logged out users
	logoutUserCh chan *user.User

	dirty bool

	users []*user.User // to be removed
}

func NewManager(seed1, seed2 uint64) (*Manager, error) {

	rng := rand.New(rand.NewPCG(seed1, seed2))

	userGenerator, err := generator.NewUserGenerator(rng)

	if err != nil {
		return nil, err
	}
	m := &Manager{
		userGenerator: userGenerator,
		dirty:         false,
		injectUserCh:  make(chan *user.User),
		logoutUserCh:  make(chan *user.User),
	}

	return m, nil
}

func (m *Manager) load() error {
	file, err := os.Open(usersFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			m.registeredUsers = []*user.User{}
			return nil
		}
		return err
	}
	info, err := file.Stat()

	if err != nil {
		return err
	}

	if info.Size() == 0 {
		m.registeredUsers = []*user.User{}
		return nil
	}

	if err = json.NewDecoder(file).Decode(&m.registeredUsers); err != nil {
		return err
	}
	if m.registeredUsers == nil {
		m.registeredUsers = []*user.User{}
	}
	return nil
}

func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.dirty {
		return nil
	}

	tmpFile := usersFile + ".tmp"

	file, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(m.registeredUsers); err != nil {
		file.Close()
		os.Remove(tmpFile)
		return err
	}

	if err := file.Close(); err != nil {
		os.Remove(tmpFile)
		return err
	}

	if err := os.Rename(tmpFile, usersFile); err != nil {
		os.Remove(tmpFile)
		return err
	}

	m.dirty = false

	return nil
}

func (m *Manager) AddUser(user *user.User) {
	log.Println("NOT HERE?")
	m.mu.Lock()
	defer m.mu.Unlock()
	log.Println("not here after getting lock")
	m.users = append(m.users, user)
	log.Println("ADDED TO LIST?")
	m.dirty = true
}

// func (m *Manager) Users(count int) ([]*user.User, error) {
// 	if err := m.load(); err != nil {
// 		return nil, err
// 	}

// 	// m.mu.Lock()
// 	// defer m.mu.Unlock()

// 	log.Println("GET USERS", len(m.users))
// 	if len(m.users) == 0 {
// 		for i := 0; i < count; i++ {
// 			log.Println("IN FOR LOOP")
// 			user := m.userGenerator.Generate()
// 			log.Println(user)
// 			m.AddUser(&user)
// 		}
// 	}

// 	return m.users, nil
// }

func (m *Manager) periodicSave() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("running periodic save")
		if err := m.Save(); err != nil {
			log.Printf("failed to save users: %v", err)
		}
	}
}

func (m *Manager) shouldCreateNewUser() bool {
	// 20% chance of creating a new user.
	return rand.IntN(100) < 20
}

func (m *Manager) nextInjectionInterval() time.Duration {
	// Random interval between 5 and 30 seconds.
	seconds := 5 + rand.IntN(26)

	return time.Duration(seconds) * time.Second
}

func (m *Manager) nextCooldown() time.Duration {
	// Random cooldown between 5 and 15 seconds.
	seconds := 5 + rand.IntN(11)

	return time.Duration(seconds) * time.Second
}

func (m *Manager) createUser() *user.User {
	u := m.userGenerator.Generate()

	m.registeredUsers = append(m.registeredUsers, &u)
	m.dirty = true

	return &u
}

func (m *Manager) selectAvailableUser() *user.User {
	// lock is needed when reinserting logged out users in the future
	// m.mu.Lock()
	// defer m.mu.Unlock()

	index := rand.IntN(len(m.availableUsers))
	u := m.availableUsers[index]

	// Remove the selected user using swap-with-last.
	last := len(m.availableUsers) - 1
	m.availableUsers[index] = m.availableUsers[last]
	m.availableUsers[last] = nil // Allow GC
	m.availableUsers = m.availableUsers[:last]
	log.Println("[MANAGER] Len of available users:", len(m.availableUsers))
	return u
}

func (m *Manager) injectionLoop() {
	for {
		log.Println("[MANAGER] INjection loop started")
		time.Sleep(m.nextInjectionInterval())

		if m.shouldCreateNewUser() {
			log.Println("[MANAGER] creating new user")
			u := m.createUser()
			m.injectUserCh <- u
			continue
		}

		if len(m.availableUsers) == 0 {
			log.Println("[MANAGER] no available users")
			continue
		}
		log.Println("[MANAGER] selecting available user")
		u := m.selectAvailableUser()
		m.injectUserCh <- u
	}
}

func (m *Manager) Start() error {
	if err := m.load(); err != nil {
		return err
	}

	m.availableUsers = append([]*user.User(nil), m.registeredUsers...)

	go m.periodicSave()
	go m.injectionLoop()
	go m.logoutLoop()

	return nil
}

func (m *Manager) cooldown(u *user.User) {
	time.Sleep(m.nextCooldown())

	m.mu.Lock()
	m.availableUsers = append(m.availableUsers, u)
	log.Println("[MANAGER] added user back to available user")
	m.mu.Unlock()
}

func (m *Manager) Users() <-chan *user.User {
	return m.injectUserCh
}

func (m *Manager) ReturnedUsers() chan<- *user.User {
	return m.logoutUserCh
}

func (m *Manager) logoutLoop() {
	for u := range m.logoutUserCh {
		log.Printf("User returned: %s", u.Identity.Username)
		go m.cooldown(u)
	}
}
