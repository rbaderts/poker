package domain

import (
	"fmt"
	"github.com/rbaderts/pokerlib"
	"strings"
)

type GameResults struct {
	seatResults map[int]*SeatResult     // map of seat # to SeatResult
	orderIndex []*SeatResult            // ordered from highest to lowest rank
	potResults []*PotResult            // ordered from highest to lowest rank
}

type SeatResult struct {
	seat        int
	handValue   pokerlib.HandRank
	description string
	cards       [2]pokerlib.Card
}

type PotResult struct {
	amount          int
	winners         []int
	participants    []int
	message         string
}


func (this *GameResults) String() string {
	var buf strings.Builder
	for _, s := range this.seatResults {
		buf.WriteString(fmt.Sprintf("seat %d:   handValue (%v)\n", s.seat, s.handValue))
	}
	for i, p := range this.potResults {
		buf.WriteString(fmt.Sprintf("pot %d:   winners (%v), message: %s\n", i, p.winners, p.message))
	}

	return buf.String()

}

func (this *GameState) CalculateResults(foldWinner int) *GameResults {

	results := new(GameResults)
	results.seatResults = make (map[int]*SeatResult)
	results.potResults = make ([]*PotResult, 0)

	if (foldWinner >= 0) {
		result := new(SeatResult)
		result.seat = foldWinner
		result.description = "By Default"
		results.seatResults[foldWinner] = result

	} else {
		for _, s := range this.Players {
			if !s.folded {
				cards, rank := s.calculateBestHand()

				fmt.Printf("best hand for seat %d has communityCards %v and is %v\n", s.Seatnum, cards,
					rank.DescribeWithColor())

				result := new(SeatResult)
				result.seat = s.Seatnum
				result.handValue = rank
				result.description = rank.Describe()
				result.cards[0] = s.cards[0]
      			result.cards[1] = s.cards[1]
				results.seatResults[s.Seatnum] = result
			}
		}
	}
    it := this.Pots.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		pot := value.(*Pot)
		pr := new(PotResult)
		if foldWinner >= 0 {
			//sr := results.seatResults[foldWinner]
			pr.message = fmt.Sprintf("Seat %d Wins", foldWinner)
			pr.amount = pot.GetValue()
			pr.winners = []int{foldWinner}
			pr.participants = []int{foldWinner}
			results.potResults = append(results.potResults, pr)

		} else {
			participants := pot.GetNonFoldedParticpants()
			pr.winners = results.getWinners(participants)
			pr.participants = participants
			sr := results.seatResults[pr.winners[0]]
			pr.message = results.Describe(pr.winners, sr.description)
			pr.amount = pot.GetValue()
			results.potResults = append(results.potResults, pr)
		}

	}

	return results

}


// Given a set of seats, determine and return the winner (or winners if split)
func (this *GameResults) getWinners(seats []int)  []int {

	r := make([]int, 0)

	var highHand uint32 = 0
	for _, snum := range seats {
		seat, ok := this.seatResults[snum]
		if !ok {
			continue
		}
		if uint32(seat.handValue) > highHand {
			r = make([]int, 0)
			r = append(r, snum)
			highHand = uint32(seat.handValue)
		} else if uint32(seat.handValue) == highHand {
			r = append(r, snum)
		}
	}

	fmt.Printf("Winning seats: %v\n", r)
	return r

}

func  (this *GameResults) Describe(winners []int, handDescription string) string {

	var buf strings.Builder

	if len(winners) > 1 {
		fmt.Fprintf(&buf, "Players ")

		for i, s := range winners {
			if i < len(winners) {
				fmt.Fprintf(&buf, " %d and ", s)
			} else {
				fmt.Fprintf(&buf, " %d ", s)
			}
		}
		fmt.Fprintf(&buf, " split with a %s\n",  handDescription)

	} else {
		fmt.Fprintf(&buf, "Seat %d wins with %s\n", winners[0], handDescription)
	}
	return buf.String()

	//	this.Payout(currentHighHandHolders, message)

	//	current = this.ActivePlayers

}
