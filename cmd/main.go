package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"

	"github.com/OTumanov/go_final_project"
	"github.com/OTumanov/go_final_project/pkg/handler"
	"github.com/OTumanov/go_final_project/pkg/repository"
	"github.com/OTumanov/go_final_project/pkg/service"
)

const (
	START_MESSAGE = "Поехали!!! =)"
	ENV_PORT      = "TODO_PORT"
)

func main() {
	logrus.Println(START_MESSAGE)

	if err := initConfig(); err != nil {
		logrus.Fatal(err)
	}

	port := app.EnvPORT(ENV_PORT)
	r := repository.NewRepository(repository.GetDB())
	s := service.NewService(r)
	handlers := handler.NewHandler(s)

	serv := new(app.Server)

	if err := serv.Run(port, handlers.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
