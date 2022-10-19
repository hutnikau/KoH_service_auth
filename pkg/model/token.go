package model

type Token struct {
	Token     string `json:"token"`
	UserId    string `json:"userId"`
	CreatedAt int64  `json:"createdAt"`
}
