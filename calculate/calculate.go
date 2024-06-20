package calculate

import (
	"errors"
	"strings"
	"time"
)

const (
	dateFormat = "20060102"
)

var ruleLettersTask = []string{"d", "y", "w", "m"}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", nil
	}
	
	ruleLetter, err := parseRepeat(repeat)
	if err != nil {
		return "", err
	}

	switch {
	case ruleLetter == "d":
		dayRepeat, err := normilizeDays(repeat)
		if err != nil {
			return "", err
		}
		
		nextDateTask, err := nextDayCalc(now, startDate, dayRepeat)
		if err != nil {
			return "", err
		}
		
		return nextDateTask.Format(dateFormat), nil

	case ruleLetter == "y":
	 	err := normilizeYear(repeat)
		if err != nil {
			return "", err
		}

		nextYearTask := nextYearCalc(now, startDate)
		
		return nextYearTask.Format(dateFormat), nil

	case ruleLetter == "w":
		daysWeek, err := normilizeDaysWeek(repeat)
		if err != nil {
			return "", err
		}
		
		addDays, err :=	nextDayWeekCalc(daysWeek, int(now.Weekday()))
		if err != nil {
			return "", err
		}
		nextWeekTask := now.AddDate(0, 0, addDays)
		return nextWeekTask.Format(dateFormat), nil
		
	case ruleLetter == "m":
		parts := strings.Split(repeat, " ")
		if len(parts) > 3 || len(parts) < 2 {
			err := errors.New("invalid input data value")
			return "", err
		}
		months, err := normilizeMonths(repeat)
		if err != nil {
			return "", err
		}
		daysMonth, err := normilizeDaysForMonth(now, repeat)
		if err != nil {
			return "", err
		}
		nextMonthTask, err := nextDayMonth(now, repeat, parts, daysMonth, months)
		if err != nil {
			return "", err
		}
		return nextMonthTask.Format(dateFormat), nil
			
	default:
		return "", nil
	}
}