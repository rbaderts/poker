package adapters

import (
	"context"
	"fmt"
	"github.com/rbaderts/poker/domain"
)

type MemoryUserRepository struct {
	users map[int]*domain.User
	counter int
}

func NewMemoryUserRepository() domain.UserRepository {
	r := new(MemoryUserRepository)
	r.users = make(map[int]*domain.User)
	return r
}




func (this * MemoryUserRepository) GetAllUsers(ctx context.Context) []*domain.User {

	r := make([]*domain.User, len(this.users))

	i := 0
	for _, v := range this.users {
		r[i] = v
		i += 1
	}

	return r
}

func (this * MemoryUserRepository) GetUser(ctx context.Context, id int) *domain.User {
	k, ok := this.users[id]

	if !ok {
		return nil
	}
	return k
}


func (this *MemoryUserRepository) AddUser(ctx context.Context, id int, name string) (*domain.User, error) {

	fmt.Printf("MemoryUserRepository: NewUser\n")
	user := new(domain.User)
	user.Id = id
	user.GivenName = name
	this.users[id] = user
	return user, nil
}

func (this *MemoryUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {

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

func (this *MemoryUserRepository) AddUserFromProfile(ctx context.Context, accountId int, profile map[string]interface{}) (*domain.User, error) {
	return nil, nil
}
func (this* MemoryUserRepository) AddProvidedUser(ctx context.Context, accountId int, email string, provider string, subject string, givenName string) (*domain.User, error) {
	return nil, nil
}
func (this* MemoryUserRepository) LoadUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (this* MemoryUserRepository) LoadUserById(ctx context.Context, uid int) (*domain.User, error) {
	return this.users[uid], nil
}
func (this* MemoryUserRepository) LoadUserBySubject(ctx context.Context, subject string) (*domain.User, error) {
	return nil, nil
}

func (this* MemoryUserRepository) AddMoney(ctx context.Context, id int, amount int) error {
	return nil
}
func (this* MemoryUserRepository) DeductMoney(ctx context.Context, id int, amount int) error {
	return nil

}
