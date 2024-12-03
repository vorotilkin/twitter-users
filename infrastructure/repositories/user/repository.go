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

func (r *Repository) UsersByIDs(ctx context.Context, ids []int32) ([]models.User, error) {
	userIDs := lo.Map(ids, func(id int32, _ int) postgres.Expression {
		return postgres.Int(int64(id))
	})

	followingSubquery :=
		table.Follow.SELECT(
			postgres.Raw("ARRAY_AGG(follow.following_user_id)")).WHERE(
			table.Follow.UserID.EQ(table.User.ID),
		)

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
			followingSubquery.AS("following_ids"),
		).
		WHERE(table.User.ID.IN(userIDs...)).
		Sql()

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.User, error) {
		user := model.User{}

		var followingIDs []int32

		err := row.Scan(
			&user.ID,
			&user.Name,
			&user.PasswordHash,
			&user.Username,
			&user.Email,
			&user.Bio,
			&user.ProfileImage,
			&user.CoverImage,
			&followingIDs,
		)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, err
		}

		return toDomain(user, followingIDs), nil
	})
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

	return toDomain(user, nil), nil
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

	return toDomain(user, nil), nil
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

func (r *Repository) Follow(ctx context.Context, userID, targetUserID int32) (bool, error) {
	query, args := table.Follow.
		INSERT(table.Follow.UserID, table.Follow.FollowingUserID).
		MODEL(model.Follow{
			UserID:          userID,
			FollowingUserID: targetUserID,
		}).
		Sql()

	commandTag, err := r.conn.Exec(ctx, query, args...)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func (r *Repository) Unfollow(ctx context.Context, userID, targetUserID int32) (bool, error) {
	query, args := table.Follow.
		DELETE().WHERE(
		table.Follow.UserID.EQ(postgres.Int(int64(userID))).
			AND(
				table.Follow.FollowingUserID.EQ(postgres.Int(int64(targetUserID))),
			),
	).
		Sql()

	commandTag, err := r.conn.Exec(ctx, query, args...)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
