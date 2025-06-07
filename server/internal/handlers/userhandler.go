package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/realtime_chat_app/internal/config"
	"github.com/puremike/realtime_chat_app/internal/model"
	"github.com/puremike/realtime_chat_app/internal/services"
)

type UserHandler struct {
	service *services.UserService
	app     *config.Application
}

func NewUserHandler(service *services.UserService, app *config.Application) *UserHandler {
	return &UserHandler{service: service, app: app}
}

func (u *UserHandler) CreateUser(c *gin.Context) {

	var payload model.CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	createdUser, err := u.service.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a user"})
		return
	}

	response := model.CreateUserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
	}

	c.JSON(http.StatusCreated, response)
}

func (u *UserHandler) Login(c *gin.Context) {

	var payload model.LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.service.Login(c.Request.Context(), &payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials"})
		return
	}

	c.SetCookie("jwt", user.Token, int(u.app.Config.AuthConfig.TokenExp.Seconds()), "/", "localhost", false, true)
	c.SetSameSite(http.SameSiteLaxMode)

	// Set refresh token cookie (long-lived) - valid for 7 days
	c.SetCookie("refresh_token", user.RefreshToken, 604800, "/", "localhost", false, true)
	c.SetSameSite(http.SameSiteLaxMode)

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (u *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.SetSameSite(http.SameSiteLaxMode)

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
