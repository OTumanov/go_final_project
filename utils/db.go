package utils

import (
	"database/sql"
	"fmt"
	"github.com/OTumanov/go_final_project/model"
	"github.com/OTumanov/go_final_project/settings"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	ENV_DBFILE                            = "TODO_DBFILE"
	INFO_GETTING_PORT_FROM_ENVIRONMENT    = "Получаем порт из окружения..."
	INFO_GETTING_DB_NAME_FROM_ENVIRONMENT = "Получаем имя БД из окружения..."
	INFO_USING_DEFAULT_PORT               = "Порт не задан. Будем использовать 7540"
	PORT_SET                              = "Порт задан - "
	DB_NAME_SET                           = "Имя БД задано -- "
	INFO_DB_NAME_NOT_SET_USING_DEFAULT    = "Имя БД не задано. Будем использовать scheduler.db"
	SQL_DRIVER                            = "sqlite3"
	FAILED_TO_OPEN_DATABASE               = "Не удалось открыть БД: "
	TABLE_CREATION_ERROR                  = "Упс!.. Ошбика при создании таблицы: "
	INDEX_CREATION_ERROR                  = "Упс!.. Ошбика при создании индекса: "
	SQL_CREATE_TABLES                     = "CREATE TABLE IF NOT EXISTS scheduler " +
		"(id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"date TEXT, " +
		"title TEXT, " +
		"comment TEXT, " +
		"repeat VARCHAR(128));"
)

func CheckDb() bool {
	dbName := EnvDBFILE(ENV_DBFILE)

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	_, err = os.Stat(dbFile)

	if err != nil {
		ok, err := installDB(dbName)
		if err != nil && !ok {
			log.Fatal(err)
			return false
		}
	}
	return true
}

func installDB(dbName string) (bool, error) {
	db, err := sql.Open(SQL_DRIVER, dbName)
	if err != nil {
		log.Fatal(FAILED_TO_OPEN_DATABASE, err)
	}
	defer db.Close()

	createTableSQL := SQL_CREATE_TABLES
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(TABLE_CREATION_ERROR, err)
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		log.Fatal(INDEX_CREATION_ERROR, err)
	}
	return true, err
}

func EnvPORT(key string) string {
	log.Println(INFO_GETTING_PORT_FROM_ENVIRONMENT)
	port := os.Getenv(key)
	if len(port) == 0 {
		log.Println(INFO_USING_DEFAULT_PORT)
		port = settings.Port
	} else {
		log.Println(PORT_SET + port)
	}
	return ":" + port
}
func EnvDBFILE(key string) string {
	log.Println(INFO_GETTING_DB_NAME_FROM_ENVIRONMENT)
	dbName := os.Getenv(key)
	if len(dbName) == 0 {
		log.Println(INFO_DB_NAME_NOT_SET_USING_DEFAULT)
		dbName = settings.DBFile
	} else {
		log.Println(DB_NAME_SET + dbName)
	}
	return dbName
}

func addingTask(db *sql.DB, task model.Task) (int64, error) {
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.LastInsertId())
	return res.LastInsertId()
}

func getDB() *sql.DB {
	db, err := sql.Open(SQL_DRIVER, EnvDBFILE(ENV_DBFILE))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
