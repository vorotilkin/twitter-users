package user

import (
	"context"
	"errors"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-users/domain/models"
	"github.com/vorotilkin/twitter-users/pkg/database"
	"github.com/vorotilkin/twitter-users/schema/gen/my_database/public/model"
	"github.com/vorotilkin/twitter-users/schema/gen/my_database/public/table"
)

type Repository struct {
	conn *database.Database
}

func (r *Repository) UpdateByID(ctx context.Context, userToUpdate models.UserOption) (bool, error) {
	columns, dbUser := columnsAndModelToUpdate(userToUpdate)
	if len(columns) == 0 {
		return false, models.ErrNothingToUpdate
	}

	query, args := table.User.
		UPDATE(columns).
		WHERE(table.User.ID.EQ(postgres.Int(int64(dbUser.ID)))).
		MODEL(dbUser).
		Sql()

	tag, err := r.conn.Exec(ctx, query, args...)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}

func columnsAndModelToUpdate(userToUpdate models.UserOption) (postgres.ColumnList, model.User) {
	user := model.User{
		ID: userToUpdate.ID,
	}

	columns := make(postgres.ColumnList, 0, 10)

	userToUpdate.Name.ForEach(func(name string) {
		columns = append(columns, table.User.Name)
		user.Name = name
	})

	userToUpdate.Username.ForEach(func(username string) {
		columns = append(columns, table.User.Username)
		user.Username = username
	})

	userToUpdate.Bio.ForEach(func(bio string) {
		columns = append(columns, table.User.Bio)
		user.Bio = lo.ToPtr(bio)
	})

	userToUpdate.ProfileImage.ForEach(func(profileImage string) {
		columns = append(columns, table.User.ProfileImage)
		user.ProfileImage = lo.ToPtr(profileImage)
	})

	userToUpdate.CoverImage.ForEach(func(coverImage string) {
		columns = append(columns, table.User.CoverImage)
		user.CoverImage = lo.ToPtr(coverImage)
	})

	return columns, user
}

func (r *Repository) UserByID(ctx context.Context, id int32) (models.User, error) {
	query, args := table.User.
		SELECT(
			table.User.ID,
			table.User.Name,
			table.User.PasswordHash,
			table.User.Username,
			table.User.Email,
			table.User.Bio,
			table.User.ProfileImage,
			table.User.CoverImage,
		).
		WHERE(table.User.ID.EQ(postgres.Int(int64(id)))).
		Sql()

	row := r.conn.QueryRow(ctx, query, args...)
	user := model.User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.PasswordHash,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.ProfileImage,
		&user.CoverImage,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, err
	}

	return toDomain(user), nil
}

func (r *Repository) UserByEmail(ctx context.Context, email string) (models.User, error) {
	query, args := table.User.
		SELECT(
			table.User.ID,
			table.User.Name,
			table.User.PasswordHash,
			table.User.Username,
			table.User.Email,
			table.User.Bio,
			table.User.ProfileImage,
			table.User.CoverImage,
		).
		WHERE(table.User.Email.EQ(postgres.Text(email))).
		Sql()

	row := r.conn.QueryRow(ctx, query, args...)
	user := model.User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.PasswordHash,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.ProfileImage,
		&user.CoverImage,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, err
	}

	return toDomain(user), nil
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

func (r *Repository) FetchPasswordHashByEmail(ctx context.Context, email string) (string, error) {
	query, args := table.User.
		SELECT(table.User.PasswordHash).
		WHERE(table.User.Email.EQ(postgres.Text(email))).
		Sql()

	row := r.conn.QueryRow(ctx, query, args...)
	user := model.User{}

	err := row.Scan(
		&user.PasswordHash,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	return user.PasswordHash, nil
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
