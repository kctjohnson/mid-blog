package models

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type Blogger struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Age       int    `db:"age"`
	Gender    Gender `db:"gender"`
	Bio       string `db:"bio"`
}

func (m Blogger) SelectString() []string {
	return []string{
		"blogger.id",
		"blogger.first_name",
		"blogger.last_name",
		"blogger.email",
		"blogger.age",
		"blogger.gender",
		"blogger.bio",
	}
}

func (m Blogger) TableString() string {
	return "blogger"
}
