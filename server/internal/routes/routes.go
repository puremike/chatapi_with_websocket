package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/puremike/realtime_chat_app/internal/config"
	"github.com/puremike/realtime_chat_app/internal/handlers"
	"github.com/puremike/realtime_chat_app/internal/services"
	"github.com/puremike/realtime_chat_app/internal/ws"
)

func Routes(app *config.Application) http.Handler {

	g := gin.Default()

	userService := services.NewUserService(app.Store.User, app)
	userHandler := handlers.NewUserHandler(userService, app)
	wsHandler := ws.NewWSHandler(app.Hub)

	api := g.Group("/api/v1")
	{
		api.GET("/healthcheck", handlers.HealthCheck)
		api.POST("/signup", userHandler.CreateUser)
		api.POST("/login", userHandler.Login)
		api.POST("/logout", userHandler.Logout)
		api.POST("/refresh", handlers.Refresh(app))
		api.POST("/ws/createRoom", wsHandler.CreateRoom)
		api.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
		api.GET("/ws/getRooms", wsHandler.GetRooms)
		api.GET("/ws/getClients/:roomId", wsHandler.GetClients)
	}

	return g
}

func Server(mux http.Handler, PORT string) error {

	server := &http.Server{
		Addr:         ":" + PORT,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("starting server on port: %s", PORT)
	return server.ListenAndServe()
}
