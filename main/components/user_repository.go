package main

import (
	"errors"
	"sync"
)

type InMemoryUserStorage struct {
	lock    sync.RWMutex
	storage map[string]User
}

func NewInMemoryUserStorage() *InMemoryUserStorage {
	return &InMemoryUserStorage{
		lock:    sync.RWMutex{},
		storage: make(map[string]User),
	}
}

func (data *InMemoryUserStorage) Add(login string, user User) error {
	data.lock.Lock()
	_, ok := data.storage[login]
	data.lock.Unlock()
	if ok {
		return errors.New("User exists")
	}
	data.lock.Lock()
	data.storage[login] = user
	data.lock.Unlock()

	return nil
}

func (data *InMemoryUserStorage) Update(login string, user User) error {

	data.lock.Lock()
	_, currentUser := data.storage[login]
	data.lock.Unlock()

	if !currentUser {
		return errors.New("there is no user")
	}

	data.lock.Lock()
	data.storage[login] = user
	data.lock.Unlock()

	return nil
}

func (data *InMemoryUserStorage) Delete(login string) (User, error) {
	data.lock.Lock()
	currentUser, ok := data.storage[login]
	data.lock.Unlock()

	if !ok {
		return User{}, errors.New("there is no this user")
	}

	data.lock.Lock()
	delete(data.storage, login)
	data.lock.Unlock()

	return currentUser, nil
}

func (data *InMemoryUserStorage) Get(login string) (User, error) {
	data.lock.Lock()
	usr, ok := data.storage[login]
	data.lock.Unlock()

	if !ok {
		return User{}, errors.New("there is no user")
	}

	return usr, nil
}
