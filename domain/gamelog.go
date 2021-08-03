package domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)
/*
  A GameLog is a full record of completed game.    It contains descriptive info,
     am order set of Actions that took place and a description of the outcome
 */

type GameLog struct {
	Game        GameHeader   `json:"game"`
	Actions     []*ActionLog `json:"actions"`
	orderedKeys []int
}

type ActionLog struct {
	Sequence int                    `json:"sequence"`
	Typ      string                 `json:"actionType"`
	Seat     int                    `json:"seat"`
	Params   map[string]interface{} `json:"params"`
	Round    int                    `json:"round"`
}

func ActionLogFromAction(a *Action) *ActionLog {
	al := new(ActionLog)
	al.Sequence = a.Sequence
	al.Typ = a.Typ
	al.Seat = a.Seat
	al.Sequence = a.Sequence
	al.Params = make(map[string]interface{}, 0)
	for k, v := range a.Params {
		al.Params[k] = v
	}
	if a.round != nil {
		al.Round = a.round.round
	}

	return al
}

type SeatInfo struct {
	Seatnum  int    `json:"seatnum"`
	Username string `json:"username"`
	UserId   int    `json:"userId"`
	Stack    int    `json:"stack"`
}

type GameHeader struct {
	BigBlind int               `json:"bigblind"`
	Seats    map[int]*SeatInfo `json:"seats"`
	Button   int               `json:"button"`
}

func (this *GameState) dumpGame(filename string) {

	gl := new(GameLog)
	gl.Game.BigBlind = this.BigBlind
	gl.Game.Seats = make(map[int]*SeatInfo, 0)
	gl.Game.Button = this.DealerSeat
	for _, ss := range this.Players {
		seatinfo := new(SeatInfo)
		seatinfo.UserId = ss.UserId
		seatinfo.Username = ss.Seat.User.GivenName
		seatinfo.Stack = ss.gameStack
		gl.Game.Seats[ss.Seatnum] = seatinfo
	}
	//for _, ac := range this.actionHistory {
	//}
	gl.Actions = this.actionHistory
	data, err := json.MarshalIndent(gl, "", "  ")

	if err != nil {
		fmt.Printf("Error marshalling json: %v\n", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func (this *GameState) readGameLog(data []byte) *GameLog {

	//actions := make([]*Action, 0)
	//buf := bytes.NewBuffer(data)

	var gamelog GameLog
	err := json.Unmarshal(data, &gamelog)
	if err != nil {
		fmt.Printf("Error reading log: %v\n", err)
	}

	return &gamelog
}

