package repos

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kctjohnson/mid-blog/internal/db"
	"github.com/kctjohnson/mid-blog/internal/db/models"
)

type PostRepository struct {
	db *db.DB
}

func NewPostRepository(db *db.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// Create post
type PostInsertParameters struct {
	BloggerID int    `db:"blogger_id"`
	Title     string `db:"title"`
	Content   string `db:"content"`
}

func (r PostRepository) Insert(newPost PostInsertParameters) (*models.Post, error) {
	query, args, err := sq.
		Insert(models.Post{}.TableString()).
		Columns("create_date", "blogger_id", "title", "content", "likes", "dislikes").
		Values(time.Now(), newPost.BloggerID, newPost.Title, newPost.Content, 0, 0).
		ToSql()
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	lastInserted, err := tx.LastInsertId()
	if err != nil {
		return nil, err
	}

	createdModel, err := r.FindByID(int(lastInserted))
	if err != nil {
		return nil, err
	}

	return createdModel, nil
}

// Get post by id
func (r PostRepository) FindByID(id int) (*models.Post, error) {
	query, args, err := sq.
		Select(models.Post{}.SelectString()...).
		From(models.Post{}.TableString()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	var found models.Post
	err = r.db.Get(&found, query, args...)
	if err != nil {
		return nil, err
	}

	return &found, nil
}

// Get all posts
func (r PostRepository) All() ([]models.Post, error) {
	query, args, err := sq.
		Select(models.Post{}.SelectString()...).
		From(models.Post{}.TableString()).
		OrderBy("(likes - dislikes) DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var found []models.Post
	err = r.db.Select(&found, query, args...)
	if err != nil {
		return nil, err
	}

	return found, nil
}

// Get comments on post
func (r PostRepository) Comments(id int) ([]models.Comment, error) {
	query, args, err := sq.
		Select(models.Comment{}.SelectString()...).
		From(models.Comment{}.TableString()).
		Join("post ON comment.post_id = post.id").
		OrderBy("create_date DESC").
		Where(sq.Eq{"post_id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	var found []models.Comment
	err = r.db.Select(&found, query, args...)
	if err != nil {
		return nil, err
	}

	return found, nil
}

type PostUpdateParameters struct {
	ID       int     `db:"id"`
	Title    *string `db:"title"`
	Content  *string `db:"content"`
	Likes    *int    `db:"likes"`
	Dislikes *int    `db:"dislikes"`
}

func (r PostRepository) Update(updatePost PostUpdateParameters) (*models.Post, error) {
	update := sq.Update(models.Post{}.TableString())

	if updatePost.Title != nil {
		update = update.Set("title", *updatePost.Title)
	}

	if updatePost.Content != nil {
		update = update.Set("content", *updatePost.Content)
	}

	if updatePost.Likes != nil {
		update = update.Set("likes", *updatePost.Likes)
	}

	if updatePost.Dislikes != nil {
		update = update.Set("dislikes", *updatePost.Dislikes)
	}

	query, args, err := update.Where(sq.Eq{"id": updatePost.ID}).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	found, err := r.FindByID(updatePost.ID)
	if err != nil {
		return nil, err
	}

	return found, nil
}

func (r PostRepository) Delete(id int) error {
	query, args, err := sq.
		Delete(models.Post{}.TableString()).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
