package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"pokedex/internal/pokecache"
	"pokedex/internal/pokedex"
	"time"
)

type Client struct {
	BaseURL    string
	HttpClient *http.Client
	Config     *Config
	Cache      *pokecache.Cache
	Pokedex    *pokedex.Pokedex
}

func NewClient() *Client {
	return &Client{
		BaseURL:    "https://pokeapi.co/api/v2/",
		HttpClient: &http.Client{Timeout: 10 * time.Second},
		Config:     &Config{},
		Cache:      pokecache.NewCache(5 * time.Minute),
		Pokedex:    &pokedex.Pokedex{Caught: make(map[string]pokedex.Pokemon)},
	}
}

func makeRequest[T any](c *Client, url string) (*T, error) {
	if val, ok := c.Cache.Get(url); ok {
		var rsp T
		err := json.Unmarshal(val, &rsp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall cached JSON: %w", err)
		}
		return &rsp, nil
	}
	res, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}

	c.Cache.Add(url, body)

	var rsp T
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall JSON: %w", err)
	}
	return &rsp, nil
}

func (c *Client) CommandMap(config *Config) error {
	url := c.BaseURL + "location-area"
	if config.Next != nil {
		url = *config.Next
	}
	rsp, reqErr := makeRequest[LocationAreaListResponse](c, url)
	if reqErr != nil {
		return reqErr
	}
	config.Next = rsp.Next
	config.Previous = rsp.Previous
	for _, area := range rsp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func (c *Client) CommandMapb(config *Config) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	url := *config.Previous
	rsp, reqErr := makeRequest[LocationAreaListResponse](c, url)
	if reqErr != nil {
		return reqErr
	}
	config.Next = rsp.Next
	config.Previous = rsp.Previous
	for _, area := range rsp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func (c *Client) CommandExplore(locationArea string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + locationArea + "/"
	rsp, reqErr := makeRequest[LocationAreaDetailResponse](c, url)
	if reqErr != nil {
		return reqErr
	}
	fmt.Printf("Exploring %s...\n", locationArea)
	fmt.Println("Found Pokemon:")
	for _, encounter := range rsp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func (c *Client) CommandCatch(pokemon string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon + "/"
	rsp, reqErr := makeRequest[pokedex.Pokemon](c, url)
	if reqErr != nil {
		return reqErr
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", rsp.Name)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	catchRoll := r.Intn(rsp.BaseXP)
	if catchRoll < 70 {
		fmt.Printf("%s was caught!\n", rsp.Name)
		c.Pokedex.AddPokemonToPokedex(*rsp)
		fmt.Printf("%s's info was added to the PokeDex!\n", rsp.Name)
	} else {
		fmt.Printf("%s escaped!\n", rsp.Name)
	}
	return nil
}

func (c *Client) CommandInspect(pokemon string) error {
	caught, exists := c.Pokedex.Caught[pokemon]
	if exists {
		fmt.Printf("Name: %s\n", caught.Name)
		fmt.Printf("Height: %d\n", caught.Height)
		fmt.Printf("Weight: %d\n", caught.Weight)
		fmt.Println("Stats:")
		for _, statInfo := range caught.Stats {
			fmt.Printf(" - %s: %d\n", statInfo.Stat.Name, statInfo.BaseStat)
		}
		fmt.Println("Types:")
		for _, typeInfo := range caught.Types {
			fmt.Printf("  - %s\n", typeInfo.Type.Name)
		}
	} else {
		fmt.Println("You have not caught that pokemon")
	}
	return nil
}

func (c *Client) CommandPokedex() error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range c.Pokedex.Caught {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}
