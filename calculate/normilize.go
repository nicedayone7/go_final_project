package calculate

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

func normilizeDays(rule string) (int, error) {
	parts := strings.Split(rule, " ")
	if len(parts) != 2 {
		err := errors.New("invalid input data value")
		return 0, err
	}

	part := strings.Split(rule, " ")[1]
	dayRepeat, err := strconv.Atoi(part)
	if err != nil {
		return 0, err
	}

	if dayRepeat > 400 || dayRepeat < 1 {
		err := errors.New("invalid input data value")
		return 0, err
	}

	return dayRepeat, nil
}

func normilizeYear(rule string) (error) {
	if len(rule) != 1 {
		err := errors.New("invalid input data value")
		return err
	}
	return nil
}

func normilizeDaysWeek(rule string) ([]int, error) {
	var week []int
	parts := strings.Split(rule, " ")
	daysWeek := strings.Split(parts[1], ",")

	if len(daysWeek) < 1 || len(daysWeek) > 7 {
		err := errors.New("invalid input data value")
		return nil, err
	}

	for _, dayWeek := range daysWeek {
		day, err := strconv.Atoi(dayWeek)
		if err != nil {
			return nil, err
		}
		if day < 1 || day > 7 {
			err := errors.New("invalid input data value")
			return nil, err
		}
		week = append(week, day)
	}

	return week, nil
}

func normilizeDaysForMonth(now time.Time,rule string) ([]int, error) {
	var daysForSort []int
	part := strings.Split(rule, " ")[1]
	daysMounth := strings.Split(part, ",")
	for _, dayStr := range daysMounth {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return nil, err
		}

		if day > 31 || day < -2 {
			err := errors.New("invalid input data value")
			return nil, err
		}

		if day == -1 || day == -2 {
			firstDayNextMonth := time.Date(now.Year(),now.Month()+1, day+1,0,0,0,0,time.UTC)
			daysForSort = append(daysForSort, firstDayNextMonth.Day())
			continue
		}
		daysForSort = append(daysForSort, day)
	}
	sort.Ints(daysForSort)
	return daysForSort, nil
}

func normilizeMonths(rule string) ([]int, error) {
	var monthsForSort []int
	part := strings.Split(rule, " ")[2]
	months := strings.Split(part, ",")
	for _, m := range months {
		month, err := strconv.Atoi(m)
		if err != nil {
			return nil, err
		}
		if month > 12 || month < 1 {
			err := errors.New("invalid input data value")
			return nil, err
		}
		monthsForSort = append(monthsForSort, month)	
	}
	sort.Ints(monthsForSort)

	return monthsForSort, nil
}

