package model

type NextDate struct {
	Now    string `form:"now"`
	Date   string `form:"date"`
	Repeat string `form:"repeat"`
}
