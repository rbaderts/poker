package domain

import "math"

type Pot struct {
	number             int
	participatingSeats []int
	//amount             int
	startAmount        int
	max                int
	allInBy            int
	contributions      map[int]int      // for each seat how much
	game               *GameState
}

func (this *Pot) GetNonFoldedParticpants() []int {
	r := make([]int, 0)
	for s, _ := range this.contributions {
		if (this.game.Players[s].folded == false) {
			r = append(r, s)
		}
	}
	return r
}

func (this *Pot) GetValue() int {
	total := 0
	for _, v := range this.contributions {
		total += v
	}
	return total


}


func (this* Pot) reduceMax(newmax int) {

	oldMaxPot := this.game.GetPot(math.MaxInt32)
	newPot := this.game.NewPot(newmax)
	for s, amount := range oldMaxPot.contributions {
		if amount + oldMaxPot.startAmount > newmax {
			newPot.contributions[s] = amount - newmax
			this.contributions[s] = newmax
		}
	}
}
