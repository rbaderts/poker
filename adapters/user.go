package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rbaderts/poker/db"

	//	"github.com/georgysavva/scany/pgxscan"
	"context"
	"github.com/randallmlough/pgxscan"
	"github.com/rbaderts/poker/domain"
	"log"
	"time"
)

type UserRecord struct {
	Id              int            `json:"id"`
	AccountId       int            `json:"accountId" db:"account_id"`
	Subject         string         `json:"subject"`
	Email           string         `json:"email"`
	Provider        string         `json:"provider"`
	PictureUrl      sql.NullString `json:"pictureUrl" db:"picture_url"`
	GivenName       string         `json:"givenName" db:"given_name"`
	PreferredHandle string         `json:"preferredHandle" db:"preferred_handle"`
	LastLogin       time.Time      `json:"lastLogin" db:"last_login"`
	Stack           int            `json:"stack"  db:"stack"`
}

type User struct {
	Id         int            `json:"id"`
	Subject    string         `json:"subject"`
	Email      string         `json:"email"`
	PictureUrl sql.NullString `json:"pictureUrl" db:"picture_url"`
	GivenName  string         `json:"givenName" db:"given_name"`
	Stack      int            `json:"stack"  db:"stack"`
}

type PostgresUserRepository struct {
}

func NewPostgresUserRepository() domain.UserRepository {
	r := new(PostgresUserRepository)
	return r
}

/*
func UserRecordToUser(ur *UserRecord) *domain.User {
	u :+ new(domain.User)

}
*/

func (this *PostgresUserRepository) LoadUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	fmt.Printf("LoadUserByEmail: %s\n", email)

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	//	c := this.poker.GetDBConn(ctx, this.db, tx)

	rows, err := conn.Conn().Query(newctx,
		`SELECT id, account_id, subject, email, provider, picture_url, given_name, preferred_handle, last_login, stack 
             from users where email = $1`, email)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error loading user with enail %s\n", email))

	}
	var user UserRecord

	//	_ := pgxscan.NewScanner(row).Scan(
	if err := pgxscan.NewScanner(rows).Scan(&user); err != nil {
		//if err := pgxscan.ScanOne(&user, rows); err != nil {
		fmt.Printf("pgxscan err = %v\n", err)
		return nil, err
	}
	rows.Close()
	fmt.Printf("user = %v\n", user)

	u := UserRecordToUser(&user)
	return u, nil
}

func (this *PostgresUserRepository) LoadUserById(ctx context.Context, id int) (*domain.User, error) {

	fmt.Printf("LoadUserById: %d\n", id)
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	rows, err := conn.Conn().Query(newctx,
		`SELECT * from users where id = $1`, id)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error loading user with id %d\n", id))
	}
	var user UserRecord
	if err := pgxscan.NewScanner(rows).Scan(&user); err != nil {
		//if err := pgxscan.ScanOne(&user, rows); err != nil {
		fmt.Errorf("LoadUserById error: %v\n", err)
		return nil, err
	}
	rows.Close()

	u := UserRecordToUser(&user)
	fmt.Printf("LoadUserById - returning user %d\n", user.Id)
	return u, nil
}

func (this *PostgresUserRepository) LoadUserBySubject(ctx context.Context, subject string) (*domain.User, error) {
	fmt.Printf("LoadUserBySubject: %s\n", subject)
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Conn().Query(ctx,
		`SELECT * from users where subject = $1`, subject)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error loadingu ser with subject %s\n", subject))
	}
	var user UserRecord

	//if err := pgxscan.ScanOne(&user, rows); err != nil {
	if err := pgxscan.NewScanner(rows).Scan(&user); err != nil {
		return nil, err
	}

	rows.Close()

	u := UserRecordToUser(&user)
	return u, nil
}

func (this *PostgresUserRepository) AddMoney(ctx context.Context, userId int, amount int) error {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	fmt.Printf("Adding %d dollars to user: %d\n", amount, userId)
	user, err := this.LoadUserById(newctx, userId)
	if err != nil {
		return err
	}
	user.Stack += amount
	return this.UpdateUser(newctx, user)
}

func (this *PostgresUserRepository) DeductMoney(ctx context.Context, userId int, amount int) error {
	fmt.Printf("Deducting %d dollars from user: %d\n", amount, userId)

	fmt.Printf("DeductMoney   CTX has a TX = %v\n", db.ContextHasTx(ctx))
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

//	newctx := context.WithValue(ctx, db.CONN_KEY, conn.Conn())

	fmt.Printf("DeductMoney   newctx has a TX = %v\n", db.ContextHasTx(ctx))

	user, err := this.LoadUserById(ctx, userId)
	if err != nil {
		return err
	}
	user.Stack -= amount
	return this.UpdateUser(ctx, user)
}

func (this *PostgresUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {

	fmt.Printf("UpdateUser\n")

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Conn().Query(ctx,
		`update users set last_login=$1, stack =$2 where id=$3`,
		user.LastLogin, user.Stack, user.Id)

	if err != nil {
		return err
	}
	return nil
}

/*
func (this *PostgresUserRepository) AddUserFromProfile(ctx context.Context, accountId int, profile map[string]interface{}) (*domain.User, error) {

	conn, err := poker.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}

	var rows pgx.Rows
	rows, _ = conn.Query(ctx,
		"insert into Users(account_id, subject, email, provider, given_name, preferred_handle) values($1,$2,$3,$4,$5,$6) " +
			" ON CONFLICT (subject) DO NOTHING " +
			" returning id, account_id, subject, email, provider, given_name, preferred_handle;",
		accountId,
		getProfileString("sub", profile),
		getProfileString("email", profile),
		getProfileString("iss", profile),
		getProfileString("given_name", profile),
		getProfileString("given_name", profile))

	var user UserRecord
	wScanner(row).Scan(&user); err != nil {
		return nil, err
	}

	rows.Close()

	u := UserRecordToUser(&user)
	return u, nil
}

func (this *PostgresUserRepository) AddMoney (ctx context.Context, userId int, amount int) error {

	conn, err := poker.GetDBConn(ctx)
	if err != nil {
		return err
	}
	newctx := context.WithValue(ctx, "DBConn", conn)

	fmt.Printf("Adding %d dollars to user: %d\n", amount, userId)
	user, err := this.LoadUserById(newctx, userId)
	if err != nil {
		return  err
	}
	user.Stack += amount
	return this.UpdateUser(newctx, user)
}

*/

/*
func (this *PostgresUserRepository) DeductMoney (ctx context.Context, userId int, amount int) error {
	fmt.Printf("Deducting %d dollars from user: %d\n", amount, userId)

	conn, err := poker.GetDBConn(ctx)
	if err != nil {
		return err
	}
	newctx := context.WithValue(ctx, "DBConn", conn)

	user, err := this.LoadUserById(newctx,  userId)
	if err != nil {
		return  err
	}
	user.Stack -= amount
	return this.UpdateUser(newctx, user)
}

*/

/*
func (this *PostgresUserRepository) UpdateUser(ctx context.Context, user *domain.User) (error) {

	conn, err := poker.GetDBConn(ctx)
	if err != nil {
		return err
	}
	_, err = conn.Query(ctx,
		`update users set last_login=$1, stack =$2 where id=$3`,
		user.LastLogin, user.Stack, user.Id)

	if err != nil {
		return err
	}
	return nil
}
*/

func (this *PostgresUserRepository) AddUserFromProfile(ctx context.Context, accountId int, profile map[string]interface{}) (*domain.User, error) {

	fmt.Printf("AddUserFromProfile\n")
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Conn().Query(ctx,
		"insert into Users(account_id, subject, email, provider, given_name, preferred_handle) values($1,$2,$3,$4,$5,$6) "+
			" ON CONFLICT (subject) DO NOTHING "+
			" returning id, account_id, subject, email, provider, given_name, preferred_handle;",
		accountId,
		getProfileString("sub", profile),
		getProfileString("email", profile),
		getProfileString("iss", profile),
		getProfileString("given_name", profile),
		getProfileString("given_name", profile))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error adding user\n"))

	}
	var user UserRecord

	//	if err := pgxscan.ScanOne(&user, rows); err != nil {
	if err := pgxscan.NewScanner(rows).Scan(&user); err != nil {
		log.Fatalf("Insert User failed: %v", err)
		return nil, err
	}
	rows.Close()

	u := UserRecordToUser(&user)
	return u, nil

}

func UserRecordToUser(rec *UserRecord) *domain.User {
	user := new(domain.User)
	user.Stack = rec.Stack
	user.Id = rec.Id
	user.GivenName = rec.GivenName
	user.PictureUrl = rec.PictureUrl
	user.Email = rec.Email
	user.Subject = rec.Subject
	user.PreferredHandle = rec.PreferredHandle

	return user

}

func (this *PostgresUserRepository) AddProvidedUser(
	ctx context.Context, accountId int, email string,
	provider string, subject string, givenName string) (*domain.User, error) {
	fmt.Printf("AddProvidedUser\n")

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	rows, err := conn.Conn().Query(ctx,
		"insert into Users(account_id, subject, email, provider, given_name, preferred_handle) values($1,$2,$3,$4,$5,$6) "+
			" ON CONFLICT (subject) DO NOTHING "+
			" returning id, account_id, subject, email, provider, given_name, preferred_handle;",
		accountId, subject, email, provider, givenName, givenName)

	if err != nil {
		log.Fatalf("Insert User failed: %v", err)
		return nil, err
	}
	var user UserRecord
	//if err := pgxscan.ScanOne(&user, rows); err != nil {
	if err := pgxscan.NewScanner(rows).Scan(&user); err != nil {
		log.Fatalf("Insert User failed: %v", err)
		return nil, err
	}
	rows.Close()
	if err != nil {
		log.Fatalf("Insert User failed on scan: %v", err)
		return nil, err
	}

	u := UserRecordToUser(&user)
	return u, nil

}

func getProfileString(key string, profile map[string]interface{}) string {

	v := profile[key]
	if v == nil {
		return ""
	}
	return v.(string)

}

func (this *PostgresUserRepository) AddUser(ctx context.Context, id int, name string) (*domain.User, error) {
	return nil, nil
}
