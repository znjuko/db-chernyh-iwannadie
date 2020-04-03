package models


type Forum struct {
	Posts int64 `json:"posts,omitempty"`
	Threads int `json:"threads,omitempty"`
	Slug string `json:"slug,omitempty"`
	Title string `json:"title,omitempty"`
	User string `json:"user,omitempty"`
}
