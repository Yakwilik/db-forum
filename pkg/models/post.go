package models

type PostUpdate struct {
	Message string `json:"message"`
}

type Post struct {
	Id     int64  `json:"id"`
	Parent int64  `json:"parent"`
	Author string `json:"author"`
	PostUpdate
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int32  `json:"thread"`
	Created  string `json:"created"`
}

type Posts []Post

type PostFull struct {
	Post   `json:"post"`
	User   `json:"author"`
	Thread `json:"thread"`
	Forum  `json:"forum"`
}
