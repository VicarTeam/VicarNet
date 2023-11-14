package auth

type registerRequestDto struct {
	Alias string `json:"alias" binding:"required"`
	Email string `json:"email" binding:"required"`
}
