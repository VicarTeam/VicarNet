package auth

import (
	"log"
	"time"
	"vicarnet/db"
	"vicarnet/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Init(r *gin.RouterGroup) {
	r.POST("/register", register)
	r.PUT("/activate/:activationCode", activate)
	r.POST("/login", login)
	r.PATCH("/recover/:email/begin", beginRecoverAccount)
	r.PUT("/recover/:recoverCode/finish", finishRecoverAccount)
}

func Authenticate(c *gin.Context) (bool, *db.User) {
	userID := c.Request.Header.Get("X-User-ID")
	passkey := c.Request.Header.Get("X-User-Passkey")

	user := &db.User{}
	if err := db.DB.Where("id = ? and passkey = ?", userID, passkey).First(user).Error; err != nil {
		return false, nil
	}

	return true, user
}

func register(c *gin.Context) {
	dto := &registerRequestDto{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.JSON(400, gin.H{"message": "bad request"})
		return
	}

	if err := db.DB.Where("alias = ? or email = ?", dto.Alias, dto.Email).First(&db.User{}).Error; err == nil {
		c.JSON(400, gin.H{"message": "alias or email already exists"})
		return
	}

	activationCode := uuid.New().String()
	activationExp := time.Now().Add(time.Minute * 15)
	db.Cache.Set("register:"+activationCode, dto, &activationExp)
	log.Println("activation code: " + activationCode)

	util.SendMail(dto.Email, "VicarNet - Activation", "Your activation code is: "+activationCode+"\n\nThis code will expire in 15 minutes.")

	c.JSON(202, gin.H{"message": "await activation"})
}

func activate(c *gin.Context) {
	activationCode := c.Param("activationCode")

	dto, ok := db.Cache.Get("register:" + activationCode)
	if !ok {
		c.JSON(400, gin.H{"message": "invalid activation code"})
		return
	}

	passkey, err := util.GenerateRandomString(128)
	if err != nil {
		panic(err)
	}

	user := &db.User{
		Alias:   dto.(*registerRequestDto).Alias,
		Email:   dto.(*registerRequestDto).Email,
		Passkey: passkey,
	}

	if err := db.DB.Create(user).Error; err != nil {
		panic(err)
	}

	c.JSON(201, gin.H{
		"message": "activated",
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"alias":   user.Alias,
			"passkey": user.Passkey,
		},
	})
}

func login(c *gin.Context) {
	if ok, _ := Authenticate(c); ok {
		c.JSON(200, gin.H{"message": "logged in"})
		return
	}

	c.JSON(401, gin.H{"message": "unauthorized"})
}

func beginRecoverAccount(c *gin.Context) {
	email := c.Param("email")

	user := &db.User{}
	if err := db.DB.Where("email = ?", email).First(user).Error; err != nil {
		c.JSON(400, gin.H{"message": "invalid email"})
		return
	}

	recoverCode := uuid.New().String()
	recoverExp := time.Now().Add(time.Minute * 15)
	db.Cache.Set("recover:"+recoverCode, user, &recoverExp)
	log.Println("recover code: " + recoverCode)

	util.SendMail(user.Email, "VicarNet - Recover", "Your recover code is: "+recoverCode+"\n\nThis code will expire in 15 minutes.")

	c.JSON(202, gin.H{"message": "await recover"})
}

func finishRecoverAccount(c *gin.Context) {
	recoverCode := c.Param("recoverCode")

	userInterface, ok := db.Cache.Get("recover:" + recoverCode)
	if !ok {
		c.JSON(400, gin.H{"message": "invalid recover code"})
		return
	}

	user := userInterface.(*db.User)
	passkey, err := util.GenerateRandomString(128)
	if err != nil {
		panic(err)
	}

	user.Passkey = passkey
	if err := db.DB.Model(user).Save(user).Error; err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"message": "recovered",
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"alias":   user.Alias,
			"passkey": user.Passkey,
		},
	})
}
