package models

type ThreadUpdate struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Thread struct {
	Id int64 `json:"id"`
	ThreadUpdate
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Votes   int32  `json:"votes"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
}

type Threads []Thread
