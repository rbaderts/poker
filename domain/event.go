package domain

import (
	"github.com/rbaderts/pokerlib"
	"time"
)

/*
   A domain event represents some kind of state change in the domain.
       All User actions produce a domain event
       Any state change to a table or to a game produces a domain event

   A action request is sent to a specific client when player input is required.

*/

// A dispatcher needs to be provided by the application
type MessageDispatcher interface {
	BroadcastToTable(table *Table, typ EventType, data interface{}) error
	PublishToTable(table *Table, data interface{}) error
	SendRequestAndWait(table *Table, seatnum int, typ EventType, data interface{}, handler ResponseHandler) (interface{}, error)
	SendToSeat(table *Table, seatnum int, typ EventType, data interface{}) (*Event, error)
	RespondToRequest(requestId int, msg *Message) error
}

type ResponseHandler interface {
	Handle(*Message) (interface{}, error)
}

type EventType int

var eventCounter int = 0

type MessageHeader struct {
	Id         int       `json:"id"`
	Time       time.Time `json:"time"`
	Initiator  int       `json:"initiator"`
	ResponseTo int       `json:"responseTo"`
}

type Message struct {
	MessageHeader
	Data interface{}
}

type Event struct {
	Id   int         `json:"id"`
	Typ  EventType   `json:"typ"`
	Data interface{} `json:"data"`
}

/*
type Event struct {
	Id          int
	Typ         EventType
	Data        interface{}
	Time        time.Time
	Initiator   int
}

*/

func NewEvent(typ EventType, data interface{}) *Event {
	e := new(Event)
	e.Id = eventCounter
	eventCounter += 1
	//e.Time = time.Now()
	e.Typ = typ
	e.Data = data
	return e
}

// Event types, Request types
const (
	BetEvent   EventType = 1
	BetRequest EventType = 2 // total bet

	GameHoleCardsDraw EventType = 3
	GamePlayerAction  EventType = 4

	PlayerPaidout               EventType = 5
	GamePlayerActionBet         EventType = 7 // player
	GamePlayerActionFold        EventType = 8 // player
	BettingRoundComplete        EventType = 9
	GameCardsDealt              EventType = 10 // player, amount
	RequestStatus               EventType = 11 // player, amount
	TablePlayerMessage          EventType = 12 // player, amount
	GamePlayerActionMuckCards   EventType = 13 // player
	GamePlayerActionRequestTime EventType = 14 // player
	GamePlayerActionSitOut      EventType = 15 // player
	GamePlayerActionReturned    EventType = 16 // player

	GamePlayerActionRequested EventType = 17 // player, time

	GameShowdown       EventType = 21 // winning player, winning hand
	GamePotDistributed EventType = 22 // amount

	TablePlayerJoined      EventType = 23 // player, time
	TablePlayerLeft        EventType = 24 // player
	TablePlayerSittingOut  EventType = 25 // player, time
	TablePlayerReturned    EventType = 26 // player
	TablePlayerStackChange EventType = 27 // player, amount
	TableGameStarted       EventType = 28 // game
	TableStatusUpdated     EventType = 32 // game
	TableGameCompleted     EventType = 29 // game
	TableMessage           EventType = 35
	PlayerUpdated          EventType = 30
	TableOddsUpdated       EventType = 33
	GameSeatUpdated        EventType = 34
)

func (this EventType) String() string {
	switch this {

	case BetRequest:
		return "BetRequest"

	case GameHoleCardsDraw:
		return "GameHoleCardsDraw"
	case GamePlayerAction:
		return "GamePlayerAction"
	case PlayerPaidout:
		return "PlayerPaidout"
	case GamePlayerActionBet:
		return "GamePlayerActionBet"
	case GamePlayerActionFold:
		return "GamePlayerActionFold"
	case BettingRoundComplete:
		return "BettingRoundComplete"
	case GameCardsDealt:
		return "GameCardsDealt"
	case RequestStatus:
		return "RequestStatus"
	case TablePlayerMessage:
		return "TablePlayerMessage"

	case GamePlayerActionMuckCards:
		return "GamePlayerActionMuckCards"
	case GamePlayerActionRequestTime:
		return "GamePlayerActionRequestTime"
	case GamePlayerActionSitOut:
		return "GamePlayerActionSitOut"
	case GamePlayerActionReturned:
		return "GamePlayerActionReturned"

	case GamePlayerActionRequested:
		return "GamePlayerActionRequested"

	case GameShowdown:
		return "GameShowdown"
	case GamePotDistributed:
		return "GamePotDistributed"
	case TablePlayerJoined:
		return "TablePlayerJoined"
	case TablePlayerLeft:
		return "TablePlayerLeft"
	case TablePlayerSittingOut:
		return "TablePlayerSittingOut"
	case TablePlayerReturned:
		return "TablePlayerReturned"
	case TablePlayerStackChange:
		return "TablePlayerStackChange"
	case TableGameStarted:
		return "TableGameStarted"
	case TableGameCompleted:
		return "TableGameCompleted"
	case PlayerUpdated:
		return "PlayerUpdated"
	case TableOddsUpdated:
		return "TableOddsUpdated"
	case GameSeatUpdated:
		return "GameSeatUpdated"
	case TableMessage:
		return "TableMessage"

	}
	return "Unknown event type"

}

type PlayerBetData struct {
	SeatNum     int     `json:"seatNum"`
	Amount      int     `json:"amount"`
	RaiseAmount int     `json:"raise"`
	BetType     BetType `json:"betType"`
}

type BettingRoundCompleteData struct {
	TotalAmount int `json:"total"`
}

type BetRequestData struct {
	CurrentBet int `json:"currentBet"`
	CallAmount int `json:"callAmount"`
}

type PlayerJoinedData struct {
	PlayerName string `json:"playerName"`
	SeatNum    int    `json:"seatNum"`
	Stack      int    `json:"Stack"`
}

type PlayerLeftData struct {
	SeatNum int `json:"seatNum"`
}

type TableMessageData struct {
	Message string `json:"message"`
}

type GameShowdownData struct {
	Message string   `json:"message"`
	Payouts []Payout `json:"payouts"`
	//	Cards   [][2]pokerlib.Card         `json:"cards"`
	//	Hands   map[int][2]pokerlib.Card   `json:"hands"`
}

type OddsData struct {
	Odds map[int]*pokerlib.Odds `json:"odds"`
}

type GameSeatData struct {
	BestHandDescription string `json:"hand"`
}
