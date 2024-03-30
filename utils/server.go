package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OTumanov/go_final_project/model"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	webDir       = "./web"
	ENV_PORT     = "TODO_PORT"
	API_NEXTDATE = "/api/nextdate"
	API_TASK     = "/api/task"
	NOW          = "now"
	DATE_EVENT   = "date"
	REPEAT_EVENT = "repeat"
)

type ErrorResponse struct {
	Message string `json:"error"`
}

type TaskIdResponse struct {
	Id string `json:"id"`
}

func StartHTTPServer() {
	listenPort := EnvPORT(ENV_PORT)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.Handle(API_NEXTDATE, http.HandlerFunc(nextDate))
	http.Handle(API_TASK, http.HandlerFunc(task))
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
func task(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//taskGet(w, r)
	case http.MethodPost:
		taskPost(w, r)
	//case http.MethodDelete:
	//	taskDel(w, r)
	//case http.MethodPut:
	//	taskPut(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func taskPost(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	task, err = checkTask(&task)
	if err != nil {
		responseError := ErrorResponse{Message: err.Error()}
		jsonResponse, err := json.Marshal(responseError)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	lastId, err := addingTask(getDB(), task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(lastId)
	w.WriteHeader(http.StatusOK)

	responseTaskId := TaskIdResponse{Id: strconv.Itoa(int(lastId))}
	jsonResponse, err := json.Marshal(responseTaskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

//func taskGet(w http.ResponseWriter, r *http.Request) {
//	// TODO
//}
//
//func taskPut(w http.ResponseWriter, r *http.Request) {
//	// TODO
//}
//
//func taskDel(w http.ResponseWriter, r *http.Request) {
//	// TODO
//}
