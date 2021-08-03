package domain

import (
	"context"
	"time"
)




type TransactionRepository interface {
	GetAllTransactionsForTable(ctx context.Context, tableId int) ([]*Transaction, error)
	DeleteAllTransactionsForTable(ctx context.Context, tableId int) error
	RecordTransaction(ctx context.Context, typ TransactionType, amount int, userId int, tableId int) error

}

type TransactionType int

const (
	UserAccountToTable       = 1
	TableToUserAccount       = 2
	ExternalToUserAccount    = 3
	UserAccountToExternal    = 4
)


type Transaction struct {
	Id                int             `json:"id" db:"id"`
	UserId            int             `json:"userId" db:"user_id"`
	TableId           int             `json:"tableId" db:"table_id"`
	Typ               TransactionType `json:"type" db:"typ"`
	Amount            int             `json:"amount" db:"amount"`
	TransactionDate   time.Time       `json:"date" db:"transaction_date"`
}


