package vicartt

import (
	"vicarnet/controllers/auth"

	"github.com/gin-gonic/gin"
)

var server *Server

func Init(r *gin.RouterGroup) {
	server = NewServer()

	r.GET("/stream", server.HeadersMiddleware(), server.PrepareMiddleware(), server.HandleSSE())
	r.POST("/dice-roll/:accessKey", postDiceRoll)
}

func postDiceRoll(c *gin.Context) {
	username := ""

	_, user := auth.Authenticate(c)
	if user != nil {
		username = user.Alias
	}

	accessKey := c.Param("accessKey")
	if accessKey == "" {
		c.JSON(400, gin.H{"error": "accessKey is required"})
		return
	}

	msg := DiceRollMessage{}
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	msg.Username = username

	server.RelayMessage(accessKey, msg)

	c.JSON(200, gin.H{"status": "ok"})
}
