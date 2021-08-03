package domain

/*

      A poker application consists of a lobby, 1 or more Tables, a cashier.

      A poker application allows users to login.   The application maintains an account balance for each User.
      A User can deposit or withdrawl money.     A User can sit down at a table with seat avaialbility at which point
      the User becomes a "player" at that table.   When sitting down at a table a User must move some amount of money
      from her account to the table (becoming her Stack).    When a player exists a table their current Stack is
       added to their account balance.

      A User log's into poker via a client application.


 */




type ApplicationState struct {
	//sessions []*UserSession
	//UserRepository *adapters.UserRepository
	//UserService *app.UserServiceImpl
	//DB *pgxpool.Pool
	tables map[string]*Table
	dispatcher MessageDispatcher
}

func NewApplication(dispatcher MessageDispatcher) {
    a := new(ApplicationState)
	a.tables = make(map[string]*Table, 0)
	a.dispatcher = dispatcher

}
