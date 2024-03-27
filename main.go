package main

import (
	"net/http"
	"os"
)

const webDir = "./web"

func listenPort(key string) string {
	port := ":" + os.Getenv(key)
	return port
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(listenPort("TODO_PORT"), nil)
	if err != nil {
		panic(err)
	}
}
