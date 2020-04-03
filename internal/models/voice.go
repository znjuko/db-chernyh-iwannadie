package models

type Voice struct {
	Nickname string `json:"nickname,omitempty"`
	Voice    int    `json:"voice,omitempty"`
}
