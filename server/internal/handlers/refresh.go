package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/realtime_chat_app/internal/config"
)

func Refresh(app *config.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
			return
		}

		// Verify refresh token in DB
		userID, err := app.Store.User.ValidateRefreshToken(c.Request.Context(), refreshToken)
		if err != nil || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		// Generate new JWT
		claims := jwt.MapClaims{
			"sub": userID,
			"iss": app.Config.AuthConfig.Iss,
			"aud": app.Config.AuthConfig.Aud,
			"iat": time.Now().Unix(),
			"nbf": time.Now().Unix(),
			"exp": time.Now().Add(app.Config.AuthConfig.TokenExp).Unix(),
		}

		newToken, err := app.JwtAuth.GenerateToken(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		// Set new JWT cookie
		c.SetCookie("jwt", newToken, int(app.Config.AuthConfig.TokenExp.Seconds()), "/", "localhost", false, true) // 15 min
		c.SetSameSite(http.SameSiteLaxMode)

		c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
	}
}
