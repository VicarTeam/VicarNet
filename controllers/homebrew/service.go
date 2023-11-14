package homebrew

import (
	"encoding/json"
	"log"
	"time"
	"vicarnet/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HomebrewType string

type InviteData struct {
	Id   uint         `json:"id"`
	Type HomebrewType `json:"type"`
}

const (
	HomebrewTypeClan       HomebrewType = "clan"
	HomebrewTypeDiscipline HomebrewType = "discipline"
)

func findCreatorName(id uint) string {
	var user db.User

	where := "id = ?"
	result := db.DB.Where(where, id).First(&user)

	if result.Error != nil {
		log.Println(result.Error)
		return "gelöschter Benutzer"
	}

	return user.Alias
}

func generateHomebrewInviteCode(hbType HomebrewType, id uint) string {
	code := uuid.New().String()
	codeExp := time.Now().Add(5 * time.Minute)
	data := InviteData{
		Id:   id,
		Type: hbType,
	}
	db.Cache.Set("homebrew_invite_code_"+code, data, &codeExp)

	return code
}

func validateHomebrewInviteCode(code string) (bool, HomebrewType, uint) {
	key := "homebrew_invite_code_" + code

	result, ok := db.Cache.Get(key)
	if !ok {
		return false, "", 0
	}

	data, ok := result.(InviteData)
	if !ok {
		return false, "", 0
	}

	db.Cache.Delete(key)

	return true, data.Type, data.Id
}

func parseDisciplinesJson(jsonStr string) []uint {
	var disciplines []uint

	if err := json.Unmarshal([]byte(jsonStr), &disciplines); err != nil {
		log.Println(err)
	}

	return disciplines
}

func getClans(search string, page int, limit int) ([]*db.HomebrewClan, int64) {
	var clans []*db.HomebrewClan
	var count int64

	where := "is_public = 1 AND (name LIKE ? OR slogan LIKE ? OR curse LIKE ? OR description LIKE ?)"

	db.DB.Model(&db.HomebrewClan{}).Where(where, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").Count(&count)
	db.DB.Model(&db.HomebrewClan{}).Where(where, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").Offset((page - 1) * limit).Limit(limit).Find(&clans)

	return clans, count
}

func getClan(user *db.User, id uint, force bool) *db.HomebrewClan {
	var clan db.HomebrewClan

	if res := db.DB.Where("id = ?", id).First(&clan); res.Error != nil {
		log.Println(res.Error)
		return nil
	}

	if !clan.IsPublic && !force {
		if user != nil && user.ID == clan.CreatorID {
			return &clan
		}

		return nil
	}

	return &clan
}

func getPrivateClan(id uint, creatorId uint) *db.HomebrewClan {
	var clan db.HomebrewClan

	if res := db.DB.Where("id = ? AND creator_id = ?", id, creatorId).First(&clan); res.Error != nil {
		log.Println(res.Error)
		return nil
	}

	return &clan
}

func transformClanToDto(clan *db.HomebrewClan) gin.H {
	if clan == nil {
		return nil
	}

	disciplinesArr := []uint{}
	if err := json.Unmarshal([]byte(clan.Disciplines), &disciplinesArr); err != nil {
		log.Println(err)
	}

	return gin.H{
		"id":          clan.ID,
		"creator":     findCreatorName(clan.CreatorID),
		"name":        clan.Name,
		"slogan":      clan.Slogan,
		"curse":       clan.Curse,
		"symbol":      clan.Symbol,
		"description": clan.Description,
		"disciplines": disciplinesArr,
	}
}

func applyClanDtoToEntity(clan *db.HomebrewClan, dto homebrewClanRequestDto) {
	clan.Name = dto.Name
	clan.Slogan = dto.Slogan
	clan.Curse = dto.Curse
	clan.Description = dto.Description
	clan.Symbol = dto.Symbol

	disciplinesArr, err := json.Marshal(dto.Disciplines)
	if err != nil {
		log.Println(err)
	}

	clan.Disciplines = string(disciplinesArr)
}

func getDisciplines(search string, page int, limit int) ([]*db.HomebrewDiscipline, int64) {
	var disciplines []*db.HomebrewDiscipline
	var count int64

	where := "is_public = 1 AND (name LIKE ? OR summary LIKE ?)"

	db.DB.Model(&db.HomebrewDiscipline{}).Where(where, "%"+search+"%", "%"+search+"%").Count(&count)
	db.DB.Model(&db.HomebrewDiscipline{}).Where(where, "%"+search+"%", "%"+search+"%").Offset((page - 1) * limit).Limit(limit).Find(&disciplines)

	return disciplines, count
}

func getDiscipline(user *db.User, id uint, force bool) *db.HomebrewDiscipline {
	var discipline db.HomebrewDiscipline

	if res := db.DB.Where("id = ?", id).First(&discipline); res.Error != nil {
		log.Println(res.Error)
		return nil
	}

	if !discipline.IsPublic && !force {
		if user != nil && user.ID == discipline.CreatorID {
			return &discipline
		}

		return nil
	}

	return &discipline
}

func getPrivateDiscipline(id uint, creatorId uint) *db.HomebrewDiscipline {
	var discipline db.HomebrewDiscipline

	if res := db.DB.Where("id = ? AND creator_id = ?", id, creatorId).First(&discipline); res.Error != nil {
		log.Println(res.Error)
		return nil
	}

	return &discipline
}

func transformDisciplineToDto(discipline *db.HomebrewDiscipline) gin.H {
	if discipline == nil {
		return nil
	}

	jsonMap := make(map[uint][]interface{})
	if err := json.Unmarshal([]byte(discipline.Levels), &jsonMap); err != nil {
		log.Println(err)
	}

	return gin.H{
		"id":      discipline.ID,
		"creator": findCreatorName(discipline.CreatorID),
		"name":    discipline.Name,
		"summary": discipline.Summary,
		"levels":  jsonMap,
	}
}

func applyDisciplineDtoToEntity(discipline *db.HomebrewDiscipline, dto homebrewDisciplineRequestDto) {
	discipline.Name = dto.Name
	discipline.Summary = &dto.Summary

	jsonBytes, err := json.Marshal(dto.Levels)
	if err != nil {
		panic(err)
	}

	discipline.Levels = string(jsonBytes)
}
