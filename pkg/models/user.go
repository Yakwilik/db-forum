//go:generate easyjson -all user.go
package models

type UserUpdate struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type User struct {
	Nickname string `json:"nickname,omitempty"`
	UserUpdate
}

//easyjson:json
type Users []User
