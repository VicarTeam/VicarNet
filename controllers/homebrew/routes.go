package homebrew

import (
	"log"
	"strconv"
	"vicarnet/controllers/auth"
	"vicarnet/db"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.RouterGroup) {
	r.GET("/clans", getHomebrewClans)
	r.GET("/clans/:id", getHomebrewClan)
	r.POST("/clans/:id/invite", giveAccessToClan)
	r.POST("/clans", createHomebrewClan)
	r.PATCH("/clans/:id", updateHomebrewClan)
	r.DELETE("/clans/:id", deleteHomebrewClan)
	r.GET("/disciplines", getHomebrewDisciplines)
	r.GET("/disciplines/:id", getHomebrewDiscipline)
	r.POST("/disciplines/:id/invite", giveAccessToDiscipline)
	r.POST("/disciplines", createHomebrewDiscipline)
	r.PATCH("/disciplines/:id", updateHomebrewDiscipline)
	r.DELETE("/disciplines/:id", deleteHomebrewDiscipline)
	r.GET("/my-content", getMyHomebrewContent)
	r.GET("/from-invite/:inviteCode", getHomebrewFromInvite)
}

func parsePagination(c *gin.Context) (int, int) {
	page := c.Query("page")
	limit := c.Query("limit")

	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	parsedPage, err := strconv.Atoi(page)
	if err != nil {
		panic(err)
	}

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		panic(err)
	}

	return parsedPage, parsedLimit
}

func getHomebrewClans(c *gin.Context) {
	searchText := c.Query("search")
	page, limit := parsePagination(c)

	clans, count := getClans(searchText, page, limit)

	clansDto := []gin.H{}
	for _, clan := range clans {
		clansDto = append(clansDto, transformClanToDto(clan))
	}

	c.JSON(200, gin.H{
		"total": count,
		"items": clansDto,
	})
}

func getHomebrewClan(c *gin.Context) {
	_, user := auth.Authenticate(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid id"})
		return
	}

	clan := getClan(user, uint(id), false)
	if clan == nil {
		c.JSON(404, gin.H{"message": "Clan not found"})
		return
	}

	clanDto := transformClanToDto(clan)
	neededHomebrewDisciplines := []gin.H{}
	disciplines := parseDisciplinesJson(clan.Disciplines)

	for _, disciplineId := range disciplines {
		if disciplineId >= 3021 {
			realDisciplineId := disciplineId - 3021
			discipline := getDiscipline(nil, uint(realDisciplineId), true)
			if discipline != nil {
				dto := transformDisciplineToDto(discipline)
				dto["id"] = disciplineId
				neededHomebrewDisciplines = append(neededHomebrewDisciplines, dto)
			}
		}
	}

	c.JSON(200, gin.H{
		"clan":                      clanDto,
		"neededHomebrewDisciplines": neededHomebrewDisciplines,
	})
}

func getHomebrewDisciplines(c *gin.Context) {
	searchText := c.Query("search")
	page, limit := parsePagination(c)

	disciplines, count := getDisciplines(searchText, page, limit)

	disciplinesDto := []gin.H{}
	for _, discipline := range disciplines {
		disciplinesDto = append(disciplinesDto, transformDisciplineToDto(discipline))
	}

	c.JSON(200, gin.H{
		"total": count,
		"items": disciplinesDto,
	})
}

func getHomebrewDiscipline(c *gin.Context) {
	_, user := auth.Authenticate(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid id"})
		return
	}

	discipline := getDiscipline(user, uint(id), false)
	if discipline == nil {
		c.JSON(404, gin.H{"message": "Discipline not found"})
		return
	}

	disciplineDto := transformDisciplineToDto(discipline)

	c.JSON(200, &disciplineDto)
}

func giveAccessToClan(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	clanId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid id"})
		return
	}

	clan := getPrivateClan(uint(clanId), user.ID)
	if clan == nil {
		c.JSON(404, gin.H{"message": "Clan not found"})
		return
	}

	if clan.IsPublic {
		c.JSON(409, gin.H{"message": "Clan is already public"})
		return
	}

	inviteCode := generateHomebrewInviteCode(HomebrewTypeClan, clan.ID)

	c.JSON(200, gin.H{"inviteCode": inviteCode})
}

func giveAccessToDiscipline(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	disciplineId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid id"})
		return
	}

	discipline := getPrivateDiscipline(uint(disciplineId), user.ID)
	if discipline == nil {
		c.JSON(404, gin.H{"message": "Discipline not found"})
		return
	}

	if discipline.IsPublic {
		c.JSON(409, gin.H{"message": "Discipline is already public"})
		return
	}

	inviteCode := generateHomebrewInviteCode(HomebrewTypeDiscipline, discipline.ID)

	c.JSON(200, gin.H{"inviteCode": inviteCode})
}

func createHomebrewClan(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	var dto homebrewClanRequestDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	clan := &db.HomebrewClan{
		CreatorID: user.ID,
	}

	applyClanDtoToEntity(clan, dto)

	if res := db.DB.Create(clan); res.Error != nil {
		panic(res.Error)
	}

	c.JSON(200, transformClanToDto(clan))
}

func createHomebrewDiscipline(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	var dto homebrewDisciplineRequestDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	discipline := &db.HomebrewDiscipline{
		CreatorID: user.ID,
	}

	applyDisciplineDtoToEntity(discipline, dto)

	if res := db.DB.Create(discipline); res.Error != nil {
		panic(res.Error)
	}

	c.JSON(200, transformDisciplineToDto(discipline))
}

func updateHomebrewClan(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid id"})
		return
	}

	clan := getPrivateClan(uint(id), user.ID)
	if clan == nil {
		c.JSON(404, gin.H{"message": "Clan not found"})
		return
	}

	var dto homebrewClanRequestDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	applyClanDtoToEntity(clan, dto)

	if res := db.DB.Save(clan); res.Error != nil {
		panic(res.Error)
	}

	c.JSON(200, transformClanToDto(clan))
}

func updateHomebrewDiscipline(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid id"})
		return
	}

	discipline := getPrivateDiscipline(uint(id), user.ID)
	if discipline == nil {
		c.JSON(404, gin.H{"message": "Discipline not found"})
		return
	}

	var dto homebrewDisciplineRequestDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	applyDisciplineDtoToEntity(discipline, dto)

	if res := db.DB.Save(discipline); res.Error != nil {
		panic(res.Error)
	}

	c.JSON(200, transformDisciplineToDto(discipline))
}

func deleteHomebrewClan(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid id"})
		return
	}

	clan := getPrivateClan(uint(id), user.ID)
	if clan == nil {
		c.JSON(404, gin.H{"message": "Clan not found"})
		return
	}

	if res := db.DB.Delete(clan); res.Error != nil {
		panic(res.Error)
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func deleteHomebrewDiscipline(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid id"})
		return
	}

	discipline := getPrivateDiscipline(uint(id), user.ID)
	if discipline == nil {
		c.JSON(404, gin.H{"message": "Discipline not found"})
		return
	}

	if res := db.DB.Delete(discipline); res.Error != nil {
		panic(res.Error)
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func getMyHomebrewContent(c *gin.Context) {
	ok, user := auth.Authenticate(c)
	if !ok {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	clans := []*db.HomebrewClan{}
	disciplines := []*db.HomebrewDiscipline{}

	db.DB.Where("creator_id = ?", user.ID).Find(&clans)
	db.DB.Where("creator_id = ?", user.ID).Find(&disciplines)

	clansDto := []gin.H{}
	for _, clan := range clans {
		clansDto = append(clansDto, transformClanToDto(clan))
	}

	disciplinesDto := []gin.H{}
	for _, discipline := range disciplines {
		disciplinesDto = append(disciplinesDto, transformDisciplineToDto(discipline))
	}

	c.JSON(200, gin.H{
		"clans":       clansDto,
		"disciplines": disciplinesDto,
	})
}

func getHomebrewFromInvite(c *gin.Context) {
	inviteCode := c.Param("inviteCode")
	if inviteCode == "" {
		c.JSON(400, gin.H{"message": "Invalid invite code"})
		return
	}

	ok, hbType, id := validateHomebrewInviteCode(inviteCode)
	if !ok {
		c.JSON(400, gin.H{"message": "Invalid invite code"})
		return
	}

	if hbType == HomebrewTypeClan {
		clan := getClan(nil, id, true)
		if clan == nil {
			c.JSON(404, gin.H{"message": "Clan not found"})
			return
		}

		clanDto := transformClanToDto(clan)
		neededHomebrewDisciplines := []gin.H{}
		disciplines := parseDisciplinesJson(clan.Disciplines)

		for _, disciplineId := range disciplines {
			log.Println(disciplineId)
			if disciplineId >= 3021 {
				realDisciplineId := disciplineId - 3021
				log.Println(realDisciplineId)
				discipline := getDiscipline(nil, uint(realDisciplineId), true)
				if discipline != nil {
					log.Println("found discipline")
					dto := transformDisciplineToDto(discipline)
					dto["id"] = disciplineId
					neededHomebrewDisciplines = append(neededHomebrewDisciplines, dto)
				}
			}
		}

		c.JSON(200, gin.H{
			"type": "clan",
			"content": gin.H{
				"clan":                      clanDto,
				"neededHomebrewDisciplines": neededHomebrewDisciplines,
			},
		})
	} else {
		discipline := getDiscipline(nil, id, true)
		if discipline == nil {
			c.JSON(404, gin.H{"message": "Discipline not found"})
			return
		}

		disciplineDto := transformDisciplineToDto(discipline)

		c.JSON(200, gin.H{
			"type":    "discipline",
			"content": disciplineDto,
		})
	}

}
