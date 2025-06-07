package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var APIVersion = "1.1"

type HealthResponse struct {
	ApiVersion string `json:"apiVersion"`
	Status     string `json:"status"`
	Info       string `json:"info"`
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		ApiVersion: APIVersion,
		Status:     "OK",
		Info:       "API is running",
	})
}
