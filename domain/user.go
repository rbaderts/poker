package domain

import (
	"context"
	"database/sql"
	"time"
)

type UserId int

type User struct {
	Id            int    `json:"id"`
	AccountId     int    `json:"accountId" db:"account_id"`
	Subject       string `json:"subject"`
	Email         string `json:"email"`
	Provider      string `json:"provider"`
	PictureUrl    sql.NullString `json:"pictureUrl" db:"picture_url"`
	GivenName     string `json:"givenName" db:"given_name"`
	PreferredHandle string `json:"preferredHandle" db:"preferred_handle"`
	LastLogin     time.Time `json:"lastLogin" db:"last_login"`
	Stack         int        `json:"Stack"  db:"Stack"`

}

type UserRepository interface {
	AddUserFromProfile(ctx context.Context, accountId int, profile map[string]interface{}) (*User, error)
    AddProvidedUser(ctx context.Context, accountId int, email string, provider string, subject string, givenName string) (*User, error)
	AddUser(ctx context.Context, id int, name string) (*User, error)
	LoadUserByEmail(ctx context.Context, email string) (*User, error)
	LoadUserById(ctx context.Context, uid int) (*User, error)
	LoadUserBySubject(ctx context.Context, subject string) (*User, error)
	AddMoney(ctx context.Context, id int, amount int) error
	DeductMoney(ctx context.Context, id int, amount int) error
	UpdateUser(ctx context.Context, user *User) error
}



