package app

import (
	"context"
	"fmt"
	"github.com/rbaderts/poker/domain"
)

/*jjjjj
type TableRepository interface {
	GetAllTables()  []domain.Table
	AddTable(table domain.Table)

}

*/

type TableService interface {
    NewTable(seats int) *domain.Table
	NewReplayTable(log *domain.GameLog) *domain.Table
	GetTable(id int) *domain.Table
    ListTables() ([]*domain.Table)
}

type TableServiceState struct {
}

var GblTableService TableService

func NewTableService() TableService {
	s := new(TableServiceState)
	return s
}


func (this *TableServiceState) NewTable(seats int) *domain.Table {

	fmt.Printf("App New table for %d\n", seats)
	//ctx := context.Background()
//	table := domain.GblTableRepository.NewTable(ctx, "test", 2, 6)
	table := domain.NewTable("test", 2, 6)
	fmt.Printf("NewTable: id = %s\n", table.Id)
	go table.TableLoop()
	return table;
}

func (this *TableServiceState) NewReplayTable(log *domain.GameLog) *domain.Table {

	fmt.Printf("App New table for %d\n", len(log.Game.Seats))
//	ctx := context.Background()
//	table := domain.GblTableRepository.NewTable(ctx, "test", log.Game.BigBlind, len(log.Game.Seats))
//	fmt.Printf("NewReplayTable: id = %s\n", table.Id)
/*

	activeSeats  := make([]int, len(log.Game.Seats))
	for _, seat := range log.Game.Seats {
		user, err := domain.GblUserRepository.LoadUserById(ctx, seat.UserId)
		if err != nil {
			fmt.Errorf("%v\n", err)
			continue
		}
		table.AssignSeat(context.Background(), user, seat.Stack, "")
		activeSeats = append(activeSeats, seat.Seatnum)

	}

	*/
	//table.:= domain.NewGame(table, activeSeats, log.Game.Button)

//	table.ReplayGame(log)
	//go table.TableLoop()
	return nil;
}


func (this *TableServiceState) ListTables() ([]*domain.Table) {
	ctx := context.Background()
	res, err := domain.GblTableRepository.GetAllTables(ctx)
	if err != nil {
		return nil
	}
	return res

}

func (this *TableServiceState) GetTable(id int) *domain.Table {
	ctx := context.Background()
	table, err :=  domain.GblTableRepository.GetTable(ctx, id)
	if err != nil {
		return nil
	}
	return table
}
