package app

import (
	"context"
	"fmt"
	"github.com/rbaderts/poker/domain"
)

type AccountService interface {
    NewAccount(name string) *domain.Account
}

type AccountServiceState struct {
}

var GblAccountService AccountService

func NewAccountService() AccountService {
	s := new(AccountServiceState)
	return s
}


func (this *AccountServiceState) NewAccount(name string) *domain.Account {

	ctx := context.Background()
	account, err := domain.GblAccountRepository.CreateAccount(ctx, name)
	_ = err
	fmt.Printf("NewAccount: id = %s\n", account.Id)
	return account;
}

