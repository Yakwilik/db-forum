//go:generate easyjson -all vote.go
package models

type VoteValue int32

const (
	Up   VoteValue = 1
	Down VoteValue = -1
)

type Vote struct {
	Nickname string    `json:"nickname"`
	Voice    VoteValue `json:"voice"`
}
