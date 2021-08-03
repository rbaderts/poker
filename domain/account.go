package domain

import (
	"context"
)


type Account struct {
	Id            int    `json:"id"`
	AdminUUser    int    `json:"AdminUser" db:"admin_user_id"`
	AccountName   string `json:"accountName" db:"account_name"`
}

type AccountRepository interface {
    LoadAccountByName(context.Context, string) (*Account, error)
	CreateAccount(context.Context, string) (*Account, error)

}



