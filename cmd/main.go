package main

import (
	"github.com/gin-gonic/gin"
	"os"

	"github.com/OTumanov/go_final_project"
	"github.com/OTumanov/go_final_project/pkg/handler"
	"github.com/OTumanov/go_final_project/pkg/repository"
	"github.com/OTumanov/go_final_project/pkg/service"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	StartMessage            = "Поехали!!! =)"
	EnvPort                 = "TODO_PORT"
	dirDBfile               = "db"
	CheckDBDir              = "Проверка существования каталога: %v..."
	ErrCreateDirectory      = "Ошибка при создании каталога: %v"
	CreatingDirectory       = "Каталог %v не существует. Создаем..."
	SuccessDirectoryCreated = "Каталог успешно создан: %v"
	DirectoryExists         = "Каталог существует"
	InitConfig              = "Инициализация конфигурации..."
	InitConfigDone          = "Конфигурация успешно загружена"
	ErrServerStartReason    = "Ошибка при запуске сервера: %v"
)

func main() {
	logrus.Println(StartMessage)
	checkDBDir()
	gin.SetMode(gin.ReleaseMode)

	if err := initConfig(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Println(InitConfigDone)

	port := app.EnvPORT(EnvPort)
	repo := repository.NewRepository(repository.GetDB())
	srvr := service.NewService(repo)
	handlers := handler.NewHandler(srvr)
	serv := new(app.Server)
	err := serv.Run(port, handlers.InitRoutes())
	if err != nil {
		logrus.Fatalf(ErrServerStartReason, err)
	}
}

func checkDBDir() {
	dirName := dirDBfile
	logrus.Printf(CheckDBDir, dirName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		logrus.Warnf(CreatingDirectory, dirName)
		err := os.Mkdir(dirName, 0700)
		if err != nil {
			logrus.Fatalf(ErrCreateDirectory, err)
			return
		}
		logrus.Printf(SuccessDirectoryCreated, dirName)
	} else {
		logrus.Println(DirectoryExists)
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	logrus.Println(InitConfig)
	return viper.ReadInConfig()
}
