package utils

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ERR_INVALID_REPEAT_VALUE = "Некорректное значение повтора."
	ERR_REPEAT               = "Обрати внимание вот сюда: %v."
	//ERR_REPEAT_NOT_SET                       = "Не задан повтор."
	ALLOWED_DAYS_RANGE                       = "Допускается от 1 до 400 дней."
	ALLOWED_WEEKDAYS_RANGE                   = "Допускается w <через запятую от 1 до 7>."
	ALLOWED_DAYS_AND_MONTH_RANGE             = "m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12]"
	INVALID_DATE_MESSAGE                     = "Некорректная дата."
	UNABLE_TO_CONVERT_REPEAT_TO_NUMBER_ERROR = "Не удалось сконвертировать повтор в число."
	ERROR_PARSE_DATE                         = "Не смог разобрать вот это: %s."
	MONTH_NUMBER_RANGE                       = "Допускается m <через запятую от 1 до 12>."
	WRONG_DATE                               = "некорректная дата создания события: %v"
	PREFIX_REPEAT_M_                         = "m "
	PREFIX_REPEAT_D_                         = "d "
	PREFIX_REPEAT_W_                         = "w "
	PREFIX_DAY                               = 'd'
	PREFIX_WEEK                              = 'w'
	PREFIX_MONTH                             = 'm'
	PREFIX_YEAR                              = 'y'
	SEPARATOR_SPACE                          = " "
	SEPARATOR_COMMA                          = ","
	DATE_FORMAT_YYYYMMDD                     = "20060102"
	FIRST_DAY                                = 1
	MINUS_ONE_DAY                            = -1
	ADDING_ONE_MOUNTH                        = 1
	//MIN_REPEAT_LEN                           = 3
	MIN_REPEAT_INTERVAL_DAYS      = 1
	MAX_REPEAT_INTERVAL_DAY       = 31
	MIN_MINUS_REPEAT_INTERVAL_DAY = -2
	MAX_REPEAT_INTERVAL_DAYS      = 400
	MIN_MONTHS                    = 1
	MAX_MONTHS                    = 12
	MIN_WEEK                      = 1
	MAX_WEEK                      = 7
	ONE_WEEK
)

type NextDate struct {
	date   string
	repeat string
	want   string
}

func NextDateSearch(now time.Time, date, repeat string) (string, error) {

	if repeat == "" {
		return time.Now().Format(DATE_FORMAT_YYYYMMDD), nil
	}
	if !regexp.MustCompile(`^(w|d|m|y)(\s+\S*)?$`).MatchString(repeat) {
		return "не прокатило", fmt.Errorf("херовый повтор", repeat)
	}

	switch repeat[0] {
	case PREFIX_DAY:
		repeatIntervalDays, err := findRepeatIntervalDays(now, NextDate{date, repeat, ""})
		return repeatIntervalDays, err
	case PREFIX_YEAR:
		repeatIntervalYears, err := findRepeatIntervalYears(now, NextDate{date, repeat, ""})
		return repeatIntervalYears, err
	case PREFIX_WEEK:
		repeatIntervalWeeks, err := findRepeatIntervalWeeks(now, NextDate{date, repeat, ""})
		return repeatIntervalWeeks, err
	case PREFIX_MONTH:
		repeatIntervalMonths, err := findRepeatIntervalMonths(now, NextDate{date, repeat, ""})
		return repeatIntervalMonths, err
	default:
		return "", nil
	}
}

func findRepeatIntervalDays(now time.Time, nd NextDate) (string, error) {
	stringRepeatIntervalDays := strings.TrimPrefix(nd.repeat, PREFIX_REPEAT_D_)
	repeatIntervalDays, err := strconv.Atoi(stringRepeatIntervalDays)
	if err != nil {
		return ERR_INVALID_REPEAT_VALUE, fmt.Errorf(ERR_REPEAT, nd.repeat)
	}
	if repeatIntervalDays < MIN_REPEAT_INTERVAL_DAYS || repeatIntervalDays > MAX_REPEAT_INTERVAL_DAYS {
		return ERR_INVALID_REPEAT_VALUE + ALLOWED_DAYS_RANGE, fmt.Errorf(ERR_REPEAT, nd.repeat)
	}
	searchDate := nd.date

	for searchDate <= now.Format(DATE_FORMAT_YYYYMMDD) || searchDate <= nd.date {
		d, err := time.Parse(DATE_FORMAT_YYYYMMDD, searchDate)
		if err != nil {
			return INVALID_DATE_MESSAGE, fmt.Errorf(ERROR_PARSE_DATE, nd.date)
		}
		searchDate = d.AddDate(0, 0, repeatIntervalDays).Format(DATE_FORMAT_YYYYMMDD)
	}
	return searchDate, nil
}
func findRepeatIntervalMonths(now time.Time, nd NextDate) (string, error) {
	repeatSrt := strings.TrimPrefix(nd.repeat, PREFIX_REPEAT_M_)
	isConstainsNumMonth := strings.Contains(repeatSrt, SEPARATOR_SPACE)
	monthsSlice := make([]string, 0)
	months := make([]int, 0)
	if isConstainsNumMonth {
		IndexSep := strings.Index(repeatSrt, SEPARATOR_SPACE)
		repeatSrtMounth := repeatSrt[IndexSep+1:]
		repeatSrt = repeatSrt[:IndexSep]
		monthsSlice = strings.Split(repeatSrtMounth, SEPARATOR_COMMA)

		for _, v := range monthsSlice {
			vi, err := strconv.Atoi(strings.TrimSpace(v))
			if err != nil {
				return ERR_INVALID_REPEAT_VALUE, err
			}
			if vi < MIN_MONTHS || vi > MAX_MONTHS {
				return ERR_INVALID_REPEAT_VALUE + MONTH_NUMBER_RANGE,
					fmt.Errorf(ERR_REPEAT, nd.repeat)
			}
			months = append(months, vi)
		}
	}

	monthDaysSlice := strings.Split(repeatSrt, SEPARATOR_COMMA)

	monthDays := make([]int, 0)

	for _, v := range monthDaysSlice {
		vi, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return UNABLE_TO_CONVERT_REPEAT_TO_NUMBER_ERROR, err
		}
		if vi < MIN_MINUS_REPEAT_INTERVAL_DAY || vi > MAX_REPEAT_INTERVAL_DAY {
			return ERR_INVALID_REPEAT_VALUE + ALLOWED_DAYS_AND_MONTH_RANGE, fmt.Errorf(ERR_REPEAT, nd.repeat)
		}
		monthDays = append(monthDays, vi)
	}
	nextDates := make([]time.Time, 0)

	if len(months) > 0 {
		for i := 0; i < len(months); i++ {
			m := months[i]
			for j := 0; j < len(monthDays); j++ {
				d := monthDays[j]
				nd := findDayOfMonth(now, nd.date, m, d)
				nextDates = append(nextDates, nd)
			}
		}
	} else if len(monthDays) > 0 {
		for _, d := range monthDays {
			nextDates = append(nextDates, findDayOfMonth(now, nd.date, 0, d))
		}
	}

	findNearestDate := findNearestDate(now, nd.date, nextDates)
	return findNearestDate.Format(DATE_FORMAT_YYYYMMDD), nil
}
func findRepeatIntervalWeeks(now time.Time, nd NextDate) (string, error) {
	weekdayNumber := strings.TrimPrefix(nd.repeat, PREFIX_REPEAT_W_)
	weekDaysSlice := strings.Split(weekdayNumber, SEPARATOR_COMMA)

	for _, v := range weekDaysSlice {
		vi, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return UNABLE_TO_CONVERT_REPEAT_TO_NUMBER_ERROR, err
		}
		if vi < MIN_WEEK || vi > MAX_WEEK {
			return ERR_INVALID_REPEAT_VALUE + ALLOWED_WEEKDAYS_RANGE, fmt.Errorf(ERR_REPEAT, nd.repeat)
		}
		findWeekday, err := findNextWeekDay(now, nd.date, vi)
		if err != nil {
			return ERR_INVALID_REPEAT_VALUE, fmt.Errorf(ERR_REPEAT, nd.repeat)
		}
		if findWeekday > nd.date && findWeekday > now.Format(DATE_FORMAT_YYYYMMDD) {
			return findWeekday, nil
		}
	}
	return ERR_INVALID_REPEAT_VALUE + ALLOWED_WEEKDAYS_RANGE, fmt.Errorf(ERR_REPEAT, nd.repeat)
}
func findRepeatIntervalYears(now time.Time, nd NextDate) (string, error) {
	formattedNow := now.Format(DATE_FORMAT_YYYYMMDD)
	searchDate := nd.date

	for searchDate <= formattedNow || searchDate <= nd.date {
		d, err := time.Parse(DATE_FORMAT_YYYYMMDD, searchDate)
		if err != nil {
			return INVALID_DATE_MESSAGE, fmt.Errorf(ERR_REPEAT, nd.date)
		}
		searchDate = d.AddDate(1, 0, 0).Format(DATE_FORMAT_YYYYMMDD)
	}
	return searchDate, nil
}
func findNextWeekDay(now time.Time, date string, weekday int) (string, error) {
	eventDate, err := time.Parse(DATE_FORMAT_YYYYMMDD, date)
	if err != nil {
		return "", fmt.Errorf(WRONG_DATE, err)
	}
	currentWeekday := int(now.Weekday())
	daysUntilWeekday := (weekday - currentWeekday + ONE_WEEK) % ONE_WEEK
	nextWeekday := now.AddDate(0, 0, daysUntilWeekday)

	if nextWeekday.Before(eventDate) {
		nextWeekday = eventDate.AddDate(0, 0, (ONE_WEEK-currentWeekday+weekday)%ONE_WEEK)
	}

	return nextWeekday.Format(DATE_FORMAT_YYYYMMDD), nil
}
func findDayOfMonth(now time.Time, date string, month, repeat int) time.Time {
	var searchDay time.Time
	maxDate, err := time.Parse(DATE_FORMAT_YYYYMMDD, date)
	if err != nil {
		log.Fatal(err)
	}

	if maxDate.Before(now) {
		maxDate = now
	}

	if month == 0 {
		month = int(maxDate.Month())
	}

	lastDayOfMonth := lastDayMonth(maxDate.Year(), time.Month(month))
	if repeat > lastDayOfMonth {
		searchDay = time.Date(maxDate.Year(), time.Month(month+ADDING_ONE_MOUNTH), repeat, 0, 0, 0, 0, time.UTC)
	} else if repeat < lastDayOfMonth && repeat > 0 {
		searchDay = time.Date(maxDate.Year(), time.Month(month), repeat, 0, 0, 0, 0, time.UTC)
	} else if repeat < 0 {
		searchDay = time.Date(maxDate.Year(), time.Month(month), lastDayOfMonth+ADDING_ONE_MOUNTH, 0, 0, 0, 0, time.UTC)
		searchDay = searchDay.AddDate(0, 0, repeat)
	}

	if searchDay.Before(maxDate) {
		searchDay = searchDay.AddDate(0, ADDING_ONE_MOUNTH, 0)
	}

	return searchDay
}
func lastDayMonth(year int, month time.Month) int {
	nextMonth := time.Date(year, month+ADDING_ONE_MOUNTH, FIRST_DAY, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := nextMonth.AddDate(0, 0, MINUS_ONE_DAY)
	return lastDayOfMonth.Day()
}
func findNearestDate(now time.Time, date string, dates []time.Time) time.Time {
	if len(dates) == 1 {
		return dates[0]
	}

	var nearestDate time.Time
	dateToDate, err := time.Parse(DATE_FORMAT_YYYYMMDD, date)
	if err != nil {
		fmt.Println(err)
	}
	minDifference := math.MaxInt64

	for _, d := range dates {
		if d.After(now) && d.After(dateToDate) {
			difference := int(d.Sub(now).Seconds())
			if difference < minDifference {
				minDifference = difference
				nearestDate = d
			}
		}
	}
	return nearestDate
}
