package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Cache *cache
)

func Init() {
	Cache = newCache()
	initDB()
}

func initDB() {
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{}, &HomebrewDiscipline{}, &HomebrewClan{})

	DB = db
}
