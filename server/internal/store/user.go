package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/puremike/realtime_chat_app/internal/model"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeOut)
	defer cancel()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email`

	if err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeOut)
	defer cancel()

	query := `SELECT id, username, email, password FROM users WHERE email = $1`

	user := model.User{}

	if err := u.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserStore) StoreRefreshToken(ctx context.Context, userID int, token string, expires time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeOut)
	defer cancel()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Delete existing tokens for the user
	deleteQuery := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, userID)
	if err != nil {
		return err
	}

	// Insert new token
	insertQuery := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, insertQuery, userID, token, expires)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *UserStore) ValidateRefreshToken(ctx context.Context, token string) (int, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeOut)
	defer cancel()

	var userID int
	var expires time.Time

	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`
	err := u.db.QueryRowContext(ctx, query, token).Scan(&userID, &expires)
	if err != nil {
		return 0, err
	}
	if time.Now().After(expires) {
		return 0, errors.New("refresh token expired")

	}
	return userID, nil
}
