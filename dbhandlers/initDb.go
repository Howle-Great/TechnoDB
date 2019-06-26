package dbhandlers

import(
	"github.com/jackc/pgx"
	"io/ioutil"
)

type DataBase struct {
	Pool *pgx.ConnPool
}

var DB DataBase
const dbSchema = "./db/init.sql"

func (db *DataBase) Connetc() error {
	conConfig := pgx.ConnConfig {
		Host: 			"127.0.0.1",
		Port: 			5432,
		Database: 		"zxc",
		User: 			"docker",
		Password: 		"docker",
		TLSConfig: 		nil,
		UseFallbackTLS: false,
		RuntimeParams: 	map[string]string{"application_name": "dz"},
	}
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     conConfig,
		MaxConnections: 25,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	con, err := pgx.NewConnPool(poolConfig)
	db.Pool = con
	if err != nil {
		return err
	}
	err = LoadSchemaSQL()
	
	return err
}

func LoadSchemaSQL() error {
	if DB.Pool == nil {
		return pgx.ErrDeadConn
	}

	content, err := ioutil.ReadFile(dbSchema)
	if err != nil {
		return err
	}

	tx, err := DB.Pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(string(content)); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func ErrorCode(err error) (string) {
	pgerr, ok := err.(pgx.PgError)
	if !ok {
		return pgxOK
	}
	return pgerr.Code
}