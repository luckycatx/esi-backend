package user

import (
	"context"
	"database/sql"

	"esi/internal/pkg/db/sqlc"
)

type Repo struct {
	// Perform db operations using generated query struct
	// instead of using the db interface directly
	qry sqlc.Querier
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{qry: sqlc.New(db)}
}

func (ur *Repo) create(ctx context.Context, u *User) error {
	if err := ur.qry.CreateUser(ctx, &sqlc.CreateUserParams{
		Username: u.Username,
		Email:    u.Email,
		Pwd:      u.Pwd,
	}); err != nil {
		return err
	}
	return nil
}

func (ur *Repo) fetch(ctx context.Context) ([]*User, error) {
	users, err := ur.qry.FetchUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (ur *Repo) getByID(ctx context.Context, id UUID) (*User, error) {
	user, err := ur.qry.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *Repo) getByEmail(ctx context.Context, email string) (*User, error) {
	user, err := ur.qry.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *Repo) update(ctx context.Context, u *User) error {
	if err := ur.qry.UpdateUser(ctx, &sqlc.UpdateUserParams{
		Username: u.Username,
		Email:    u.Email,
		Pwd:      u.Pwd,
		ID:       u.ID,
	}); err != nil {
		return err
	}
	return nil
}
