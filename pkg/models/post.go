//go:generate easyjson -all post.go
package models

type PostUpdate struct {
	Message string `json:"message"`
}

type Post struct {
	Id     int64  `json:"id"`
	Parent int64  `json:"parent"`
	Author string `json:"author"`
	PostUpdate
	IsEdited bool   `json:"isEdited" db:"is_edited"`
	Forum    string `json:"forum"`
	Thread   int32  `json:"thread"`
	Created  string `json:"created"`
}

//easyjson:json
type Posts []Post

type PostFull struct {
	Post    `json:"post"`
	*User   `json:"author,omitempty"`
	*Thread `json:"thread,omitempty"`
	*Forum  `json:"forum,omitempty"`
}
