package calculate

import (
	"errors"
	"strings"
	"time"
)

func nextDayCalc(now time.Time, date time.Time, dayRepeat int) (time.Time, error) {
	var resultDate time.Time
	if dayRepeat == 1 {
		if date.After(now) {
			resultDate = date.AddDate(0, 0, 1)
		} else {
			resultDate = now.AddDate(0, 0, 1)
		}
		return resultDate, nil
	}

	diffDays := now.Sub(date).Hours()/24 // result different now and start task date
	remainingDays := int(diffDays)%dayRepeat
	
	if remainingDays != 0 {
		countRepeat := int(diffDays)/dayRepeat
		dayToNextTask := (countRepeat+1)*dayRepeat
		resultDate := date.AddDate(0, 0, dayToNextTask)
		return resultDate, nil
	}
	
	return resultDate, nil
}

func nextYearCalc(now time.Time,startDate time.Time) time.Time {
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
		return 0, nil
	}

	for i, day := range days {		
		if nowDay > day {
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
	
	return result, nil
}

func nextDayMonth(now time.Time, parts []string, sortDays, sortMonths []int) (time.Time, error) {
	if len(parts) == 2 {
		for _, day := range sortDays {
			if now.Day() < day {
				nextDateTask := time.Date(now.Year(), now.Month(), day, 0,0,0,0,time.UTC)
				return nextDateTask, nil
			}
		}
		
		nextDateTask := time.Date(now.Year(), now.Month()+1, sortDays[0], 0, 0, 0, 0, time.UTC)
		return nextDateTask, nil
	}

	if len(parts) == 3 {
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

	err := errors.New("fail to find next date task")
	return time.Time{}, err	
}