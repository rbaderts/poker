package adapters

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/rbaderts/poker/db"
	"strconv"

	//"github.com/randallmlough/pgxscan"
	"github.com/randallmlough/pgxscan"
	//"github.com/georgysavva/scany/pgxscan"

	"github.com/rbaderts/poker/domain"
	"log"
)

type PostgresTransactionRepository struct {
	counter int
}

func NewPostgresTransactionRepository() domain.TransactionRepository {
	r := new(PostgresTransactionRepository)
	return r
}
func (this *PostgresTransactionRepository) DeleteAllTransactionsForTable(ctx context.Context, tableId int) error {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}

	tableIdStr := strconv.Itoa(tableId)
	_, err = conn.Conn().Query(ctx,
		"Delete from Transactions where table_id=$1", tableIdStr)

	return err
}

func (this *PostgresTransactionRepository) GetAllTransactionsForTable(ctx context.Context, tableId int) ([]*domain.Transaction, error) {
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	var transactions []*domain.Transaction

	tableIdStr := strconv.Itoa(tableId)
	rows, err := conn.Conn().Query(ctx,
		"Select * from Transactions where table_id=$1", tableIdStr)
	if err != nil {
		return nil, err
	}

	//	err = pgxscan.Select(ctx, conn.Conn(), &transactions, "SELECT * from Transactions where table_id="+tableIdStr)
	err = pgxscan.NewScanner(rows).Scan(&transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (this *PostgresTransactionRepository) RecordTransaction(
	ctx context.Context,
	typ domain.TransactionType,
	amount int,
	userId int,
	tableId int) error {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var rows pgx.Rows
	rows, _ = conn.Conn().Query(ctx,
		"insert into Transactions(typ, user_id, table_id, amount) values($1,$2,$3,$4) "+
			" returning id, user_id, table_id, typ, amount, transaction_date;",
		typ, userId, tableId, amount)

	var transactionRecord domain.Transaction
	//	if err := pgxscan.ScanOne(&transactionRecord, rows); err != nil {
	if err := pgxscan.NewScanner(rows).Scan(&transactionRecord); err != nil {

		log.Fatalf("Insert Transaction failed: %v", err)
		return err
	}
	rows.Close()

	return nil

}
