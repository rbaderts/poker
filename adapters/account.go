package adapters

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	//"github.com/randallmlough/pgxscan"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/rbaderts/poker/db"
	"github.com/rbaderts/poker/domain"
	"log"
)

type PostgresAccountRepository struct {
}

func NewPostgresAccountRepository() domain.AccountRepository {
	r := new(PostgresAccountRepository)
	return r
}

func (this *PostgresAccountRepository) LoadAccountByName(ctx context.Context, name string) (*domain.Account, error) {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Conn().Query(ctx,
		`SELECT * from accounts where account_name = $1`, name)
	if err != nil {
		fmt.Printf("Account query error: %v\n", err)
		return nil, err

	}
	var act domain.Account
	//	if err := pgxscan.NewScanner(rows).Scan(&act); err != nil {
	if err := pgxscan.ScanOne(&act, rows); err != nil {
		fmt.Printf("pgxscan err = %v\n", err)
		return nil, err
	}
	rows.Close()
	fmt.Printf("act = %v\n", act)

	return &act, nil
}

func (this *PostgresAccountRepository) CreateAccount(ctx context.Context, name string) (*domain.Account, error) {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var rows pgx.Rows
	rows, _ = conn.Conn().Query(ctx,
		"insert into Accounts(admin_user_id, account_name) values($1,$2) "+
			" returning id, admin_user_id, account_name",
		0, name)

	var account domain.Account

	//if err := pgxscan.NewScanner(rows).Scan(&account); err != nil {
	if err := pgxscan.ScanOne(&account, rows); err != nil {
		log.Fatalf("Insert Account failed: %v", err)
		return nil, err
	}

	return &account, nil
}
