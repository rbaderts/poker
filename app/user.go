package app

import (
	"context"
	"github.com/rbaderts/poker/adapters"
	"github.com/rbaderts/poker/domain"
)


type UserService interface {

	NewUser(userId int, name string) (*adapters.User,error)
	LoadUserById(userId int) (*adapters.User,error)
	AddMoney(userId int, amount int) (error)
	DeductMoney(userId int, amount int) (error)

}
var GblUserService UserService

type UserServiceImpl struct {
	userRepository domain.UserRepository;
}

func (this *UserServiceImpl)  AddMoney(userId int, amount int) (error)  {
	ctx := context.Background()
	err := this.userRepository.AddMoney(ctx, userId, amount)
	if err != nil {
		return err
	}
	return nil
}

func (this *UserServiceImpl) DeductMoney(userId int, amount int) (error)  {
	ctx := context.Background()
	err := this.userRepository.DeductMoney(ctx, userId, amount)
	if err != nil {
		return err
	}

	return nil
}

func (this *UserServiceImpl)  LoadUserById(userId int) (*adapters.User, error)  {
	ctx := context.Background()
	userRec, err := this.userRepository.LoadUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	user := &adapters.User{userRec.Id, userRec.Subject, userRec.Email, userRec.PictureUrl,
		userRec.GivenName, userRec.Stack}
	return user, nil

}

func (this *UserServiceImpl) NewUser(id int, name string) (*adapters.User,error) {
	ctx := context.Background()
	userRec, err := this.userRepository.AddUser(ctx, id, name)
	if err != nil {
		return nil, err
	}
	user := &adapters.User{userRec.Id, userRec.Subject, userRec.Email, userRec.PictureUrl,
		userRec.GivenName, userRec.Stack}
	return user, nil
}

func NewUserService(repo domain.UserRepository) *UserServiceImpl {

	this := new(UserServiceImpl)
	this.userRepository = repo
	return this

}

