package app

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	INFO_GETTING_PORT_FROM_ENVIRONMENT = "Получаем порт из окружения..."
	INFO_USING_DEFAULT_PORT            = "Порт не задан. Будем использовать из конфига - "
	PORT_SET                           = "Порт задан - "
)

type Server struct {
	httpserver *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpserver = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return s.httpserver.ListenAndServe()
}

func EnvPORT(key string) string {
	logrus.Println(INFO_GETTING_PORT_FROM_ENVIRONMENT)
	port := os.Getenv(key)
	if len(port) == 0 {
		port = viper.Get("Port").(string)
		logrus.Println(INFO_USING_DEFAULT_PORT + port)
	} else {
		logrus.Println(PORT_SET + port)
	}
	return port
}
