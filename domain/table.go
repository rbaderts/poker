package domain

import (
	"container/ring"
	"context"
	"errors"
	"fmt"
	"github.com/rbaderts/poker/db"
	"os"
	"strconv"
	"time"
)

/*
   A table represents a game table.    At table has a fixed # of seats.

   At any given time each Seat may be empty or populated by a player
   A table continually plays games as long a 2 or more seats are occupied.






*/

/*
    A record of how much money the player has won/lost while sitting at the table.
 */
type TablePlayerState struct {
	TableId        int          `json:"userId" db:"table_id"`
	UserId         int          `json:"userId" db:"user_id"`
	SeatNum        int          `json:"seatnum" db:"seatnum"`
	StartingStack  int          `json:"startingStack" db:"starting_stack"`
	StackChange    int          `json:"stackChange" db:"stack_change"`

}

type TableRepository interface {
	GetAllTables(context.Context) ([]*Table, error)
	UpdateTable(context.Context, *Table) error
	AddTable(context.Context, *Table) (*Table, error)

	GetTable(context.Context, int) (*Table, error)
//	NewTable(ctx context.Context, name string, blinds int, numSeats int) (*Table, error)

	UpdatePlayerState(ctx context.Context, table *Table, seatnum int) error
	UpdateAllPlayerStates(ctx context.Context, table *Table) error
	GetAllPlayerStates(ctx context.Context, tableId int) ([]*TablePlayerState, error)
	ClearPlayerStates(ctx context.Context, tableId int) error
}

type Seat struct {
	Seatnum            int
	UserId             int
	SittingOut         bool
	StartingStack      int
	Stack              int
	Channel            string
	User               *User
}

func (this *Seat) clearSeat() {
	this.UserId = -1
	this.User = nil
	this.SittingOut = false
	this.Stack = 0
	this.StartingStack = 0
	this.Seatnum = -1

}

func (this *Seat) GetUser() *User {
	if this.User == nil {
		if this.UserId != -1 {
			ctx := context.Background()
			var err error
			this.User, err = GblUserRepository.LoadUserById(ctx, this.UserId)
			if err != nil {
				this.User = nil
				return nil
			}
		}
	}
	return this.User

}

func (this *Table) NewSeat(seatnum int) *Seat {
	seat := new(Seat)
	seat.Channel = "t_" + strconv.Itoa(this.Id) + "_s_" + strconv.Itoa(seatnum)
	seat.Seatnum = seatnum
	seat.SittingOut = false
	seat.UserId = -1
	return seat
}

type Table struct {
	Id           int     `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	Channel      string  `json:"channel" db:"channel"`
	Button       int     `json:"button" db:"button"`
	BigBlind     int     `json:"bigblind" db:"bigblind"`
	NumSeats     int     `json:"numseats" db:"num_seats"`
	Seats        []*Seat `json:"seats"`
	currentGame  *GameState
	EventChannel *string
}

func NewTable(name string, blind int, seats int) *Table {
	t := new(Table)
	t.BigBlind = blind
	t.Name = name
	t.Seats = make([]*Seat, 0)
	t.NumSeats = seats
	t.Button = 0
	ctx := context.Background()
	t, err := GblTableRepository.AddTable(ctx, t)
	if err != nil {
		return nil
	}
	for i := 0; i < seats; i++ {
		t.Seats = append(t.Seats, t.NewSeat(i))
	}
	t.Channel = "t_" + strconv.Itoa(t.Id)
	GblTableRepository.UpdateTable(ctx, t)

	err = GblRecoveryRepository.CreateTableTombstone(ctx, t.Id)
	if err != nil {
		fmt.Printf("Error  creating table tombstone: %v\n", err)
	}


	return t

}

/*
func (this *Table) ReplayGame(log *GameLog) {

	ctx := context.Background()
	activeSeats := make([]int, len(log.Game.Players))

	for _, seat := range log.Game.Players {
		user, err := GblUserRepository.LoadUserById(ctx, seat.UserId)
		_ = user
		if err != nil {
			fmt.Errorf("%v\n", err)
			continue
		}
		//this.AssignSeat(user, seat.Stack, "")
		activeSeats = append(activeSeats, seat.Seatnum)

	}

	this.currentGame = NewGame(this, activeSeats, log.Game.Button)

	/*
		for _, action := range log.Actions {

		}
	*/
/*
}
*/

func (this *Table) VacateSeat(ctx context.Context, seatnum int) {

	tx, err := GblDBPool.Begin(ctx)
	if err != nil {
		fmt.Errorf("Error creating tx: %v\n", err)
	}
	defer tx.Commit(ctx)


	txctx := context.WithValue(ctx, db.TX_KEY, &tx)

	ontable := this.Seats[seatnum].Stack

	// return money to user account
	GblUserRepository.AddMoney(txctx, this.Seats[seatnum].UserId, ontable)
	this.Seats[seatnum].User = nil
	this.Seats[seatnum].UserId = -1
	this.Seats[seatnum].Stack = 0
	this.Seats[seatnum].SittingOut = false
	GblTransactionRepository.RecordTransaction(txctx, TableToUserAccount, ontable, this.Seats[seatnum].UserId, this.Id)
	GblTableRepository.UpdatePlayerState(txctx, this, seatnum)

}

var MIN_STAKE int = 100

func (this *Table) AssignSeat(ctx context.Context, user *User, stack int, sessionId string) (int, error) {

	if stack < MIN_STAKE {
		return 0, errors.New("Not enough money to join table")

	}
	/*
	tx, err := GblDBPool.Begin(ctx)
	if err != nil {
		fmt.Errorf("Error creating tx: %v\n", err)
	}
	txctx := db.NewContextWithTx(ctx, &tx)
	 */

	seatNum := -1
	for _, s := range this.Seats {
		if s.UserId == -1 {
			s.UserId = user.Id
			//s.occupantUserId = User.Id
			s.SittingOut = false
			s.Stack = stack
			s.StartingStack = stack
			// Take money from user account and put on table
			GblUserRepository.DeductMoney(ctx, s.UserId, stack)

			s.User = user
			data := PlayerJoinedData{user.GivenName, s.Seatnum, user.Stack}
			///	this.BroadcastEvent(*msg)
			_ = GblMessageDispatcher.BroadcastToTable(this, TablePlayerJoined, &data)

			GblTransactionRepository.RecordTransaction(ctx, UserAccountToTable, stack, user.Id, this.Id)

			seatNum = s.Seatnum
			break
		}
	}

	GblTableRepository.UpdatePlayerState(ctx, this, seatNum)

	/*
	if seatNum == -1 {
		err := tx.Rollback(ctx)
		if err != nil {
			return -1, err
		}
		return -1, errors.New("No seats available")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}
	 */
	return seatNum, nil
}

func (this *Table) FindUsersSeat(userId int) int {

	for _, s := range this.Seats {
		if s.UserId == userId {
			return s.Seatnum
		}
	}
	return -1

}

func (this *Table) GetActivePlayers() int {

	count := 0
	for _, s := range this.Seats {
		if s.UserId != -1 && !s.SittingOut {
			count += 1
		}
	}
	return count
}

func (this *Table) moveButton() {
	this.Button = this.seatAfter(this.Button)

}

func (this *Table) seatAfter(seat int) int {

	//	if seat == 0 {
	//		return 1
	//	}
	for i := seat+1; i < len(this.Seats); i++ { //seat + 1; i < len(this.Players); i++ {
//	for i := seat + 1; i < len(this.Players); i++ {
        if this.Seats[i].UserId != -1 && this.Seats[i].SittingOut == false {
        	return i
        }
	}

	for i := 0; i < seat; i++ {
		if this.Seats[i].UserId != -1 {
			return i
		}
	}

	fmt.Errorf("Unable to move the button\n")
	os.Exit(-1)
	return -1

}

func (this *Table) TableLoop() {

	//	r := make(chan int)
	for {

		playersReady := this.GetActivePlayers()
		fmt.Printf("TableLoop - players Ready = %d\n", playersReady)
		if this.currentGame != nil && playersReady >= 2 {

			// update := CalculateUpdate(this)

			//this.Print()
			fmt.Printf("Starting game\n")
			//this.currentGame.GameFlow(*update)
			this.currentGame.GameFlow()

			this.currentGame.dumpGame("GameLog-" + this.currentGame.Id)
			this.currentGame = nil

		}
		time.Sleep(6 * time.Second)
	}

}

func (this *Table) SetupGame() {

	var players *ring.Ring
	tmp := make([]int, 0)
	for _, seat := range this.Seats {
		if seat != nil && !seat.SittingOut && seat.UserId > 0 {
			next := ring.New(1)
			next.Value = seat.Seatnum
			if players != nil {
				players.Link(next)
	 		} else {
	 			players = next
			}
			//tmp = append(tmp, seat.Seatnum)
		}
	}
	fmt.Printf("seats = %v\n",tmp)

	/*
		activeSeats := make([]int, 0)
		for i := this.Button; i < len(tmp); i++ {
			if !this.Players[i].SittingOut && this.Players[i].UserId > 0 {
				activeSeats = append(activeSeats, i)
			}
		}
		for i := 0; i < this.Button; i++ {
			if !this.Players[i].SittingOut && this.Players[i].UserId > 0 {
				activeSeats = append(activeSeats, i)
			}
		}
	*/

	//	this.Button = this.getNextDealer()

	this.moveButton()

	/*
	activeSeats := make([]int, 0)
	for i, s := range tmp {
		if !this.Players[s].SittingOut && this.Players[i].UserId > 0 {
			activeSeats = append(activeSeats, this.Players[i].Seatnum)
		}
	}
	 */
	this.currentGame = NewGame(this, players, this.Button)

}

func (this *Table) GetSeat(seatNum int) *Seat {
	return this.Seats[seatNum]
}

func (this *Table) RemovePlayer(seatNum int, sessionId string) error {

	for _, s := range this.Seats {
		if s.Seatnum == seatNum {
			s.clearSeat()

			data := PlayerLeftData{s.Seatnum}
			///	this.BroadcastEvent(*msg)
			_ = GblMessageDispatcher.BroadcastToTable(this, TablePlayerLeft, &data)

			return nil
		}
	}

	GblTableRepository.UpdatePlayerState(context.Background(), this, seatNum)

	//	token := auth.GetJWTToken(r)

	return nil

}

func (this *Table) getNextDealer() int {

	next := false
	for _, seat := range this.Seats {
		if next == true {
			return seat.Seatnum
		}
		if this.Button == -1 {
			return seat.Seatnum
		}
		if this.Button == seat.Seatnum {
			next = true
		}
	}
	return -1
}

func (this *Table) ReturnPlayerCash() {

	ctx := context.Background()
	for _, s := range this.Seats {
		GblUserRepository.AddMoney(ctx, s.UserId, s.Stack)
		GblTransactionRepository.RecordTransaction(ctx, TableToUserAccount, s.Stack, s.UserId, this.Id)
	}
}
