package models

type UserModel struct {
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Email string `json:"email,omitempty"`
	About string `json:"about,omitempty"`
}
