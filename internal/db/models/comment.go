package models

import "time"

type Comment struct {
	ID         int `db:"id"`
	UserID     int `db:"user_id"`
	User       *User
	PostID     int `db:"post_id"`
	Post       *Post
	CreateDate time.Time `db:"create_date"`
	Content    string    `db:"content"`
	Likes      int       `db:"likes"`
	Dislikes   int       `db:"dislikes"`
}

func (m Comment) SelectString() []string {
	return []string{
		"comment.id",
		"comment.user_id",
		"comment.post_id",
		"comment.create_date",
		"comment.content",
		"comment.likes",
		"comment.dislikes",
	}
}

func (m Comment) TableString() string {
	return "comment"
}
