package calculate

import (
	"errors"
	"strings"
	"time"

	"go_final_project/pkg/normilize"
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
		dayRepeat, err := normilize.Days(repeat)
		if err != nil {
			return "", err
		}
		
		nextDateTask, err := nextDayCalc(now, startDate, dayRepeat)
		if err != nil {
			return "", err
		}
		
		return nextDateTask.Format(dateFormat), nil

	case ruleLetter == "y":
	 	err := normilize.Year(repeat)
		if err != nil {
			return "", err
		}

		nextYearTask := nextYearCalc(now, startDate)
		
		return nextYearTask.Format(dateFormat), nil

	case ruleLetter == "w":
		daysWeek, err := normilize.DaysWeek(repeat)
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
		if now.Before(startDate) {
			now = startDate
		}
		parts := strings.Split(repeat, " ")
		if len(parts) > 3 || len(parts) < 2 {
			err := errors.New("invalid input data value")
			return "", err
		}
		if len(parts) == 2 {
			normilizeDays, err := normilize.DaysForMonth(now, parts[1])
			if err != nil {
				return "", err
			}
			nextMonthTask, err := nextDayWithDays(now, normilizeDays)
			if err != nil {
				return "", err
			}
			return nextMonthTask.Format(dateFormat), nil
		}
	
		if len(parts) == 3 {
			normilizeDays, err := normilize.DaysForMonth(now, parts[1])
			if err != nil {
				return "", err
			}
			months, err := normilize.Months(parts[2])
			if err != nil {
				return "", err
			}
			nextMonthTask, err := nextDayWithDaysMonths(now, normilizeDays, months)
			if err != nil {
				return "", err
			}
			return nextMonthTask.Format(dateFormat), nil
		}

		return "", err
			
	default:
		err := errors.New("not found next date")
		return "", err
	}
}