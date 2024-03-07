package models

import "time"

type User struct {
	ID         int       `db:"id"`
	Username   string    `db:"username"`
	Password   string    `db:"password"`
	Email      string    `db:"email"`
	CreateDate time.Time `db:"create_date"`
	Avatar     int       `db:"avatar"`
}

func (m User) SelectString() []string {
	return []string{
		"user.id",
		"user.username",
		"user.password",
		"user.email",
		"user.create_date",
		"user.avatar",
	}
}

func (m User) TableString() string {
	return "user"
}
