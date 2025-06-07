package main

import (
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/puremike/realtime_chat_app/db"
	"github.com/puremike/realtime_chat_app/internal/auth"
	"github.com/puremike/realtime_chat_app/internal/config"
	"github.com/puremike/realtime_chat_app/internal/routes"
	"github.com/puremike/realtime_chat_app/internal/store"
	"github.com/puremike/realtime_chat_app/internal/ws"
	"github.com/puremike/realtime_chat_app/pkg"
)

func main() {

	cfg := &config.Config{
		Port:   pkg.GetString("PORT", "8080"),
		Env:    pkg.GetString("ENV", "development"),
		DbAddr: pkg.GetString("DB_ADDR", ""),
		AuthConfig: &config.AuthConfig{
			Secret:   pkg.GetString("JWT_SECRET", "a9545b8bb27c6291ba5c7dda13b661a508735b31b2f4165957bb509bc09aa164"),
			Iss:      pkg.GetString("JWT_ISS", "realtimechatapp"),
			Aud:      pkg.GetString("JWT_AUD", "realtimechatapp"),
			TokenExp: pkg.GetDuration("JWT_EXP", 15*time.Minute),
		},
	}

	db, err := db.NewPostGresDB(cfg.DbAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("DB connected successfully")

	app := &config.Application{
		Config:  cfg,
		Store:   store.NewStorage(db),
		JwtAuth: auth.NewJWTAuthenticator(cfg.AuthConfig.Secret, cfg.AuthConfig.Iss, cfg.AuthConfig.Aud),
		Hub:     ws.NewHub(),
	}

	go app.Hub.Run()

	mux := routes.Routes(app)
	log.Fatal(routes.Server(mux, app.Config.Port))
}
