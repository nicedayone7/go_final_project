package calculate

import (
	"errors"
	"strings"
	"time"
)

func nextDayCalc(now time.Time, date time.Time, dayRepeat int) (time.Time, error) {

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	if now.Before(date) {
		return date.AddDate(0, 0, dayRepeat), nil
	}

	diff := date.Sub(now)
	daysDifference := int(diff.Hours() / 24)

	if daysDifference >= 0 {
		if daysDifference == 0 {
			return date.AddDate(0, 0, dayRepeat), nil
		} else if daysDifference < dayRepeat {
			return date.AddDate(0, 0, dayRepeat), nil
		} else {
			countRepeats := daysDifference / dayRepeat
			daysToAdd := (countRepeats + 1) * dayRepeat
			return date.AddDate(0, 0, daysToAdd), nil
		}
	} 
	// result different now and start task date
	diffDays := now.Sub(date).Hours() / 24 
	remainingDays := int(diffDays) % dayRepeat

	if remainingDays != 0 {
		countRepeat := int(diffDays) / dayRepeat
		dayToNextTask := (countRepeat + 1) * dayRepeat
		resultDate := date.AddDate(0, 0, dayToNextTask)
		return resultDate, nil
	}
	
	return now, nil
}

func nextYearCalc(now time.Time, startDate time.Time) time.Time {
	for {
		startDate = startDate.AddDate(1, 0, 0)
		if startDate.After(now) {
			break
		}
	}
	return startDate
}

func parseRepeat(dayRepeat string) (string, error) {
	ruleLetter := strings.Split(dayRepeat, " ")[0]
	for _, rule := range ruleLettersTask {
		if rule == ruleLetter {
			return ruleLetter, nil
		}
	}
	err := errors.New("invalid character")
	return "", err
}

func nextDayWeekCalc(days []int, nowDay int) (int, error) {
	var result int
	var previousDay int

	if days[0] == nowDay {
		return 7, nil
	}

	for i, day := range days {
		if nowDay >= day {
			result = 7 - nowDay + day
		} else {
			result = day - nowDay
		}

		if i == 0 {
			previousDay = result
			continue
		}

		if previousDay > result {
			previousDay = result
		}
	}

	return previousDay, nil
}

func nextDayWithDays(now time.Time, sortDays []int) (time.Time, error) {
	if sortDays[0] == 31 {
		specMonth := lastDayMonthCheck(now)
		nextDateTask := time.Date(now.Year(), time.Month(specMonth), sortDays[0], 0, 0, 0, 0, time.UTC)
		return nextDateTask, nil
	}

	for _, day := range sortDays {
		if now.Day() < day {
			nextDateTask := time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, time.UTC)
			return nextDateTask, nil
		}
	}
	nextDateTask := time.Date(now.Year(), now.Month()+1, sortDays[0], 0, 0, 0, 0, time.UTC)
	return nextDateTask, nil
}

func nextDayWithDaysMonths(now time.Time, sortDays, sortMonths []int) (time.Time, error) {
	for _, month := range sortMonths {
		for _, day := range sortDays {
			findDate := time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
			if findDate.After(now) {
				return findDate, nil
			}
		}
	}
	findDate := time.Date(now.Year()+1, time.Month(sortMonths[0]), sortDays[0], 0, 0, 0, 0, time.UTC)
	return findDate, nil
}

func lastDayMonthCheck(date time.Time) int {
	specMonths := [7]int{1, 3, 5, 7, 8, 10, 12}
	for _, month := range specMonths {
		if int(date.Month()) <= month {
			return month
		}
	}
	return 0
}
