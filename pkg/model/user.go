package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Token    string `json:"token"`
	Password string `json:"password,omitempty"`
}

func (u User) IsPasswordValid(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	return err == nil
}
