package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rbaderts/poker/domain"
	"strconv"
)

type MemoryTableRepository struct {
	BaseTableRepository
	tables  map[int]*domain.Table
	counter int
}

func NewMemoryTableRepository() domain.TableRepository {
	r := new(MemoryTableRepository)
	r.tables = make(map[int]*domain.Table)
	return r
}

func (this MemoryTableRepository) getNextId() int {
	highest := 0
	for i, _ := range this.tables {
		if i > highest {
			highest = i
		}
	}
	return highest + 1
}

func (this *MemoryTableRepository) UpdateTable(ctx context.Context, table *domain.Table) error {
	this.tables[table.Id] = table
	return nil
}
func (this *MemoryTableRepository) AddTable(ctx context.Context, table *domain.Table) (*domain.Table, error) {

	id := this.getNextId()
	table.Id = id
	this.tables[id] = table

	return table, nil

}

func (this *MemoryTableRepository) GetAllTables(ctx context.Context) ([]*domain.Table, error) {

	r := make([]*domain.Table, len(this.tables))

	i := 0
	for _, v := range this.tables {
		r[i] = v
		i += 1
	}

	return r, nil
}

func (this *MemoryTableRepository) GetTable(ctx context.Context, id int) (*domain.Table, error) {
	k, ok := this.tables[id]

	if !ok {
		return nil, errors.New("Table Not found")
	}
	return k, nil
}

func (this *MemoryTableRepository) NewTable(ctx context.Context, name string, blinds int, numSeats int) (*domain.Table, error) {
	fmt.Printf("MemoryTableRepository: NewTable\n")
	table := new(domain.Table)
	//id, err := uuid.NewV4()
	//if err != nil {
	//	log.Fatal(err)
	//}
	table.Id = this.getNextId()
	table.Name = "Table_" + strconv.Itoa(table.Id)

	table.Seats = make([]*domain.Seat, numSeats)
	table.BigBlind = blinds
	table.Button = 0
	for i := 0; i < numSeats; i++ {
		channel := "t_" + strconv.Itoa(table.Id) + "_s_" + strconv.Itoa(i)
		table.Seats[i] = &domain.Seat{i, -1, false, 0, 0, channel, nil}
	}

	table.Channel = "t_" + strconv.Itoa(table.Id)

	data, err := json.Marshal(table)

	//	table.blockedRequestChannels = make(map[int64]*chan Message)

	/*
		res, err := CentClient.Presence(ctx, table.Channel)
		if err != nil {
			fmt.Errorf("Error getting channel presense: %v\n", err)

		}

		fmt.Printf("channel presence: %v\n", res)


	*/
	ctx = context.Background()
	//table *domain.Table, typ domain.EventType, data interface{}) error {

	err = domain.GblMessageDispatcher.PublishToTable(table, data)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	//	err = CentClient.Publish(ctx, table.Channel, data)

	this.tables[table.Id] = table

	return table, nil
}
