package domain

import (
	"container/ring"
	"fmt"
	"time"
	"github.com/fatih/color"
)
//

type Bet struct {
	Amount  int    `json:"amount"`
	IsAllIn int    `json:"isAllIn"`
	IsFold  int    `json:"isFold"`
}

/*
type Bet struct {
	betType BetType
	Amount  int
}
 */

type BetResponseHandler struct {
	seat  *Player
	round *Round
}

type Round struct {
	round              int
	seatBets           map[int]int
	amount             int
	game               *GameState
	currentTurn        *ring.Ring
	last               *ring.Ring
	turn               int
}

type PlayerBet struct {
	betType BetType
	Amount  int
}


func (this *Round) playersNotFolded() int {

	start := this.currentTurn
	cur := this.currentTurn
	count := 0
	for {
		seat := this.game.Players[cur.Value.(int)]
		if !seat.folded {
			count += 1
		}
		cur = cur.Next()
		if cur == start {
			break
		}
	}
	return count

}

func (this *Round) playersLeftBetting() int {
	start := this.currentTurn
	cur := this.currentTurn
	count := 0
	for {
		seat := this.game.Players[cur.Value.(int)]
		if !seat.folded && !seat.allIn {
			count += 1
		}
		cur = cur.Next()
		if cur == start {
			break
		}
    }
    return count
}

/*
   Returns false game is forfeit
 */
func (this *Round) Execute() {

	if this.MoreBetting() == false {
		time.Sleep(2 * time.Second)
		return
	}


	color.Set(color.FgYellow)
	fmt.Printf("============= starting betting round %d with seat: %v =============\n",
		this.round, this.currentTurn.Value.(int))
	color.Unset() // Don't forget to unset

	this.last = this.currentTurn.Prev()
	fmt.Printf(" last bet will be with seat %d\n", this.last.Value.(int))

//	done := false

	oneMore := false
	for {

		p := this.playersNotFolded()
		if  p == 1 {
			break
		}

		color.Set(color.FgRed)
		fmt.Printf("Round loop:  Players not folded, total ring size: %d\n", p, this.currentTurn.Len())
		color.Unset()
		/*
		if this.currentTurn.Len() < 1 {
			break;
		}
		 */
		seatNum := this.currentTurn.Value.(int)

		if (this.game.Players[seatNum].folded || this.game.Players[seatNum].allIn) {
			this.currentTurn = this.currentTurn.Next()
			continue
		}

		betAction, err := this.getPlayerBet(seatNum)
		if err != nil {
			fmt.Printf("Error getting player bet: %v\n", err)
		}

		prev := this.currentTurn.Prev()
		betType := this.game.ProcessBetAction(this, betAction)

		if betType == Raise {
			color.Set(color.FgBlue)
			fmt.Printf("action raised the Bet")
			color.Unset()
			this.last = prev
			oneMore = false
		}
		if oneMore == true {
			break
		}
		if this.currentTurn == this.last {
			oneMore = true
		}
	}

	color.Set(color.FgRed)
	fmt.Printf("Completed Betting Round: %d\n", this.round)
	fmt.Printf("----------------------------\n")
	color.Unset()
	color.Set(color.FgGreen)
	for sn, bt := range this.seatBets {
		fmt.Printf("Seat # %d in for %d\n", sn, bt)
	}
	color.Unset()

	//this.game.GetCurrentPot().amount += this.GetTotalAmount()

	data := BettingRoundCompleteData{this.GetTotalAmount()}
	_ = GblMessageDispatcher.BroadcastToTable(this.game.table, BettingRoundComplete, &data)

}



func (this *Round) getPlayerBet(
	seatNum int) (*Action, error) {


	color.Set(color.FgMagenta)
	fmt.Printf("getPlayerBet for %d\n", seatNum)
    fmt.Printf("round %d current bet amount is %d\n", this.round, this.amount)
	color.Unset()

	currentBet := this.amount

	seatState := this.game.Players[seatNum]

	if (seatState.allIn == true || seatState.folded == true) {
		fmt.Printf("ERROR: seat is allin or folded alread\n")
		return nil, nil
	}
	color.Set(color.FgCyan)
	fmt.Printf("    Querying seat %d for a bet\n", seatNum)
	playerBet := this.seatBets[seatNum]
	//playerBet := this.game.Players[seatNum].totalGameBet

	fmt.Printf("CurrentbetAmount: %d, currentPlayerBet: %d\n",
		currentBet, playerBet)
	color.Unset()
	bet, err := this.RequestPlayerBet(
		playerBet, seatState)

	if err != nil {
		fmt.Errorf("Error getting player bet: %v\n", err)
		if currentBet > playerBet {
			bet = &Bet{0, 0, 1}
		} else {
			bet = &Bet{0, 0, 0}
		}
		//return nil, err
	}

	fmt.Printf("    %s of %d for seat: %d\n", Yellow("Received a bet"),
		bet.Amount, seatNum)

	act := this.game.CreateBetAction(seatNum, bet.Amount, bet.IsAllIn, bet.IsFold)
	return act, nil
//	act.Params["amount"] = bet.Amount

}

func (this *Round) addBet(seat int, amount int) bool {
	color.Set(color.FgGreen)
	fmt.Printf("addBet for seat %d, amount: %d\n",seat, amount)
	color.Unset() // Don't forget to unset

	fmt.Printf("  total bets for seat %d so far this round:  %d\n", seat, this.seatBets[seat])
	fmt.Printf("  total bets for seat %d so far this game:  %d\n", seat, this.game.Players[seat].totalGameBet)
	val := this.seatBets[seat]
	total := val + amount
    this.seatBets[seat] = total

	fmt.Printf("  total bets for seat %d now %d\n", seat, this.seatBets[seat])

	if total > this.amount {
		this.amount = total
		fmt.Printf("  round high bet now %d\n", this.amount)
		return true
	}
	return false

}
func (this *Round) RequestPlayerBet(
	currentPlayerBet int,
	seat *Player) (*Bet, error) {

	//fmt.Printf("GetPlayerBet: currentBet: %d, currentRoundTotal: %d\n",
	//	currentBet, currentRoundTotal)

	currentBet := this.amount
	callAmount := currentBet - currentPlayerBet

	data := BetRequestData{currentBet, callAmount}
	///cmd := NewEvent(BetRequest, data)
	//	cmd := NewEvent(this.game.tableId, this.game.Id, BetRequest, data, 0)

	h := new(BetResponseHandler)
	h.seat = seat
	h.round = this

	betData, err := GblMessageDispatcher.SendRequestAndWait(this.game.table, seat.Seatnum, BetRequest, &data, h)

	//		SendRequestAndWait(table *Table, Seatnum int, typ EventType, data interface{}, handler ResponseHandler) (interface{}, error)


	if err != nil {
		fmt.Errorf("Error get bet from client: %v\n", err)
		return nil, err
	}

	bet := betData.(Bet)
	return &bet, err
}


type BetRespHandler struct {
	seat  *Player
	round *Round
}

func (this *Round) GetTotalAmount() int {

	total := 0
	for _, b := range this.seatBets {
		total += b
	}

	return total
}

func (this *Round) MoreBetting() bool {

	activePlayers := 0
	for _, state := range this.game.Players {

		if state.folded {
			continue
		}
		if state.allIn {
			continue
		}
		activePlayers += 1
	}

	if activePlayers > 1 {
		return true
	}
	return false
}

func (this *BetResponseHandler) Handle(cmd *Message) (interface{}, error) {

	var amount int = 0
	data := cmd.Data.(map[string]interface{})
	//betType := int(data["betType"].(float64))
	isAllIn := int(data["isAllIn"].(float64))
	isFold := int(data["isFold"].(float64))
	amtData := data["amount"]
	if amtData != nil {
		amount = int(amtData.(float64))
	}

	bet := Bet{amount, isAllIn, isFold}

	fmt.Printf("Recieved bet:  amount: %d, isAllIn: %d, isFold: %d\n", amount, isAllIn, isFold)
	/*
	   if cmd.ActionType == GamePlayerActionFold {
	       amount = -1
	   } else {
	       data := cmd.Data.(map[string]interface{})
	       amount = int(data["amount"].(float64))
	   }
	*/

	//      amount := res.Data.(BetResponse).Amount
	//      amountStr := res.Data.(BetData)
	///     amountStr := res.Params[AMOUNT]
	//      amount, _ := strconv.Atoi(amountStr)

	//  _ = this.round.PlayerBet(int(amount), this.seat.seat)

	return bet, nil


}

