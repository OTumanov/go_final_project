package utils

import (
	"log"
	"os"
)

func EnvPORT(key string) string {
	log.Println("Получаем порт из окружения...")
	port := os.Getenv(key)
	if len(port) == 0 {
		log.Println("Порт не задан. Будем использовать 7540")
		port = "7540"
	} else {
		log.Println("Порт задан -- " + port)
	}
	return ":" + port
}

func EnvDBFILE(key string) string {
	log.Println("Получаем имя БД из окружения...")
	dbName := os.Getenv(key)
	if len(dbName) == 0 {
		log.Println("Имя БД не задано. Будем использовать scheduler.db")
		dbName = "scheduler.db"
	} else {
		log.Println("Имя БД задано -- " + dbName)
	}
	return dbName
}
