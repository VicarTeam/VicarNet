package controllers

import (
	"vicarnet/controllers/auth"
	"vicarnet/controllers/homebrew"
	"vicarnet/controllers/share"
	"vicarnet/controllers/sync"
	"vicarnet/vicartt"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	auth.Init(r.Group("/auth"))
	share.Init(r.Group("/share"))
	homebrew.Init(r.Group("/homebrew"))
	sync.Init(r.Group("/sync"))
	vicartt.Init(r.Group("/vicartt"))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}
