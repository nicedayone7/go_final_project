package wraper

import (
	"go_final_project/pkg/models"
	"sort"
	"time"
)

func SortedTasks(tasks []models.Task) []models.Task {
	sort.Slice(tasks, func(i, j int) bool {
		firstDate, _ := time.Parse("20060102", tasks[i].Date) 
		lastDate, _ := time.Parse("20060102", tasks[j].Date)
		return firstDate.Before(lastDate)
	})
	return tasks
}