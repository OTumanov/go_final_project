package model

type Task struct {
	ID      int    `json:"-"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
	Done    bool   `json:"-"`
}
