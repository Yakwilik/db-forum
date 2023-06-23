//go:generate easyjson -all forum.go

package models

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int32  `json:"threads,omitempty"`
}

type ServiceStatus struct {
	User   int64 `json:"user"`
	Forum  int64 `json:"forum"`
	Thread int64 `json:"thread"`
	Post   int64 `json:"post"`
}
