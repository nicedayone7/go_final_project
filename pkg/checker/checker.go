package task

import (
	"errors"
	"fmt"
	calc "go_final_project/pkg/calculate"
	"go_final_project/pkg/models"
	"time"
)

func Task(task models.Task) (models.Task, error) {
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),0,0,0,0, time.UTC)
	fmt.Println(task.Date, task.Title, task.Comment, task.Repeat)

	if task.Title == "" {
		err := errors.New("empty title fail write task")
		return models.Task{}, err
	}

	if task.Date == "" || (task.Date == "" && task.Repeat == "") {
		task.Date = time.Now().Format("20060102")
		return task, nil
	} 

	date, err := time.Parse("20060102", task.Date)
	fmt.Println(date)
	if err != nil {
		return models.Task{}, err
	}

	if date.Before(now) {
		
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
			return task, nil
		}

		task.Date, err = calc.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return models.Task{}, err
		}
		return task, nil
	}

	return task, nil
}