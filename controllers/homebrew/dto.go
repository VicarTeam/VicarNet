package homebrew

type homebrewDisciplineRequestDto struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Levels  map[uint][]struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Requirement *uint  `json:"requirement"`
		Combination *struct {
			ID    uint `json:"id"`
			Level uint `json:"level"`
		} `json:"combination"`
		MinBloodPotency *uint    `json:"minBloodPotency"`
		Summary         string   `json:"summary"`
		Costs           string   `json:"costs"`
		DiceSupplies    *string  `json:"diceSupplies"`
		System          string   `json:"system"`
		Alternatives    []string `json:"alternatives"`
		Duration        string   `json:"duration"`
	}
}

type homebrewClanRequestDto struct {
	Name        string `json:"name"`
	Slogan      string `json:"slogan"`
	Curse       string `json:"curse"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Disciplines []uint `json:"disciplines"`
}
