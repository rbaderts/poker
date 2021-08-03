package adapters

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rbaderts/poker/db"

	//	"github.com/georgysavva/scany/pgxscan"
	"github.com/randallmlough/pgxscan"
	"github.com/rbaderts/poker/domain"
)

type PostgresRecoveryRepository struct {
	db *pgxpool.Pool
	counter int
}

func NewPostgresRecoveryRepository() domain.RecoveryRepository {
	r := new(PostgresRecoveryRepository)
	return r
}

func (this *PostgresRecoveryRepository) CreateTableTombstone(ctx context.Context, tableId int) error {
	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Conn().Query(context.Background(),
		"insert into TableRecovery (table_id)  "+
			"Values ($1);", tableId)

	if err != nil {
		return err
	}
	return nil

}

type TableId struct {
	TableId int `db:"table_id"`
}

func (this *PostgresRecoveryRepository) DeleteTableTombstone(ctx context.Context, tableId int) error {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Conn().Query(ctx, "Delete from TableRecovery where table_id = $1", tableId)
	return err

}
func (this *PostgresRecoveryRepository) ReadTableTombstones(ctx context.Context) ([]int, error) {

	conn, err := db.GetDBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	var ids []*TableId

	rows, err := conn.Conn().Query(ctx, "SELECT table_id from TableRecovery")
	//pgxscan.Select(ctx, conn.Conn(), &ids, "SELECT table_id from TableRecovery");
	if err != nil {
		return nil, err
	}

	if err := pgxscan.NewScanner(rows).Scan(&ids); err != nil {
		return nil, err
	}

	result := make([]int, len(ids))
	for i, id := range ids {
		result[i] = id.TableId
	}
	return result, nil
}
