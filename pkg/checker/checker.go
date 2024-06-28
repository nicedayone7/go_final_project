package task

import (
	"errors"
	"fmt"
	calc "go_final_project/calculate"
	"go_final_project/pkg/models"
	"time"
)

func Task(task models.Task) (models.Task, error) {
	fmt.Println(task.Date, task.Title, task.Comment, task.Repeat)

	if task.Title == "" {
		err := errors.New("empty title fail write task")
		fmt.Println(err)
		return models.Task{}, err
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		fmt.Println(err)
		return models.Task{}, err
	}

	if task.Date == "" || (task.Date == "" && task.Repeat == "") {
		task.Date = time.Now().String()
		return task, nil
	} 

	date, _ := time.Parse("20060102", task.Date)

	if date.Before(time.Now()) {
		if task.Repeat == "" {
			task.Date = time.Now().String()
			return task, nil
		}

		task.Date, err = calc.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			fmt.Println(err)
			return models.Task{}, err
		}
		return task, nil
	}

	return task, nil
}