package repos

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kctjohnson/mid-blog/internal/db"
	"github.com/kctjohnson/mid-blog/internal/db/models"
)

type UserRepository struct {
	db *db.DB
}

func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CheckAvailable(username string, email string) error {
	_, err := r.FindByUsername(username)
	if err == nil {
		return fmt.Errorf("Username %s already exists", username)
	} else if err != sql.ErrNoRows {
		return err
	}

	_, err = r.FindByEmail(email)
	if err == nil {
		return fmt.Errorf("Email %s already exists", email)
	} else if err != sql.ErrNoRows {
		return err
	}
	return nil
}

// Create user
type UserInsertParameters struct {
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

func (r UserRepository) Insert(newUser UserInsertParameters) (*models.User, error) {
	if err := r.CheckAvailable(newUser.Username, newUser.Email); err != nil {
		return nil, err
	}

	query, args, err := sq.
		Insert(models.User{}.TableString()).
		Columns("create_date", "username", "password", "email").
		Values(time.Now(), newUser.Username, newUser.Password, newUser.Email).
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

// Get user by id
func (r UserRepository) FindByID(id int) (*models.User, error) {
	query, args, err := sq.
		Select(models.User{}.SelectString()...).
		From(models.User{}.TableString()).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.Get(&user, query, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) FindByEmail(email string) (*models.User, error) {
	query, args, err := sq.
		Select(models.User{}.SelectString()...).
		From(models.User{}.TableString()).
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.Get(&user, query, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) FindByUsername(username string) (*models.User, error) {
	query, args, err := sq.
		Select(models.User{}.SelectString()...).
		From(models.User{}.TableString()).
		Where(sq.Eq{"username": username}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.Get(&user, query, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) All() ([]models.User, error) {
	query, args, err := sq.
		Select(models.User{}.SelectString()...).
		From(models.User{}.TableString()).
		ToSql()
	if err != nil {
		return nil, err
	}

	var found []models.User
	err = r.db.Select(&found, query, args...)
	if err != nil {
		return nil, err
	}

	return found, nil
}

func (r UserRepository) Comments(id int) ([]models.Comment, error) {
	query, args, err := sq.
		Select(models.Comment{}.SelectString()...).
		From(models.Comment{}.TableString()).
		Join("user ON comment.user_id = user.id").
		OrderBy("create_date DESC").
		Where(sq.Eq{"user_id": id}).ToSql()
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

func (r UserRepository) Delete(id int) error {
	query, args, err := sq.
		Delete(models.User{}.TableString()).
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
