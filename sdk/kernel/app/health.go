package app

import (
	"github.com/gin-gonic/gin"

	"hei-gin/sdk/config"
)

func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": config.C.App.Name + " is running",
		"version": config.C.App.Version,
	})
}
