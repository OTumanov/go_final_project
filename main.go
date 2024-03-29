package main

import (
	"github.com/OTumanov/go_final_project/utils"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const START_MESSAGE = "Поехали!!! =)"

func main() {
	log.Println(START_MESSAGE)
	if utils.CheckDb() {
		utils.StartHTTPServer()
	}
}
