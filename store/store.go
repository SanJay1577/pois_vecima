package store

import (
	"pois/config"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// Store is a DAO handle to the backend DB/memstore. This
// application has a relatively simple schema, so this is
// currently modeled as embedded structs, but we can split
// out the individual struct keys to substores if things
// start getting too cluttered.
//
// We have a Store interface here so that we can easily
// sub out the SQL backend store in model unit tests.
type Store interface {
	Alias
	Schedule
}

// DataStore is a postgres backed Store.
type DataStore struct {
	SQLAlias
	SQLSchedule
}

func MustOpenDB(log *zaplogger.Logger, env string) *sqlx.DB {
	log.Infof("opening db connection: %v / %v", config.GetConfig().GetString("app.db.dialect"), config.GetConfig().GetString("app.db.datasource"))
	db := sqlx.MustOpen(config.GetConfig().GetString("app.db.dialect"), config.GetConfig().GetString("app.db.datasource"))
	err := db.Ping()
	if err != nil {
		log.Fatalf("could not connect to database - error: %v", err)
	}

	return db
}

// New opens a new db connection and returns a concrete Store.
func New(log *zaplogger.Logger, env string) *DataStore {
	db := MustOpenDB(log, env)
	return &DataStore{
		SQLAlias{DB: db, Log: log},
		SQLSchedule{DB: db, Log: log},
	}
}
