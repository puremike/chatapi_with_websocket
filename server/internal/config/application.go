package config

import (
	"time"

	"github.com/puremike/realtime_chat_app/internal/auth"
	"github.com/puremike/realtime_chat_app/internal/store"
	"github.com/puremike/realtime_chat_app/internal/ws"
)

type Application struct {
	Config  *Config
	Store   *store.Storage
	JwtAuth *auth.JWTAuthenticator
	Hub     *ws.Hub
}

type Config struct {
	Port       string
	Env        string
	DbAddr     string
	AuthConfig *AuthConfig
}

type AuthConfig struct {
	Secret, Iss, Aud string
	TokenExp         time.Duration
}
