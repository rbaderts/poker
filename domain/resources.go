package domain

import (
	"fmt"
	"github.com/rbaderts/pokerlib"
)

/**
    Defines the resource types that present the external view of our domain

 */
var Colors = []string{"red", "green", "blue", "yellow", "purple", "white"}

type TableResource struct {
	Id             int              `json:"id"`
	Name           string           `json:"name"`
	Channel        string           `json:"channel"`
	Button         int              `json:"button"`
	CommonCards     []pokerlib.Card `json:"commonCards"`
	Stage          GameStage        `json:"stage"`

	PotValue       int              `json:"potValue"`
///	Pots           []*Pot           `json:"pots"`
//	Players          []*SeatResource  `json:"seats"`

	Players        []*PlayerResource `json:"players"`

}


type PlayerResource struct {
	Id       int    `json:"id"`
	Name     string `json:"username"`
	Stack    int    `json:"stack"`
	OnTable  int    `json:"ontable"`
	HasCards bool   `json:"hasCards"`
	Seatnum  int    `json:"seatnum"`
	Color    string `json:"color"`
	Folded   bool   `json:"folded"`
	PercentWin  float32   `json:"percentWin"`
	PercentTie  float32   `json:"percentTie"`
	//Card1     *pokerlib.Card  `json:"card1"`
	//Card2     *pokerlib.Card  `json:"card2"`

	//	Cards     [2]pokerlib.Card  `json:"communityCards"`
}



func CalculateUpdate(table *Table) *TableResource {

	fmt.Printf("CalculateUpdate\n")
	tr := new(TableResource)
	tr.CommonCards = make([]pokerlib.Card, 0)
	tr.Channel = table.Channel
	tr.Id = table.Id
	tr.Button = table.Button

	if table.currentGame != nil {
		tr.Stage = table.currentGame.Stage
		for _, c := range table.currentGame.communityCards {
			tr.CommonCards = append(tr.CommonCards, c)
		}
		tr.PotValue = table.currentGame.PotValue

	}

	tr.Players = make([]*PlayerResource, 0)
	for _, seat := range table.Seats {
		if !seat.SittingOut && seat.UserId != -1 {

			pr := new(PlayerResource)
			fmt.Printf("seat = %v\n", seat)
			user := seat.GetUser()
			pr.Id = user.Id
			pr.Seatnum = seat.Seatnum
			pr.Stack = seat.Stack
			pr.Color = Colors[seat.Seatnum]
			pr.Name = user.PreferredHandle
			pr.HasCards = false
			if table.currentGame != nil {
				pr.PercentWin = table.currentGame.Players[seat.Seatnum].percentWin
				pr.PercentTie = table.currentGame.Players[seat.Seatnum].percentTie
				br := table.currentGame.CurrentBettingRound
				if br != nil {
					pr.OnTable = br.seatBets[seat.Seatnum]
				}
				ss := table.currentGame.Players[seat.Seatnum]
				if ss != nil && ss.cards[0].Suit == 0 {
					pr.HasCards = true
				}
			}
			//			tr.Players = append(tr.Players, pr)
			tr.Players = append(tr.Players, pr)
		}

	}
	return tr
}

