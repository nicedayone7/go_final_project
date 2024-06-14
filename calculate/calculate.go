package calculate

import (
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	mounths := map[int]string{
		1:  time.January.String(),
		2:  time.February.String(),
		3:  time.March.String(),
		4:  time.April.String(),
		5:  time.May.String(),
		6:  time.June.String(),
		7:  time.July.String(),
		8:  time.August.String(),
		9:  time.September.String(),
		10: time.October.String(),
		11: time.November.String(),
		12: time.December.String(),
	}
	var result string
	
	role, days := parseDay(repeat)
	t, _ := strconv.Atoi(days)
	switch {
	case role == "d" && t <= 400:
		result = dateFormater(now, t, "20060102")
	case role == "y":
		result = dateFormater(now, t, "20060102")
	case role == "w" && 1 <= t && t <= 7:
		parseDay(repeat)
		nextDay, err := nextDayWeek(strings.Split(days, ","), int(now.Weekday()))
		if err != nil {
			return "", err
		}
		nextDate := now.Local().AddDate(0, 0, nextDay).Format("20060102")
		return nextDate, nil
	case role == "m":
		
		

	}

	return result, nil
}

func dateFormater(now time.Time, addDays int, format string) string {
	date := now.Local().AddDate(0, 0, addDays).Format(format)
	return date
}

func parseDay(dayRepeat string) (string, string) {
	role, days := strings.Split(dayRepeat, " ")[0], strings.Split(dayRepeat, " ")[1]
	return role, days
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
