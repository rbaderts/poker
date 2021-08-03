package app

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rbaderts/poker/adapters"
	"github.com/rbaderts/poker/domain"
)

var (
	Store *sessions.FilesystemStore
)

func init() {
	gob.Register(map[string]interface{}{})
//	Store = sessions.NewFilesystemStore("", []byte("something-very-secret"))
}

var MainApp *Application

type Application struct {
	//	sessions []*UserSession
	//UserRepository domain.UserRepository
	//TableRepository domain.TableRepository
	//UserService *app.UserServiceImpl
	//MessageDispatcher domain.MessageDispatcher
	DB    *pgxpool.Pool
	table *domain.Table
	//tables map[string]*domain.Table
}

func NewReplayApplication(db *pgxpool.Pool) *Application {

	application := new(Application)
	application.DB = db

	domain.GblUserRepository = adapters.NewPostgresUserRepository()
	domain.GblTableRepository = adapters.NewPostgresTableRepository()
	domain.GblAccountRepository = adapters.NewPostgresAccountRepository()
	GblTableService = NewTableService()
	GblUserService = NewUserService(domain.GblUserRepository)
	domain.GblMessageDispatcher = adapters.NewNullDispatcher()

	MainApp = application
	return application

}

func NewApplication(db *pgxpool.Pool) *Application {
	fmt.Printf("NewApplication")
	application := new(Application)
	application.DB = db

	domain.GblUserRepository = adapters.NewPostgresUserRepository()
	domain.GblTableRepository = adapters.NewMemoryTableRepository()
	domain.GblAccountRepository = adapters.NewPostgresAccountRepository()
	//	domain.GblGameStateRepository = adapters.NewPostgresGameStateRepository()
	domain.GblTransactionRepository = adapters.NewPostgresTransactionRepository()
	domain.GblRecoveryRepository = adapters.NewPostgresRecoveryRepository()

	domain.GblDBPool = db

	GblTableService = NewTableService()
	GblUserService = NewUserService(domain.GblUserRepository)

	//application.UserService = app.NewUserService(application.UserRepository)
	domain.GblMessageDispatcher = adapters.NewNatsDispatcher()

	//app.tables = make(map[string]*Table, 0)
	MainApp = application

	return application
}

func (this *Application) SetTable(table *domain.Table) {
	this.table = table
}
func (this *Application) NewTableWithFixedDec(deckfile string, seats int) *domain.Table {
	//	if err != nil || table == nil {
	table := domain.NewTable("test", 2, 6)
	//	this.tables[table.Id] = table
	//deck := pokerlib.ReadDeck(deckfile)
	//	go table.TablePlaybackGame(deck)
	return table
}
