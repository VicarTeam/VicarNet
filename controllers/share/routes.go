package share

import (
	"vicarnet/controllers/auth"
	"vicarnet/db"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.RouterGroup) {
	r.POST("/id", setShareId)
	r.DELETE("/id", removeShareId)
	r.GET("/:alias/id", getShareId)
}

func setShareId(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	dto := &setShareIdRequestDto{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.JSON(400, gin.H{"message": "bad request"})
		return
	}

	db.Cache.Set("share:"+user.Alias, dto.ShareId, nil)

	c.JSON(200, gin.H{"message": "ok"})
}

func removeShareId(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	db.Cache.Delete("share:" + user.Alias)

	c.JSON(200, gin.H{"message": "ok"})
}

func getShareId(c *gin.Context) {
	alias := c.Param("alias")

	shareId, ok := db.Cache.Get("share:" + alias)
	if !ok {
		c.JSON(404, gin.H{"message": "not found"})
		return
	}

	c.JSON(200, gin.H{"id": shareId})
}
