package model

type Task struct {
	ID      int    `json:"-"`
	Date    string `json:"date"`
	Title   string `json:"title" binding:"required"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat" binding:"required"`
	Done    bool   `json:"-"`
}
