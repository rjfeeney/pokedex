package pokedex

type Pokemon struct {
	Name   string     `json:"name"`
	BaseXP int        `json:"base_experience"`
	Height int        `json:"height"`
	Weight int        `json:"weight"`
	Stats  []StatInfo `json:"stats"`
	Types  []TypeInfo `json:"types"`
}

type TypeInfo struct {
	Slot int        `json:"slot"`
	Type TypeDetail `json:"type"`
}

type TypeDetail struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type StatInfo struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}

type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokedex struct {
	Caught map[string]Pokemon
}

func (p *Pokedex) AddPokemonToPokedex(pokemon Pokemon) {
	p.Caught[pokemon.Name] = pokemon
}
