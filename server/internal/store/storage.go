package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/puremike/realtime_chat_app/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	StoreRefreshToken(ctx context.Context, userID int, token string, expires time.Time) error
	ValidateRefreshToken(ctx context.Context, token string) (int, error)
}

type Storage struct {
	User UserRepository
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		User: &UserStore{db},
	}
}

var (
	QueryContextTimeOut = 5 * time.Second
	ErrUserNotFound     = errors.New("user not found")
)
