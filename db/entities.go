package db

type User struct {
	ID      uint   `gorm:"primaryKey"`
	Alias   string `gorm:"unique"`
	Email   string `gorm:"unique"`
	Passkey string
}

type HomebrewDiscipline struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CreatorID uint    `json:"creatorId"` // the user who created this discipline
	IsPublic  bool    `gorm:"default:false"`
	Name      string  `json:"name"`
	Summary   *string `json:"summary"`
	Levels    string  `json:"levels"` // json object of levels
}

type HomebrewClan struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	CreatorID   uint `json:"creatorId"` // the user who created this clan
	IsPublic    bool `gorm:"default:false"`
	Name        string
	Slogan      string
	Curse       string
	Description string
	Symbol      string
	Disciplines string // json array of discipline ids
}
