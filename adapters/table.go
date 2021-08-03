package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/randallmlough/pgxscan"
	"github.com/rbaderts/poker/db"
	"github.com/rbaderts/poker/domain"
	"log"
	"strconv"
)

type BaseTableRepository struct {

}
type PostgresTableRepository struct {
	BaseTableRepository
	counter int
}

func NewPostgresTableRepository() domain.TableRepository {
	r := new(PostgresTableRepository)
	return r
}


func (this *PostgresTableRepository) AddTable(ctx context.Context, table *domain.Table) (*domain.Table, error) {
	fmt.Printf("AddTable\n")
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Conn().Query(ctx,
		"insert into Table(button, bigblind, num_seats, channel) values($1,$2,$3,$4) " +
			" returning id, button, bigblind, num_seats, channel;",
		table.Button, table.BigBlind, table.NumSeats, table.Channel)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error adding user\n"))
	}
	var tableRec domain.Table

	//	if err := pgxscan.ScanOne(&user, rows); err != nil {
	if err := pgxscan.NewScanner(rows).Scan(&tableRec); err != nil {
		log.Fatalf("Insert User failed: %v", err)
		return nil, err

	}
	rows.Close()
	return &tableRec, nil

}

func (this *PostgresTableRepository) GetTable(ctx context.Context, id int) (*domain.Table, error) {

	fmt.Printf("GetTable: %d\n", id)
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	rows, err := conn.Conn().Query(newctx,
		`SELECT * from tables where id = $1`, id)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error loading table with id %d\n", id))
	}
	var table domain.Table
	if err := pgxscan.NewScanner(rows).Scan(&table); err != nil {
		fmt.Errorf("GetTable error: %v\n", err)
		return nil, err
	}
	rows.Close()

	return &table, nil
}

func (this *PostgresTableRepository) GetAllTables(ctx context.Context) ([]*domain.Table, error) {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	var tables []*domain.Table

	rows, err := conn.Conn().Query(ctx, "Select * from Tables")
	if err != nil {
		return nil, err
	}

	err = pgxscan.NewScanner(rows).Scan(&tables)
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func (this *PostgresTableRepository) getNextId() int {
	return 1
}

func (this *PostgresTableRepository) NewTable(ctx context.Context, name string, blinds int, numSeats int) (*domain.Table, error) {
	table := new(domain.Table)
	table.Id = this.getNextId() + 1
	this.counter += 1
	table.Name = "Table" + strconv.Itoa(this.counter)

	table.Seats = make([]*domain.Seat, numSeats)
	table.Button = 0
	table.BigBlind = blinds
	for i := 0; i < numSeats; i++ {
		channel := "t_" + strconv.Itoa(table.Id) + "_s_" + strconv.Itoa(i)
		table.Seats[i] = &domain.Seat{0, -1, false, 0, 0, channel, nil}
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
	err = CentClient.Publish(ctx, table.Channel, data)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	//	err = CentClient.Publish(ctx, table.Channel, data)

	return table, nil
}

func (this *PostgresTableRepository) UpdateTable(ctx context.Context, table *domain.Table) error {

	/*
		c, err := this.db.Acquire(context.Background())
		if err != nil {
			fmt.Errorf("Error acquiring DB connection: %v\n", err)
			return err
		}
		defer c.Release()

		_, err = c.Query(context.Background(),
			`update users set last_login=$1, stack =$2 where id=$3`,
			user.LastLogin, user.Stack, user.Id)

		if err != nil {
			return err
		}
	*/
	return nil

}

func (this *BaseTableRepository) UpdatePlayerState(ctx context.Context, table *domain.Table, seatnum int) error {
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()
	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	s := table.Seats[seatnum]


	delta := s.Stack - s.StartingStack

	var ct pgconn.CommandTag
	ct, err = conn.Conn().Exec(newctx,
		"insert into TablePlayerState (table_id, user_id, seatnum, stack_change, starting_stack)  "+
			"Values ($1, $2, $3, $4, $5)", table.Id, s.UserId, s.Seatnum, delta, s.StartingStack)

	ct, err = conn.Conn().Exec(newctx,
		"update TablePlayerState set stack_change = $1, starting_stack = $2 where table_id=$3 and user_id=$4 and seatnum=$5",
		delta, s.StartingStack, table.Id, s.UserId, s.Seatnum)

	_ = ct
	/*
		var rows pgx.Rows
	rows, err = conn.Conn().Query(newctx,
		"insert into TablePlayerState (table_id, user_id, seatnum, stack_change)  "+
			"Values ($1, $2, $3, $4)"+" ON CONFLICT ON CONSTRAINT TablePlayerState_pkey "+
			"DO UPDATE SET seatnum=$5, stack_change = $6;",
		table.Id, s.UserId, s.Seatnum, delta, s.Seatnum, delta)
	 */

	fmt.Printf("Updated TablePlayerState for player %d, seat: %d, stack_change now is %d\n", s.UserId, s.Seatnum, delta)

//	_ = rows
	if err != nil {
		log.Fatalf("Insert/Update GamePlayerState for User %v failed: %v", s.UserId, err)
		return err
	}
	return nil
}

func (this *BaseTableRepository) UpdateAllPlayerStates(ctx context.Context, table *domain.Table) error {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()
	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	for _, s := range table.Seats {

		if s.UserId == -1 {
			continue
		}
		delta := s.Stack - s.StartingStack
		fmt.Printf("seat %d, startingStack = %d, Stack = %d, delta = %d\n", s.Seatnum, s.StartingStack, s.Stack, delta)

		var ct pgconn.CommandTag
		ct, err = conn.Conn().Exec(newctx,
			"insert into TablePlayerState (table_id, user_id, seatnum, stack_change, starting_stack)  "+
				"Values ($1, $2, $3, $4, $5)", table.Id, s.UserId, s.Seatnum, delta, s.StartingStack)

		ct, err = conn.Conn().Exec(newctx,
			"update TablePlayerState set stack_change = $1, starting_stack = $2 where table_id=$3 and user_id=$4 and seatnum=$5",
		      	delta, s.StartingStack, table.Id, s.UserId, s.Seatnum)

		if err != nil {
			fmt.Printf("Error updating TablePlayerState: %v\n", err)
		}
		fmt.Printf("Updated TablePlayerState for player %d, seat: %d, stack change is %d\n", s.UserId, s.Seatnum, delta)

		_ = ct
		if err != nil {
			log.Fatalf("Insert/Update GamePlayerState for User %v failed: %v", s.UserId, err)
			return err
		}
	}
	return nil
}

func (this *BaseTableRepository) GetAllPlayerStates(ctx context.Context, tableId int) ([]*domain.TablePlayerState, error) {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return  nil, err
	}

	defer conn.Release()
	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	rows, err := conn.Conn().Query(newctx,
		`SELECT table_id, user_id, seatnum, stack_change, starting_stack from TablePlayerState where table_id=$1`, tableId)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error loading players states for table:  %d\n", tableId))
	}
	var states []*domain.TablePlayerState


	//	_ := pgxscan.NewScanner(row).Scan(
	if err := pgxscan.NewScanner(rows).Scan(&states); err != nil {
		fmt.Printf("pgxscan err = %v\n", err)
		return nil, err
	}
	rows.Close()
	return states, nil

}

func (this *BaseTableRepository) ClearPlayerStates(ctx context.Context, tableId int) error {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return  err
	}

	defer conn.Release()
	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	_, err = conn.Conn().Query(newctx,
		`DELETE from TablePlayerState where table_id = $1`, tableId)

	if err != nil {
		return  err
	}

	return nil

}
