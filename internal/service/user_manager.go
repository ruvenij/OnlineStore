package service

import (
	"OnlieStore/internal/model"
	"OnlieStore/internal/util"
	"errors"
	"fmt"
	"sync"
)

type UserManager struct {
	mu              sync.RWMutex
	users           map[string]*model.User // key - user id, value - user
	usersByName     map[string]*model.User // key - user name, value - user
	latestUserIndex int
}

func NewUserManager() *UserManager {
	return &UserManager{
		users:       make(map[string]*model.User),
		usersByName: make(map[string]*model.User),
	}
}

func (um *UserManager) GetUser(id string) (*model.User, error) {
	um.mu.RLock()
	defer um.mu.RUnlock()

	u, ok := um.users[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("User not found, username: %s", id))
	}

	return u, nil
}

func (um *UserManager) AddUser(u *model.User) error {
	um.mu.Lock()
	defer um.mu.Unlock()

	u.ID = fmt.Sprintf("%s05%d", util.UserSuffix, um.latestUserIndex)

	um.users[u.ID] = u
	um.usersByName[u.Name] = u

	um.latestUserIndex++
	return nil
}

func (um *UserManager) ValidateAndGetUser(userName string, password string) (*model.User, error) {
	um.mu.RLock()
	defer um.mu.RUnlock()

	u, ok := um.usersByName[userName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("User not found, username: %s", userName))
	}

	if u.Password != password {
		return nil, errors.New(fmt.Sprintf("Invalid password, username: %s", userName))
	}

	return u, nil
}
