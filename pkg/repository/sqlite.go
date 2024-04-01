package repository

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/OTumanov/go_final_project/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	ENV_DBFILE                            = "TODO_DBFILE"
	INFO_GETTING_DB_NAME_FROM_ENVIRONMENT = "Получаем имя БД из окружения..."
	DB_NAME_SET                           = "Имя БД задано -- "
	INFO_DB_NAME_NOT_SET_USING_DEFAULT    = "Имя БД не задано. Будем использовать из конфига "
	SQL_DRIVER                            = "sqlite3"
	FAILED_TO_OPEN_DATABASE               = "Не удалось открыть БД: "
	TABLE_CREATION_ERROR                  = "Упс!.. Ошбика при создании таблицы: "
	INDEX_CREATION_ERROR                  = "Упс!.. Ошбика при создании индекса: "

	taskTable   = "scheduler"
	taskDate    = "date"
	taskTitle   = "title"
	taskComment = "comment"
	taskRepeat  = "repeat"
)

type Config struct {
	SQLDriver string
	DBFile    string
}

func NewSqlite(config *Config) (*sqlx.DB, error) {
	dbName := EnvDBFILE(ENV_DBFILE)

	db, err := sqlx.Connect(config.SQLDriver, dbName)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		logrus.Fatal(err)
	}

	return db, nil
}

func GetDB() *sqlx.DB {
	dbname, err := CheckDb()
	if err != nil {
		logrus.Fatal(err)
	}
	return sqlx.MustConnect(viper.Get("DB.SQLDriver").(string), dbname)
}

func CheckDb() (string, error) {
	dbName := EnvDBFILE(ENV_DBFILE)

	appPath, err := os.Executable()
	if err != nil {
		logrus.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	_, err = os.Stat(dbFile)

	if err != nil {
		ok, err := installDB(dbName)
		if err != nil && !ok {
			logrus.Fatal(err)
			return dbName, err
		}
	}
	return dbName, nil
}

func EnvDBFILE(key string) string {
	logrus.Println(INFO_GETTING_DB_NAME_FROM_ENVIRONMENT)
	dbName := os.Getenv(key)
	if len(dbName) == 0 {
		dbName = viper.Get("DB.DBFile").(string)
		logrus.Println(INFO_DB_NAME_NOT_SET_USING_DEFAULT + dbName)
	} else {
		logrus.Println(DB_NAME_SET + dbName)
	}
	return dbName
}

func installDB(dbName string) (bool, error) {
	db, err := sql.Open(SQL_DRIVER, dbName)
	if err != nil {
		logrus.Fatal(FAILED_TO_OPEN_DATABASE, err)
	}
	defer db.Close()

	createTableSQL := config.SQLCreateTables
	_, err = db.Exec(createTableSQL)
	if err != nil {
		logrus.Fatal(TABLE_CREATION_ERROR, err)
	}

	createIndexSQL := config.SQLCreateIndexes
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		logrus.Fatal(INDEX_CREATION_ERROR, err)
	}
	return true, err
}
