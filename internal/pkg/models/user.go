package models

type User struct {
	ID       uint64 `json:"-"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
	About    string `json:"about"`
}

func (u *User) SetNickname(nickname string) {
	u.Nickname = nickname
}
