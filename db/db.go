package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"time"
)

func DBUrl() string {
	pg_user := os.Getenv("POSTGRES_USER")
	pg_pw := os.Getenv("POSTGRES_PASSWORD")
	pg_host := os.Getenv("POKER_DB_HOST")
	pg_db := "poker"
	connStr := fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable&password=%s",
		pg_user, pg_host, pg_db, pg_pw)
	return connStr
}

var DB *pgxpool.Pool

type DBConn struct {
	conn     Querier
	poolconn *pgxpool.Conn
	new      bool
	tx       bool
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func (this *DBConn) Conn() Querier {
	if this.tx {
		return this.conn
	} else {
		return this.conn
	}
}

func (this *DBConn) Release() {
	if this.new {
		fmt.Printf("Released conn\n")
		this.poolconn.Release()
	}

}
func GetDBConn(ctx context.Context) (*DBConn, error) {
	var conn Querier
	tx, ok := FromContextWithTx(ctx)
	if tx != nil && ok {
		fmt.Printf("Found tx in context\n")
		return &DBConn{*tx, nil, false, true}, nil
	} else {
		c, ok := FromContextWithConn(ctx)
		if c != nil && ok {
			fmt.Printf("Found conn in context\n")
			return &DBConn{c, nil, false, false}, nil
		} else {
			fmt.Printf("Getting New Connection\n")
			txconn, err := DB.Acquire(ctx)
			if err != nil {
				fmt.Errorf("Unable to get db connection: %v\n", err)
				return nil, err
			}
			conn = txconn.Conn()
			return &DBConn{conn, txconn, true, false}, nil
		}
	}
}

type ContextKey string

const TX_KEY = ContextKey("TX")
const CONN_KEY = ContextKey("Conn")

func NewContextWithConn(ctx context.Context, conn *pgxpool.Conn) context.Context {
	return context.WithValue(ctx, CONN_KEY, conn)
}
func FromContextWithConn(ctx context.Context) (*pgxpool.Conn, bool) {
	c, ok := ctx.Value(CONN_KEY).(*pgxpool.Conn)
	return c, ok
}
func ContextHasConn(ctx context.Context) bool {
	_, ok := ctx.Value(CONN_KEY).(*pgxpool.Conn)
	return ok
}

func NewContextWithTx(ctx context.Context, tx *pgx.Tx) context.Context {
	return context.WithValue(ctx, TX_KEY, tx)
}
func FromContextWithTx(ctx context.Context) (*pgx.Tx, bool) {
	tx, ok := ctx.Value(TX_KEY).(*pgx.Tx)
	return tx, ok
}

func ContextHasTx(ctx context.Context) bool {
	_, ok := ctx.Value(TX_KEY).(*pgx.Tx)
	return ok
}

func SetupDB() *pgxpool.Pool {

	//conn, err := pgx.Connect(context.Background(), DBUrl())

	var dbpool *pgxpool.Pool
	var err error
	success := false
	for i := 0; i < 100; i++ {
		fmt.Printf("DBURL: %v\n", DBUrl())
		dbpool, err = pgxpool.Connect(context.Background(), DBUrl())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database, retrying\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		success = true
		fmt.Fprintf(os.Stderr, "Established Postgres connection\n")
		break
	}
	if !success {
		fmt.Fprintf(os.Stderr, "Unable to connect to database, giving up\n", err)
		os.Exit(-1)
	}
	DB = dbpool
	return DB

}
