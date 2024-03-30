package main

import (
	_ "github.com/OTumanov/go_final_project/model"
	"github.com/OTumanov/go_final_project/utils"

	"log"
)

const START_MESSAGE = "Поехали!!! =)"

func main() {
	log.Println(START_MESSAGE)
	if utils.CheckDb() {
		utils.StartHTTPServer()
	}
}
