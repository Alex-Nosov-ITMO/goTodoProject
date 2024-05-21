package db

import (
	"os"
	"path/filepath"

	//_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
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
	//connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Password, conf.Dbname)

	/*db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return &DbConnection{
		Db: db,
	}	*/
	
	// TODO: connect to db
	appPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(appPath, os.Getenv("DB_NAME"))

	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sqlx.Open(os.Getenv("DB_DRIVER"), dbFile)
	if err != nil {
		return nil, err
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
			return nil, err
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