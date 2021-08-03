package domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type AutoPlayers struct {
	Players    map[int]*AutoPlayer  `json:"players"`
}

type AutoPlayer struct {
	Seatnum int                     `json:"seatnum"`
	Actions map[int]Action          `json:"actions"`

}

func readAutomation(filename string) (*AutoPlayers, error) {
	var autoplayers AutoPlayers

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &autoplayers)
	if err != nil {
		fmt.Printf("Error unmarshalling player automation: %v\n", err)
		return nil, err
	}

	return &autoplayers, nil

}

