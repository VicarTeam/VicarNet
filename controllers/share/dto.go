package share

type setShareIdRequestDto struct {
	ShareId string `json:"shareId" binding:"required"`
}
