package domain

import (
	"container/ring"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/rbaderts/pokerlib"
)

type ActionType int

const (
	UnknownActionType          ActionType = 0
	BetAction                  ActionType = 1
	CheckAction                ActionType = 4
	DealHoleCardsAction        ActionType = 6
	DealCommonCardAction       ActionType = 7
	StartBettingRoundAction    ActionType = 8
	CompleteBettingRoundAction ActionType = 9
	SmallBlindAction           ActionType = 10
	BigBlindAction             ActionType = 11
	StartGameAction            ActionType = 12
	SettleBetsAction           ActionType = 13
	PlayCompleted              ActionType = 14
	ShowdownAction             ActionType = 15
)

func ActionTypeFromString(str string) ActionType {
	switch str {
	case "BetAction":
		return BetAction
	case "CheckAction":
		return CheckAction
	case "DealHoleCardsAction":
		return DealHoleCardsAction
	case "DealCommonCardAction":
		return DealCommonCardAction
	case "StartBettingRoundAction":
		return StartBettingRoundAction
	case "CompleteBettingRoundAction":
		return CompleteBettingRoundAction
	case "SmallBlindAction":
		return SmallBlindAction
	case "BigBlindAction":
		return BigBlindAction
	case "StartGameAction":
		return StartGameAction
	case "SettleBetsAction":
		return SettleBetsAction
	case "PlayCompleted":
		return PlayCompleted
	}
	return UnknownActionType
}

func (this ActionType) GetString() string {
	switch this {
	case BetAction:
		return "BetAction"
	case CheckAction:
		return "CheckAction"
	case DealHoleCardsAction:
		return "DealHoleCardsAction"
	case DealCommonCardAction:
		return "DealCommonCardAction"
	case StartBettingRoundAction:
		return "StartBettingRoundAction"
	case CompleteBettingRoundAction:
		return "CompleteBettingRoundAction"
	case SmallBlindAction:
		return "SmallBlindAction"
	case BigBlindAction:
		return "BigBlindAction"
	case StartGameAction:
		return "StartGameAction"
	case SettleBetsAction:
		return "SettleBetsAction"
	case PlayCompleted:
		return "PlayCompleted"

	}
	return "Unknown"

}


const (
	PreDeal = iota + 1 // hand has not yet started
	PostDeal
	PreFlop
	PostFlop
	PreTurn
	PostTurn
	PreRiver
	PostRiver
	Showdown
	Done // Winner(s) determined, pot(s) distributed
)

type GameStage int

type Action struct {
	Sequence int                    `json:"sequence"`
	Typ      string                 `json:"actionType"`
	Seat     int                    `json:"seat"`
	Params   map[string]interface{} `json:"params"`
	game     *GameState
	round    *Round
}

func (this Action) String() string {
	bytes, _ := json.Marshal(this)
	return string(bytes)
}

func (this *GameState) updateSeats() {

	for _, s := range this.Players {

		cards := make([]pokerlib.Card, 0)
		cards = append(cards, s.cards[0])
		cards = append(cards, s.cards[1])
		for _, c := range this.communityCards {
			cards = append(cards, c)
		}

		_, rank := pokerlib.Rank(cards)
		data := GameSeatData{rank.DescribeBasic()}
		GblMessageDispatcher.SendToSeat(this.table, s.Seatnum, GameSeatUpdated, data)

	}
}

func (this *GameState) updatePlayerOdds() {

	fmt.Printf("updatePlayerOdds\n")
	hands := make(map[int][2]pokerlib.Card, len(this.Players))

	for _, s := range this.Players {
		hands[s.Seatnum] = s.cards
	}

	odds := pokerlib.CalculateOdds(this.deck, hands, this.communityCards, 1000)
	data := OddsData{odds}

	GblMessageDispatcher.BroadcastToTable(this.table, TableOddsUpdated, &data)

}

func (this *GameState) DealHoleCards(seat int) {
	act := this.CreateAction(seat, DealHoleCardsAction)
	this.ProcessAction(this.CurrentBettingRound, act)
}

func (this *GameState) DealCommonCard() {
	act := this.CreateAction(0, DealCommonCardAction)
	this.ProcessAction(this.CurrentBettingRound, act)
}

func (this *GameState) PayBlinds(round *Round) {

	s := this.ActivePlayers.Next()
	//this.ActivePlayers = this.ActivePlayers.Next()
	smUser := s.Value.(int)
	//this.ActivePlayers = this.ActivePlayers.Next()
	act := this.CreateBetAction(smUser, this.BigBlind/2, 0, 0)
	this.ProcessBetAction(round, act)
	s = s.Next()

	bigUser := s.Value.(int)
	//this.ActivePlayers = this.ActivePlayers.Next()
	act = this.CreateBetAction(bigUser, this.BigBlind, 0, 0)
	this.ProcessBetAction(round, act)
	round.currentTurn = s.Next()

	//	act = this.CreateAction(0, BigBlindAction)
	//	act.Params["amount"] = this.BigBlind
	//	this.ProcessAction(act)
}

func (this *GameState) Replay(log *GameLog) {

}

func (this *GameState) GameFlow() GameStage {

	///	data := PlayerBetData{action.Seat, amount, 0, BigBlind}
	data := CalculateUpdate(this.table)
	GblMessageDispatcher.BroadcastToTable(this.table, TableGameStarted, &data)

	this.Stage = PreDeal

	r := this.ActivePlayers
	r = r.Next()
	last := this.ActivePlayers
	for {
		this.DealHoleCards(r.Value.(int))
		if r == last {
			break
		}
		r = r.Next()
	}
	/*
		this.seatLock.Lock()
		for _, s := range this.Players {
			this.DealHoleCards(s.Seatnum)
		}

	*/
	//this.seatLock.Unlock()

	go this.updatePlayerOdds()

	//data = CalculateUpdate(this.table)
	//GblMessageDispatcher.BroadcastToTable(this.table, TableStatusUpdated, &data)

	this.CurrentBettingRound = this.NewBettingRound(this.ActivePlayers)
	this.PayBlinds(this.CurrentBettingRound)

	foldWinner := -1
	for i := 0; i < 4; i++ {
		//		this.BettingRound()
		this.CurrentBettingRound.Execute()
		foldWinner = this.checkForFoldWinner()
		if foldWinner >= 0 {
			break
		}
		if i < 3 {
			if i == 0 {
				this.DealCommonCard()
				this.DealCommonCard()
			}
			this.DealCommonCard()
		}
		this.updateSeats()
		this.updatePlayerOdds()

		//data = CalculateUpdate(this.table)
		this.CurrentBettingRound = this.NewBettingRound(this.ActivePlayers.Next())
	}

	//	this.DeductCommitedCash()
	results := this.CalculateResults(foldWinner)
	fmt.Printf("Results = %v\n", results)
	this.Settle(results)

	return this.Stage
}

func (this *GameState) ProcessBetAction(round *Round, action *Action) BetType {

	var betType BetType

	//	this.actionHistory = append(this.actionHistory, action)

	isBigBlind := false
	isSmallBlind := false

	color.Set(color.FgRed)
	fmt.Printf("ProcessBetAction: %v, round: %d, turn: %d\n", action, round.round, round.turn)
	color.Unset()

	if round.round == 1 && round.turn == 1 {
		isSmallBlind = true
	} else if round.round == 1 && round.turn == 2 {
		isBigBlind = true
	}
	isAllIn := false
	if action.Params["isAllIn"].(int) != 0 {
		isAllIn = true
	}
	isFold := false
	if action.Params["isFold"].(int) != 0 {
		isFold = true
	}
	amount := action.Params["amount"].(int)
	if isBigBlind {
		if amount != this.BigBlind {
			fmt.Printf("Error: bigblind is not right amount")
		}
		this.addBet(action.Seat, amount, false)
		//		this.CurrentBettingRound.addBet(action.Seat, this.BigBlind)

		//		round.bets[action.Seat] += this.BigBlind
		data := PlayerBetData{action.Seat, amount, 0, BigBlind}
		GblMessageDispatcher.BroadcastToTable(this.table, GamePlayerActionBet, &data)
		betType = BigBlind
	} else if isSmallBlind {
		if amount != this.BigBlind/2 {
			fmt.Printf("Error: smallblind is not right amount")
		}
		this.addBet(action.Seat, amount, false)
		data := PlayerBetData{action.Seat, amount, 0, SmallBlind}
		GblMessageDispatcher.BroadcastToTable(this.table, GamePlayerActionBet, &data)
		betType = SmallBlind
	} else if isAllIn {

		totalbet := round.seatBets[action.Seat] + amount
		if totalbet > round.amount {
			round.last = round.currentTurn.Prev()
			round.amount = totalbet
		}
		this.Players[action.Seat].allIn = true
		raised := this.addBet(action.Seat, amount, true)
		if raised {
			round.last = round.currentTurn.Prev()
			//betType = Raise
		}
		data := PlayerBetData{action.Seat, amount, 0, AllIn}
		GblMessageDispatcher.BroadcastToTable(this.table, GamePlayerActionBet, &data)
		betType = AllIn

	} else if isFold {

		this.Players[action.Seat].folded = true

		data := PlayerBetData{action.Seat, 0, 0, Fold}
		GblMessageDispatcher.BroadcastToTable(this.table, GamePlayerActionFold, &data)
		fmt.Printf(" after fold: ")
	} else if amount == 0 {
		data := PlayerBetData{action.Seat, amount, 0, Check}
		GblMessageDispatcher.BroadcastToTable(this.table, GamePlayerActionBet, &data)

	} else {
		allIn := false
		if this.IsSeatAllIn(action.Seat) {
			allIn = true
		}
		raised := this.addBet(action.Seat, amount, allIn)
		fmt.Printf("addbet return raised = %v\n", raised)
		betType = Call
		if raised {
			//round.last = round.currentTurn.Prev()
			//round.amount = totalbet
			betType = Raise
		}
		data := PlayerBetData{action.Seat, amount, 0, betType}
		GblMessageDispatcher.BroadcastToTable(this.table, GamePlayerActionBet, &data)

	}

	round.currentTurn = round.currentTurn.Next()

	round.turn += 1
	return betType
}

func (this *GameState) ProcessAction(round *Round, action *Action) {
	color.Set(color.FgRed)
	fmt.Printf("ProcessAction: %s\n", action.Typ)
	color.Unset()
	this.actionHistory = append(this.actionHistory, ActionLogFromAction(action))

	switch ActionTypeFromString(action.Typ) {
	case DealHoleCardsAction:
		seat := action.Seat
		card1 := this.deck.DrawCard()
		card2 := this.deck.DrawCard()
		this.Players[seat].cards[0] = card1
		this.Players[seat].cards[1] = card2
		data := struct {
			Card1 pokerlib.Card `json:"card1"`
			Card2 pokerlib.Card `json:"card2"`
		}{
			card1,
			card2,
		}
		//		GblMessageDispatcher.BroadcastToTable(this.table, GameHoleCardsDraw, data)
		GblMessageDispatcher.SendToSeat(this.table, seat, GameHoleCardsDraw, data)

		break

	case DealCommonCardAction:
		card := this.deck.DrawCard()
		this.communityCards = append(this.communityCards, card)
		crds := []pokerlib.Card{card}
		GblMessageDispatcher.BroadcastToTable(this.table, GameCardsDealt, crds)
		break

	case StartBettingRoundAction:
		round := this.NewBettingRound(this.ActivePlayers)
		this.CurrentBettingRound = round
		round.Execute()
		break

	case CompleteBettingRoundAction:
		this.CurrentBettingRound = nil
		break

	case SmallBlindAction:
	case BigBlindAction:

	}

}

var SeqCounter int = 1

func (this *GameState) CreateBetAction(seat int, amount int, isAllIn int, isFold int) *Action {
	act := this.CreateAction(seat, BetAction)
	act.Params["amount"] = amount
	act.Params["isAllIn"] = isAllIn
	act.Params["isFold"] = isFold
	act.round = this.CurrentBettingRound
	return act
}
func (this *GameState) CreateAction(seat int, typ ActionType) *Action {
	act := new(Action)
	act.Sequence = SeqCounter
	SeqCounter += 1
	act.Seat = seat
	act.Typ = typ.GetString()
	act.game = this
	act.Params = make(map[string]interface{}, 0)

	return act
}

func (this *GameState) BettingRound() {

	this.CurrentBettingRound = this.NewBettingRound(this.ActivePlayers)
	this.CurrentBettingRound.Execute()
}

func (this *GameState) NewBettingRound(seats *ring.Ring) *Round {

	//	blinds := false
	this.bettingRoundCounter += 1
	//	if this.bettingRoundCounter == 1 {
	//		blinds = true
	//	}
	b := new(Round)
	b.round = this.bettingRoundCounter
	b.seatBets = make(map[int]int, 0)
	b.amount = 0
	b.game = this

	b.currentTurn = seats
	b.turn = 1
	return b
}

/*
func (this *GameState) DeductCommitedCash() {

	this.seatLock.Lock()
	defer this.seatLock.Unlock()
	for _, seat := range this.Players {
		seat.Stack -= seat.totalGameBet
	}
}
*/

func (this *GameState) Settle(results *GameResults) {
	ctx := context.Background()

	/*
		tx, err := adapters.DB.Begin(ctx)
		if err != nil {
			fmt.Errorf("Error create Tx: %v\n", err)
		}
	//	txctx := context.WithValue(ctx, "TX", tx)

		defer tx.Rollback(context.Background())
	*/

	// First update the players table stack to reflect bets made in the game

	for _, seat := range this.Players {
		spent := seat.gameStartingStack - seat.gameStack
		seat.Seat.Stack -= spent
		fmt.Printf("Settle: seat %d, spent=%d\n", seat.Seatnum, spent)
	}

	// Now distribute the pots to the winners
	msg := ""
	for allIn, pot := range results.potResults {

		//Payouts := make([]Payout, 0)
		msg = pot.message
		fmt.Printf("Settle pot %d:  %d particpants, %d winners\n", allIn, len(pot.participants), len(pot.winners))
		total := pot.amount
		for _, s := range pot.participants {
			portion := 0
			seat := this.Players[s]
			if containsValue(pot.winners, s) {
				portion = total / len(pot.winners)
				fmt.Printf("Seat %d get %d from pot, tableStack now is %d\n", s, portion, seat.Seat.Stack+portion)
				seat.Seat.Stack += portion
			}
			//fmt.Printf("distributing %d dollars to seat %d (userid: %v)\n", portion, seat.Seatnum, seat.UserId)
			//seat.Stack += portion
			//Payouts = append(Payouts, Payout{s, portion, seat.cards[0], seat.cards[1]})
			//data := GameShowdownData{pot.message, Payouts}
			//_ = GblMessageDispatcher.BroadcastToTable(this.table, GameShowdown, data)
		}
	}

	Payouts := make([]Payout, 0)
	for _, s := range this.Players {
		Payouts = append(Payouts, Payout{s.Seatnum, s.Seat.Stack, s.cards[0], s.cards[1]})
	}
	data := GameShowdownData{msg, Payouts}
	_ = GblMessageDispatcher.BroadcastToTable(this.table, GameShowdown, data)
	GblTableRepository.UpdateAllPlayerStates(ctx, this.table)

}

func (this *GameState) checkForFoldWinner() int {

	lastUnfoldedPlayer := -1
	unfoldedPlayers := 0

	this.seatLock.Lock()
	defer this.seatLock.Unlock()
	s := this.ActivePlayers
	last := s
	for {
		seat := this.Players[s.Value.(int)]
		if seat.folded != true {
			unfoldedPlayers += 1
			lastUnfoldedPlayer = seat.Seatnum
		}
		s = s.Next()
		if s == last {
			break
		}
	}
	if unfoldedPlayers == 1 {
		return lastUnfoldedPlayer
	}

	return -1

}

func containsValue(ids []int, id int) bool {
	for _, v := range ids {
		if id == v {
			return true
		}
	}
	return false
}
