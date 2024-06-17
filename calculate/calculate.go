package calculate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "20060102"
)
// mounths := map[int]string{
// 	1:  time.January.String(),
// 	2:  time.February.String(),
// 	3:  time.March.String(),
// 	4:  time.April.String(),
// 	5:  time.May.String(),
// 	6:  time.June.String(),
// 	7:  time.July.String(),
// 	8:  time.August.String(),
// 	9:  time.September.String(),
// 	10: time.October.String(),
// 	11: time.November.String(),
// 	12: time.December.String(),
// }

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", nil
	}
	
	role, daysRepeat, err := parseRepeat(repeat)
	if err != nil {
		return "", nil
	}

	switch {
	case role == "d":
		if len(daysRepeat) != 1 {
			err := errors.New("does not match the rules format for days")
			return "", err
		}
		nextDateTask, err := nextDayCalc(now, startDate, daysRepeat)
		if err != nil {
			return "", err
		}
		resultDate := nextDateTask.Format(dateFormat)
		if err != nil {
			return "", nil
		}
		return resultDate, nil

	
	case role == "y":
	 	if len(repeat) != 1 {
			err := errors.New("does not match the rules format for year")
			return "", err
		}

		if !startDate.After(now) {
			nextYearTask := nextYearCalc(now, startDate)
			return nextYearTask.Format(dateFormat), nil
		}

		return startDate.Format(dateFormat), nil

	case role == "w":
		if len(daysRepeat) < 1 || len(daysRepeat)  > 7 {
			err := errors.New("does not match the rules format for week: lenght repeat rule more 7")
				return "", err
		}
		addDay, err :=	nextDayWeek(daysRepeat, int(now.Weekday()))
		if err != nil {
			return "", err
		}
		nextWeekTask := now.AddDate(0,0,addDay)
		return nextWeekTask.Format(dateFormat), nil
		
	case role == "m":

		return "", nil	
		

	// 	if len()
	// 	nextDay, err := nextDayWeek(strings.Split(days, ","), int(now.Weekday()))
	// 	if err != nil {
	// 		return "", err
	// 	}
		
	// case role == "m":
		
		

	}

	return "", nil
}

func nextDayCalc(now time.Time, date time.Time, daysRepeat []int) (time.Time, error) {
	dayRepeat := daysRepeat[0]
	if dayRepeat >= 400 || dayRepeat <= 0 {
		err := errors.New("the number must not exceed 400")
		return time.Time{}, err
	}

	if dayRepeat == 1 {
		return now, nil
	}

	diffDays := now.Sub(date).Hours()/24 // result different now and start task date
	remainingDays := int(diffDays)%dayRepeat
	
	if remainingDays != 0 {
		countRepeat := int(diffDays)/dayRepeat
		dayToNextTask := countRepeat*(dayRepeat+1)
		resultDate := date.AddDate(0, 0, dayToNextTask)
		return resultDate, nil
	}

	return now, nil
}

func nextYearCalc(now time.Time, startDate time.Time) time.Time {
	for !startDate.After(now) {
		startDate = startDate.AddDate(1, 0, 0)
	}
	return startDate
}


func parseRepeat(dayRepeat string) (string, []int, error) {
	var formatDays []int
	role, days := strings.Split(dayRepeat, " ")[0], strings.Split(dayRepeat, " ")[1]
	daysList := strings.Split(days, ",")
	for _, day := range daysList {
		day, err := strconv.Atoi(day)
		if err != nil {
			return "", nil, err
		}
		formatDays = append(formatDays, day)
	}
	return role, formatDays, nil
}

func nextDayWeek(days []int, nowDay int) (int, error) {
	for _, day := range days {
		if day < 1 || day > 7 {
			err := errors.New("does not match the rules format for week: day is incorrect out of range")
			return 0, err
		}
	}
	

	var result int
	var previousDay int
	if days[0] == nowDay {
		return nowDay, nil
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
