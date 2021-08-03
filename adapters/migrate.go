package adapters

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/rbaderts/poker/db"
	"log"
)

func MigrateDB() {

	url := db.DBUrl()
	fmt.Printf("MigrateDB\n")
	m, err := migrate.New(
		"file://migrations", url)
	if err != nil {
		log.Fatal(err)
	}
	m.Up()
	fmt.Printf("migrated\n")
}

/*
func SeedAccount(db *dbr.Session) (*Account, error) {

	var count int
	err := db.Select("count(*)").From("accounts").LoadOne(&count)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return CreateAccount(db, "Default")
	}

	return nil, nil
}
*/
