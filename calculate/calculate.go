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
			err := errors.New("does not match the rules format")
			return "", err
		}
		nextDateTask, err := nextDayCalc(now, startDate, daysRepeat)
		if err != nil {
			return "", err
		}
		resultDate, err := nextDateTask.Format(dateFormat)
		if err != nil {
			return "", nil
		}
		return resultDate, nil

	// 	result = dateFormater(now, t, "20060102")
	// case role == "y":
	// 	result = dateFormater(now, t, "20060102")
	// case role == "w" && 1 <= t && t <= 7:
	// 	parseDay(repeat)
	// 	nextDay, err := nextDayWeek(strings.Split(days, ","), int(now.Weekday()))
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	nextDate := now.Local().AddDate(0, 0, nextDay).Format("20060102")
	// 	return nextDate, nil
	// case role == "m":
		
		

	}

	return result, nil
}

func nextDayCalc(now time.Time, date time.Time, daysRepeat []int) (string, error) {
	dayRepeat := daysRepeat[0]
	if dayRepeat > 400 {
		err := errors.New("the number must not exceed 400")
		return "", err
	}

	if dayRepeat == 1 {
		return now.Format(dateFormat), nil
	}

	diffDays := now.Sub(date).Hours()/24 // result different now and start task date
	remainingDays := int(diffDays)%dayRepeat
	
	if remainingDays != 0 {
		countRepeat := int(diffDays)/dayRepeat
		dayToNextTask := countRepeat*(dayRepeat+1)
		resultDate := date.AddDate(0, 0, dayToNextTask)
		return resultDate.Format(dateFormat), nil
	}

	return now.Format(dateFormat), nil
}


func dateFormater(dateFormat time.Time, format string) string {
	date := dateFormat.Format(format)
	return date
}

func parseRepeat(dayRepeat string) (string, []int, error) {
	var formatDays []int
	role, days := strings.Split(dayRepeat, " ")[0], strings.Split(dayRepeat, " ")[1]
	daysList := strings.Split(days, ",")
	for _, day := range daysList {
		day, err := strconv.Atoi(day)
		if err != nil {
			return "", "", err
		}
		formatDays = append(formatDays, day)
	}
	return role, formatDays, nil
}

func nextDayWeek(days []string, nowDay int) (int, error) {
	var result int
	var previousDay int
	for i, day := range days {
		changeDay, err := strconv.Atoi(day)
		if err != nil {
			return 0, err
		}
		
		if nowDay > changeDay {
			result = 7 - nowDay + changeDay
		} else {
			result = nowDay - changeDay
		}
		
		if i == 0 {
			previousDay = result
			continue
		} else {
			if previousDay < result {
				previousDay = result
			}
		}
	}
	
	return result, nil
}
