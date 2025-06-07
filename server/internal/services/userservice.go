package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/realtime_chat_app/internal/config"
	"github.com/puremike/realtime_chat_app/internal/model"
	"github.com/puremike/realtime_chat_app/internal/store"
	"github.com/puremike/realtime_chat_app/internal/utils"
)

type UserService struct {
	repo store.UserRepository
	app  *config.Application
}

func NewUserService(repo store.UserRepository, app *config.Application) *UserService {
	return &UserService{repo: repo, app: app}
}

var (
	QueryContextTimeOut = 5 * time.Second
)

func (u *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeOut)
	defer cancel()
	// validation
	if len(user.Username) < 4 || len(user.Password) < 8 || user.Email == "" {
		return nil, fmt.Errorf("fields are not valid")
	}

	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		return nil, err
	}

	us := &model.User{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}

	createdUser, err := u.repo.CreateUser(ctx, us)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {

	if u.app.Config.AuthConfig.Secret == "" || u.app.Config.AuthConfig.Iss == "" || u.app.Config.AuthConfig.Aud == "" {
		return nil, fmt.Errorf("auth config not initialized")
	}

	if req.Email == "" || req.Password == "" {
		return &model.LoginResponse{}, fmt.Errorf("email and password are required")
	}

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeOut)
	defer cancel()

	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return &model.LoginResponse{}, fmt.Errorf("user not found")
		}
		return &model.LoginResponse{}, fmt.Errorf("failed to get user")
	}

	if err := utils.CompareHashedPassword(user.Password, req.Password); err != nil {
		return &model.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"iss": u.app.Config.AuthConfig.Iss,
		"aud": u.app.Config.AuthConfig.Aud,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(u.app.Config.AuthConfig.TokenExp).Unix(),
	}

	token, err := u.app.JwtAuth.GenerateToken(claims)
	if err != nil {
		return &model.LoginResponse{}, fmt.Errorf("failed to generate token")
	}

	// generate refresh token
	refreshToken, err := u.app.JwtAuth.GenerateRefreshToken()
	if err != nil {
		return &model.LoginResponse{}, fmt.Errorf("failed to generate refresh token")
	}

	// Store refresh token in DB (add to UserRepository)
	if err = u.repo.StoreRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(7*24*time.Hour)); err != nil {
		return &model.LoginResponse{}, fmt.Errorf("failed to store refresh token")
	}

	return &model.LoginResponse{ID: user.ID, Username: user.Username, Token: token, RefreshToken: refreshToken}, nil
}
