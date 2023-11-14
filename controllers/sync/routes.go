package sync

import (
	"log"

	"github.com/gin-gonic/gin"
)

type SyncCharData struct {
	Data string `json:"data"`
}

type RetrieveCharsData struct {
	Ids []string `json:"ids"`
}

func Init(r *gin.RouterGroup) {
	r.POST("/characters/:id", postCharSync)
	r.POST("/characters/:id/level", postCharLevelSync)
	r.POST("/characters", getSyncedChars)
}

func postCharSync(c *gin.Context) {
	charId := c.Param("id")

	data := SyncCharData{}
	if err := c.BindJSON(&data); err != nil {
		log.Println(err)
		c.Status(400)
		return
	}

	sendCharSync(charId, data.Data)

	c.JSON(204, gin.H{})
}

func postCharLevelSync(c *gin.Context) {
	charId := c.Param("id")

	data := SyncCharData{}
	if err := c.BindJSON(&data); err != nil {
		log.Println(err)
		c.Status(400)
		return
	}

	sendCharLevelSync(charId, data.Data)

	c.JSON(204, gin.H{})
}

func getSyncedChars(c *gin.Context) {
	data := RetrieveCharsData{}
	if err := c.BindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{})
		return
	}

	result := retrieveCharSync(data.Ids)

	c.JSON(200, result)
}
