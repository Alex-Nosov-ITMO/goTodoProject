package db

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)


type DbConnection struct {
	Db *sqlx.DB
}

type ConfigDB struct {
	User     string
	Password string
	Host     string
	Port     string
	Dbname   string
}

func NewDbConnection() (*DbConnection, error) {
	
	// TODO: connect to db
	appPath, err := os.Getwd()
	if err != nil {
		log.Printf("DbConnection: can't get app path: %s\n", err)
		return nil, errors.New("can't get app path")
	}
	dbFile := filepath.Join(appPath, os.Getenv("DB_NAME"))

	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sqlx.Open(os.Getenv("DB_DRIVER"), dbFile)
	if err != nil {
		log.Printf("DbConnection: can't open db: %s\n", err)
		return nil, errors.New("can't open db")
	}

	if install {
		_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date СHAR(8) NOT NULL,
			title VARCHAR NOT NULL,
			comment TEXT,
			repeat VARCHAR(128)
		);
		CREATE INDEX idx_scheduler_date ON scheduler (date);
		`)
		if err != nil {
			log.Printf("DbConnection: can't create db: %s\n", err)
			return nil, errors.New("can't create db")
		}
	}

	return &DbConnection{
		Db: db,
	}, nil
}

// TODO: close db
func (db *DbConnection) Close() {
	db.Db.Close()
}