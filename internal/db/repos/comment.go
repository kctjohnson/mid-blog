package repos

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kctjohnson/mid-blog/internal/db"
	"github.com/kctjohnson/mid-blog/internal/db/models"
)

type CommentRepository struct {
	db *db.DB
}

func NewCommentRepository(db *db.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

// Create comment
type CommentInsertParameters struct {
	UserID  int    `db:"user_id"`
	PostID  int    `db:"post_id"`
	Content string `db:"content"`
}

func (r CommentRepository) Insert(newComment CommentInsertParameters) (*models.Comment, error) {
	query, args, err := sq.
		Insert(models.Comment{}.TableString()).
		Columns("create_date", "user_id", "post_id", "content", "likes", "dislikes").
		Values(time.Now(), newComment.UserID, newComment.PostID, newComment.Content, 0, 0).
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

// Get comment by id
func (r CommentRepository) FindByID(id int) (*models.Comment, error) {
	query, args, err := sq.
		Select(models.Comment{}.SelectString()...).
		From(models.Comment{}.TableString()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	var found models.Comment
	err = r.db.Get(&found, query, args...)
	if err != nil {
		return nil, err
	}

	return &found, nil
}

type CommentUpdateParameters struct {
	ID       int     `db:"id"`
	Content  *string `db:"content"`
	Likes    *int    `db:"likes"`
	Dislikes *int    `db:"dislikes"`
}

func (r CommentRepository) Update(updateComment CommentUpdateParameters) (*models.Comment, error) {
	update := sq.Update(models.Comment{}.TableString())

	if updateComment.Content != nil {
		update = update.Set("content", *updateComment.Content)
	}

	if updateComment.Likes != nil {
		update = update.Set("likes", *updateComment.Likes)
	}

	if updateComment.Dislikes != nil {
		update = update.Set("dislikes", *updateComment.Dislikes)
	}

	query, args, err := update.Where(sq.Eq{"id": updateComment.ID}).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	found, err := r.FindByID(updateComment.ID)
	if err != nil {
		return nil, err
	}

	return found, nil
}

func (r CommentRepository) Delete(id int) error {
	query, args, err := sq.
		Delete(models.Comment{}.TableString()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
