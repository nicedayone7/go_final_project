package models

type Task struct {
	ID string `json:"id"`
	Date string	`json:"date"`
	Title string	`json:"title"`	
	Comment string	`json:"comment"`
	Repeat string	`json:"repeat"`
}

type Sender struct {
	ID string `json:"id"`
	Err string `json:"err"`
}