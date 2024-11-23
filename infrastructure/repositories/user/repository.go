package user

import (
	"context"
	"github.com/vorotilkin/twitter-users/domain/models"
	"github.com/vorotilkin/twitter-users/pkg/database"
	"github.com/vorotilkin/twitter-users/schema/gen/my_database/public/model"
	"github.com/vorotilkin/twitter-users/schema/gen/my_database/public/table"
)

type Repository struct {
	conn *database.Database
}

func (r *Repository) Create(ctx context.Context, name, passwordHash, username, email string) (models.User, error) {
	query, args := table.User.
		INSERT(table.User.Name, table.User.PasswordHash, table.User.Username, table.User.Email).
		MODEL(model.User{
			Name:         name,
			PasswordHash: passwordHash,
			Username:     username,
			Email:        email,
		}).
		RETURNING(
			table.User.ID,
			table.User.Name,
			table.User.PasswordHash,
			table.User.Username,
			table.User.Email,
		).
		Sql()

	row := r.conn.QueryRow(ctx, query, args...)
	user := model.User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.PasswordHash,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		return models.User{}, err
	}

	return toDomain(user), nil
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
