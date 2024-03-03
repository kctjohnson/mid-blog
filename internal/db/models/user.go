package models

import "time"

type User struct {
	ID         int       `db:"id"`
	Username   string    `db:"username"`
	Password   string    `db:"password"`
	CreateDate time.Time `db:"create_date"`
}

func (m User) SelectString() []string {
	return []string{
		"user.id",
		"user.username",
		"user.password",
		"user.create_date",
	}
}

func (m User) TableString() string {
	return "user"
}
