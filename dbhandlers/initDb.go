package dbhandlers

import(
	"github.com/jackc/pgx"
)

type DataBase struct {
	pool *pgx.ConnPool
}

var DB DataBase

func (db *DataBase) Connetc() error {
	conConfig := pgx.ConnConfig {
		Host: 			"127.0.0.1",
		Port: 			5432,
		Database: 		"new_database",
		User: 			"docker",
		Password: 		"docker",
		TLSConfig: 		nil,
		UseFallbackTLS: false,
		RuntimeParams: 	map[string]string{"application_name": "new_schema"},
	}
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     conConfig,
		MaxConnections: 25,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	con, err := pgx.NewConnPool(poolConfig)
	db.pool = con
	
	return err
}

func ErrorCode(err error) (string) {
	pgerr, ok := err.(pgx.PgError)
	if !ok {
		return pgxOK
	}
	return pgerr.Code
}