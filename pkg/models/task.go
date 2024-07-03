package models

import (
	"errors"
	"time"

	"go_final_project/pkg/nextdate"
)
const dateFormat = "20060102"

type Task struct {
	ID string `json:"id"`
	Date string	`json:"date"`
	Title string	`json:"title"`	
	Comment string	`json:"comment"`
	Repeat string	`json:"repeat"`
}

func (t Task) Check() (Task, error) {
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),0,0,0,0, time.UTC)
	
	if t.Title == "" {
		err := errors.New("empty title fail write task")
		return Task{}, err
	}

	if t.Date == "" || (t.Date == "" && t.Repeat == "") {
		t.Date = time.Now().Format(dateFormat)
		return t, nil
	} 

	startDate, err := time.Parse(dateFormat, t.Date)
	
	if err != nil {
		return Task{}, err
	}

	if startDate.Before(now) {
		
		if t.Repeat == "" {
			t.Date = time.Now().Format(dateFormat)
			return t, nil
		}

		t.Date, err = nextdate.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return Task{}, err
		}
		return t, nil
	}

	return t, nil
}


