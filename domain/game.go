package domain

import (
	"container/ring"
	"context"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/fatih/color"
	"github.com/rbaderts/pokerlib"
	"strconv"
	"strings"
	"sync"
)

/*
   A Game represents a single poker hand.     A game has a state  which can be not yet started, underway or completed.
   Once started a game steps through it's various stages, requesting input from the players when required, upon
   completion of the game the winner(s) are determined and the pot distributed accordingly.

   The set of players in the given when a game is start  . The game state maintains the state of
      each seat as teh game progresses.

   Player communication:    The game needs to communicate with all the player's active in the game.   If the game
   cannot establish communication with a player when required an attempt will be made to reestablish communication
   before taking a default action on behalf of the player...

   As a game progresses various types of domain events are produced which describe state changes in the game.
   These events are consumed by clients.

*/

type BetType int

var GameCounter int = 0

const (
	Fold       BetType = 1
	Call       BetType = 2
	Check      BetType = 3
	Raise      BetType = 4
	AllIn      BetType = 5
	SmallBlind BetType = 6
	BigBlind   BetType = 7
)

type GameState struct {
	Id string `json"id"`

	// Players that entered into this hand
	Players map[int]*Player

	ActivePlayers *ring.Ring // Ring buffer contains seats still in the hand
	DealerSeat    int
	PotValue      int
	Pots          *treemap.Map

	BigBlind      int
	PlayCompleted bool

	CurrentBettingRound *Round
	Stage               GameStage
	ActionQueue         []*Action

	actionChannel       chan Action
	lock                sync.Mutex
	bettingRoundCounter int

	deck           *pokerlib.Deck
	communityCards []pokerlib.Card
	table          *Table
	potCount       int

	//actionLog           *log.Logger

	allIns map[int]int // The all in amount for each seat that has gone all in

	playerAutomationFile string

	actionHistory []*ActionLog
	seatLock      sync.Mutex // Ensure the no seat is emptied whille looping over seats
}

type GameStateRepository interface {
	GetGameStateForSeat(ctx context.Context, seatNum int)
	UpdateGameState(ctx context.Context, game *GameState) error
}

/*
func NewGameWithAutoPlayers(table *Table, seats []int, dealerSeat int, autofile string) *GameState {

	g := NewGame(table, seats, dealerSeat)
	g.playerAutomationFile = autofile
	return g
}
*/

func NewGame(table *Table, players *ring.Ring, dealerSeat int) *GameState {

	g := new(GameState)
	id := GameCounter
	GameCounter += 1
	g.Stage = PreDeal
	g.Id = strconv.Itoa(id)
	g.table = table
	g.deck = pokerlib.NewDeck()
	g.deck.Shuffle()
	g.communityCards = make([]pokerlib.Card, 0)
	g.allIns = make(map[int]int, 0)
	g.Players = make(map[int]*Player, 0)
	//g.ActivePlayers = ring.New(len(seats))
	g.ActivePlayers = players
	g.DealerSeat = dealerSeat
	g.bettingRoundCounter = 0
	g.actionHistory = make([]*ActionLog, 0)
	r := g.ActivePlayers
	last := r.Prev()
	var dealer *ring.Ring
	for {
		seatnum := r.Value.(int)
		g.Players[seatnum] = g.newPlayer(table.Seats[seatnum])
		if seatnum == dealerSeat {
			dealer = r
		}
		if r == last {
			break
		}
		r = r.Next()
	}

	g.ActivePlayers = dealer
	/*
		for _, seatnum := range seats {
			r.Value = seatnum
			r = r.Next()
			g.Players[seatnum] = g.newPlayer(table.Players[seatnum])
		}
	*/

	// Root the ring at seat after dealer
	/*
		last := r
		for {
			r = r.Next()
			if r == last {
				break
			}
			if r.Value.(int) == dealerSeat {
				g.ActivePlayers = r.Next()
			}
		}
	*/

	g.Pots = treemap.NewWithIntComparator()
	g.NewPot(0)
	g.potCount = 1
	g.BigBlind = table.BigBlind

	//now := time.Now().UnixNano()
	//timestr := strconv.FormatInt(now, 10)

	/*
		filename := "GameActionLog-"+g.Id+"-"+timestr

		if filename != "" {
			f, err := os.OpenFile(filename,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			g.actionLog = log.New(f, "", log.LstdFlags)
		}
	*/
	return g
}

func (this *GameState) addBet(seat int, amount int, isAllIn bool) bool {

	this.Players[seat].gameStack -= amount
	totalBet := this.Players[seat].totalGameBet + amount

	pot := this.GetPot(totalBet)

	if pot.max > 0 && pot.max < totalBet {
		extra := totalBet - pot.max
		fmt.Printf("pot.max = %d, totalBet = %d\n", pot.max, totalBet)
		newpot := this.NewPot(pot.max)
		newpot.contributions[seat] = extra
	} else if pot.max > 0 && totalBet < pot.max {
		pot.contributions[seat] = totalBet
		fmt.Printf("reduceMax\n")
		pot.reduceMax(totalBet)

	} else {
		pot.contributions[seat] += amount
		if isAllIn == true {
			fmt.Printf("allIn\n")
			pot.max = totalBet
		}
	}

	this.Players[seat].totalGameBet = totalBet
	this.PotValue += amount
	return this.CurrentBettingRound.addBet(seat, amount)

}

type Game interface {
	Play(TableResource)
}

func (this *GameState) newPlayer(seat *Seat) *Player {

	ss := new(Player)
	ss.Seat = seat
	ss.Seatnum = seat.Seatnum
	ss.UserId = seat.UserId
	//ss.User = seat.User
	//ss.Stack = seat.Stack
	ss.gameStack = seat.Stack
	ss.gameStartingStack = seat.Stack
	//ss.Channel = seat.Channel
	ss.game = this
	return ss

}

type Payout struct {
	Seat   int           `json:"seat"`
	Amount int           `json:"amount"`
	Card1  pokerlib.Card `json:"card1"`
	Card2  pokerlib.Card `json:"card2"`
}

/*
   For each seat participating in a game

   A stack amount is maintained for each seat at a table.   The TableStack.
   The TableStack is not effected until a game is Settled.

    When a new game starts each seats GameStack is equal to their TableStack,
     during the Game bets are deducted from the GameStack.     When the Game
      is settled the TableStack is set to the current GameStack value + any
      winnings from the game.
*/
type Player struct {
	//	Seat;
	Seatnum int
	UserId  int
	Seat    *Seat

	folded            bool
	allIn             bool
	leftTable         bool
	currentTurnBet    int
	cards             [2]pokerlib.Card
	totalGameBet      int
	gameStack         int
	gameStartingStack int

	percentWin float32
	percentTie float32

	bestHand  pokerlib.HandRank
	bestCards []pokerlib.Card
	game      *GameState
}

func (this *Player) calculateBestHand() ([]pokerlib.Card, pokerlib.HandRank) {

	cards := make([]pokerlib.Card, 7)
	for i, c := range this.game.communityCards {
		cards[i] = c
	}
	cards[5] = this.cards[0]
	cards[6] = this.cards[1]

	return pokerlib.Rank(cards)

}

/*
   There is always a current pot that bets go into.    When a player goes all in
      the current pot is closed at that amount and a new default pot is created
      for the remaining players bets, if any
*/

type AllInEvent struct {
	remainingSeats []int
	betAmount      int
}

func (this *GameState) IsSeatAllIn(seat int) bool {

	if this.Players[seat].gameStack == 0 {
		return true
	}
	return false

}
func (this *GameState) GetPot(amount int) *Pot {
	key, p := this.Pots.Floor(amount)
	if key != nil {
		return p.(*Pot)
	}
	return nil
}

func (this *GameState) NewPot(startAmount int) *Pot {
	color.Set(color.FgBlue)
	fmt.Printf("================NEWPOT------------ %d\n", startAmount)
	color.Unset()
	p := new(Pot)
	p.contributions = make(map[int]int, 0)
	p.number = this.potCount
	p.max = 0
	p.allInBy = -1
	p.game = this
	this.potCount += 1
	p.startAmount = startAmount
	this.Pots.Put(startAmount, p)
	return p
}

/*  A record of a completed game
    Sufficient to replay the game move by move
*/
type GameRecord struct {
	///	players map[int]
	cards [5]pokerlib.Card

	eventLog []Event
}

type HandResult struct {
	desc        string
	seat        int
	finishOrder int
}

func (this *GameState) determineWinners() (string, []int) {

	var highestHandValue pokerlib.HandRank = 0
	currentHighHandHolders := make([]int, 0)
	this.seatLock.Lock()
	defer this.seatLock.Unlock()

	current := this.ActivePlayers
	last := current
	for {
		curSeat := current.Value.(int)
		seatState := this.Players[curSeat]
		cards := make([]pokerlib.Card, 7)
		for i, c := range this.communityCards {
			cards[i] = c
		}
		cards[5] = seatState.cards[0]
		cards[6] = seatState.cards[1]

		bestCards, bestHand := pokerlib.Rank(cards)

		seatState.bestHand = bestHand
		seatState.bestCards = bestCards

		if uint32(bestHand) > uint32(highestHandValue) {
			currentHighHandHolders = make([]int, 0)
			currentHighHandHolders = append(currentHighHandHolders, curSeat)
			highestHandValue = bestHand
		} else if uint32(bestHand) == uint32(highestHandValue) {
			currentHighHandHolders = append(currentHighHandHolders, curSeat)
		}
		current = current.Next()
		if current == last {
			break
		}

	}

	winningHand := highestHandValue.Describe()
	var buf strings.Builder
	fmt.Printf("winning hand description: %s\n", winningHand)
	winners := len(currentHighHandHolders)
	if winners > 1 {
		fmt.Fprintf(&buf, "Players ")

		for i, s := range currentHighHandHolders {
			if i < winners {
				fmt.Fprintf(&buf, " %d and ", s)
			} else {
				fmt.Fprintf(&buf, " %d ", s)
			}
		}
		fmt.Fprintf(&buf, " split with a %s\n", highestHandValue.DescribeWithColor())

	} else {
		fmt.Fprintf(&buf, "Seat %d wins with %s\n",
			currentHighHandHolders[0],
			highestHandValue.DescribeWithColor())
	}
	return buf.String(), currentHighHandHolders

}

func (this *GameState) playerLeft(seatnum int) {

	this.seatLock.Lock()

	current := this.ActivePlayers
	first := current
	for {
		if current == nil {
			break
		}
		if current.Value == seatnum {
			current = current.Prev().Unlink(1)
			break
		}
		current = current.Next()
		if current == first {
			break
		}
	}
	delete(this.Players, seatnum)
}
