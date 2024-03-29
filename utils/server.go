package utils

import (
	"net/http"
	"time"
)

const (
	webDir       = "./web"
	ENV_PORT     = "TODO_PORT"
	API_NEXTDATE = "/api/nextdate"
	NOW          = "now"
	DATE_EVENT   = "date"
	REPEAT_EVENT = "repeat"
)

func StartHTTPServer() {
	listenPort := EnvPORT(ENV_PORT)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.Handle(API_NEXTDATE, http.HandlerFunc(nextDate))
	err := http.ListenAndServe(listenPort, nil)
	if err != nil {
		panic(err)
	}
}

func nextDate(w http.ResponseWriter, r *http.Request) {
	nowTime, err := time.Parse(DATE_FORMAT_YYYYMMDD, r.URL.Query().Get(NOW))
	if err != nil {
		return
	}
	nextDate, err := NextDateSearch(nowTime, r.URL.Query().Get(DATE_EVENT), r.URL.Query().Get(REPEAT_EVENT))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
