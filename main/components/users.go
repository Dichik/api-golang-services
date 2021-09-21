package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"net/http"
)

type UserService struct {
	repository UserRepository
}

type UserRegisterParams struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FavoriteCake string `json:"favorite_cake"`
}

func validateRegisterParams(p *UserRegisterParams) error {
	if !(emailValid(p.Email) && len(p.Password) >= 8 && checkFavorite(p.FavoriteCake)) {
		return errors.New("error")
	}
	return nil
}

func checkFavorite(name string) bool {
	for v := range name {
		if v < 'A' || v > 'z' {
			return false
		} else if v > 'Z' && v < 'a' {
			return false
		}
	}
	return !(name == "")
}

func emailValid(email string) bool {
	for i := 1; i < len(email); i++ {
		if email[i] == '@' {
			return true
		}
	}
	return false
}

func (u *UserService) Register(w http.ResponseWriter, r *http.Request) {
	params := &UserRegisterParams{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		handleError(errors.New("could not read params"), w)
		return
	}
	if err := validateRegisterParams(params); err != nil {
		handleError(err, w)
		return
	}
	passwordDigest := md5.New().Sum([]byte(params.Password))
	newUser := User{
		Email:          params.Email,
		PasswordDigest: string(passwordDigest),
		FavoriteCake:   params.FavoriteCake,
	}
	err = u.repository.Add(params.Email, newUser)
	if err != nil {
		handleError(err, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("registered"))
}
func handleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte(err.Error()))
}

type User struct {
	Email          string
	PasswordDigest string
	FavoriteCake   string
}

type UserRepository interface {
	Add(string, User) error
	Get(string) (User, error)
	Update(string, User) error
	Delete(string) (User, error)
}
