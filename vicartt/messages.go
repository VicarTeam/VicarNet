package vicartt

type DiceRollMessage struct {
	Username string `json:"username"`
	Vampire  struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"vampire"`
	Roll struct {
		Title       string `json:"title"`
		NormalDices int    `json:"normalDices"`
		HungerDices int    `json:"hungerDices"`
		Difficulty  *int   `json:"difficulty"`
	} `json:"roll"`
}
