package models

import "time"

type Post struct {
	ID         int       `db:"id"`
	AuthorID   int       `db:"author_id"`
	CreateDate time.Time `db:"create_date"`
	Title      string    `db:"title"`
	Content    string    `db:"content"`
	Likes      int       `db:"likes"`
	Dislikes   int       `db:"dislikes"`
}

func (m Post) SelectString() []string {
	return []string{
		"post.id",
		"post.author_id",
		"post.create_date",
		"post.title",
		"post.content",
		"post.likes",
		"post.dislikes",
	}
}

func (m Post) TableString() string {
	return "post"
}
