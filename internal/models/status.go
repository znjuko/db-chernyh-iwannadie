package models


type Status struct {
	Forum int `json:"forum,omitempty"`
	Post int64 `json:"post,omitempty"`
	Thread int `json:"thread,omitempty"`
	User int `json:"user,omitempty"`
}
