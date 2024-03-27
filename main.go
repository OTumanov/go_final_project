package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	webDir = "./web"
	dbName = "scheduler.db"
)

func needInstallDb() bool {
	log.Println("Ищем файл БД sqlite...")
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	_, err = os.Stat(dbFile)

	log.Println("Вот чего нам вернулось вместо БД -- " + err.Error())
	log.Println("А это значит, что БД не найдена. Будем создавать...")

	if err != nil {
		ok, err := installDB()
		if err != nil && !ok {
			log.Fatal(err)
			return false
		}
	}
	return true
}

func installDB() (bool, error) {
	log.Println("Создаем БД...")
	db, err := sql.Open("sqlite3", dbName)
	log.Println("Создали пустую БД")
	log.Println("Установили соединение с ней")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Ошибок пока, вроде, нет =)")
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATE,
		title TEXT,
		comment TEXT,
		repeat VARCHAR(128)
	);
	`
	log.Println("Создали таблицы в БД")

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Упс!.. Ошибка при создании таблиц: ", err)
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		log.Fatal("Упс!.. Ошбика при создании индекса: ", err)
	}

	log.Println("БД создана, настроена и подключена")
	return true, err
}

func listenPort(key string) string {
	port := ":" + os.Getenv(key)
	return port
}

func server() {
	log.Println("Запускаем HTTP-сервер...")
	log.Println("Будем использовать порт" + listenPort("TODO_PORT"))
	log.Println("Вот тут -- http://localhost" + listenPort("TODO_PORT"))
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(listenPort("TODO_PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Println("Запуск! =)")
	if needInstallDb() {
		server()
	}
}
