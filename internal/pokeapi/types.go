package pokeapi

type LocationAreaListResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetailResponse struct {
	Name              string `json:"name"`
	ID                int    `json:"id"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Config struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}
