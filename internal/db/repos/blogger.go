package repos

import (
	"time"

	"github.com/kctjohnson/mid-blog/internal/db"
	"github.com/kctjohnson/mid-blog/internal/db/models"

	sq "github.com/Masterminds/squirrel"
)

type BloggerRepository struct {
	db *db.DB
}

func NewBloggerRepository(db *db.DB) *BloggerRepository {
	return &BloggerRepository{
		db: db,
	}
}

// Create blogger
type BloggerInsertParameters struct {
	FirstName string        `db:"first_name"`
	LastName  string        `db:"last_name"`
	Email     string        `db:"email"`
	Age       int           `db:"age"`
	Gender    models.Gender `db:"gender"`
	Bio       string        `db:"bio"`
	Avatar    int           `db:"avatar"`
}

func (r BloggerRepository) Insert(newBlogger BloggerInsertParameters) (*models.Blogger, error) {
	query, args, err := sq.
		Insert(models.Blogger{}.TableString()).
		Columns("create_date", "first_name", "last_name", "email", "age", "gender", "bio", "avatar").
		Values(time.Now(), newBlogger.FirstName, newBlogger.LastName, newBlogger.Email, newBlogger.Age, newBlogger.Gender, newBlogger.Bio, newBlogger.Avatar).
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

// Get blogger by id
func (r BloggerRepository) FindByID(id int) (*models.Blogger, error) {
	query, args, err := sq.
		Select(models.Blogger{}.SelectString()...).
		From(models.Blogger{}.TableString()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	var found models.Blogger
	err = r.db.Get(&found, query, args...)
	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r BloggerRepository) All() ([]models.Blogger, error) {
	query, args, err := sq.
		Select(models.Blogger{}.SelectString()...).
		From(models.Blogger{}.TableString()).
		ToSql()
	if err != nil {
		return nil, err
	}

	var found []models.Blogger
	err = r.db.Select(&found, query, args...)
	if err != nil {
		return nil, err
	}

	return found, nil
}

type BloggerUpdateParameters struct {
	ID        int            `db:"id"`
	FirstName *string        `db:"first_name"`
	LastName  *string        `db:"last_name"`
	Email     *string        `db:"email"`
	Age       *int           `db:"age"`
	Gender    *models.Gender `db:"gender"`
	Bio       *string        `db:"bio"`
	Avatar    *int           `db:"avatar"`
}

func (r BloggerRepository) Update(updateBlogger BloggerUpdateParameters) (*models.Blogger, error) {
	update := sq.Update(models.Blogger{}.TableString())

	if updateBlogger.FirstName != nil {
		update = update.Set("first_name", *updateBlogger.FirstName)
	}

	if updateBlogger.LastName != nil {
		update = update.Set("last_name", *updateBlogger.LastName)
	}

	if updateBlogger.Email != nil {
		update = update.Set("email", *updateBlogger.Email)
	}

	if updateBlogger.Age != nil {
		update = update.Set("age", *updateBlogger.Age)
	}

	if updateBlogger.Gender != nil {
		update = update.Set("gender", *updateBlogger.Gender)
	}

	if updateBlogger.Bio != nil {
		update = update.Set("bio", *updateBlogger.Bio)
	}

	if updateBlogger.Avatar != nil {
		update = update.Set("avatar", *updateBlogger.Avatar)
	}

	query, args, err := update.Where(sq.Eq{"id": updateBlogger.ID}).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	found, err := r.FindByID(updateBlogger.ID)
	if err != nil {
		return nil, err
	}

	return found, nil
}

func (r BloggerRepository) Delete(id int) error {
	query, args, err := sq.
		Delete(models.Blogger{}.TableString()).
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
