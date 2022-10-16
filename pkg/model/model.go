package model

type User struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Token    string `json:"token"`
	Password string `json:"-"`
}

func (u User) IsPasswordValid(p string) bool {
	//check hash
	return u.Password == p
}
