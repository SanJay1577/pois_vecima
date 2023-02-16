package store_test

import (
	"os"
	"testing"

	"pois/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var DB *sqlx.DB

func MustMigrateDB(db *sqlx.DB) *sqlx.DB {
	db.MustExec(`
		DROP SCHEMA test CASCADE;
		CREATE SCHEMA test;
		GRANT ALL ON SCHEMA test TO postgres;
		GRANT ALL ON SCHEMA test TO public;
	`)

	migrate.SetTable(config.GetConfig().GetString("app.db.table"))
	migrations := &migrate.FileMigrationSource{Dir: "./pois/config"}

	_, err := migrate.Exec(db.DB, config.GetConfig().GetString("app.db.dialect"), migrations, migrate.Up)
	if err != nil {
		panic(err)
	}

	return db
}

func TestMain(m *testing.M) {
	DB = sqlx.MustOpen(config.GetConfig().GetString("app.db.dialect"), config.GetConfig().GetString("app.db.datasource"))
	err := DB.Ping()
	if err != nil {
		panic(err)
	}

	defer DB.Close()

	MustMigrateDB(DB)

	os.Exit(m.Run())
}
