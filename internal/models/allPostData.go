package models


type AllPostData struct {
	Author *UserModel `json:"author,omitempty"`
	Post *Message `json:"post,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum *Forum `json:"forum,omitempty"`
}