package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type NextDate struct {
	date   string
	repeat string
	want   string
}

func mmm(now time.Time, nd NextDate) (string, error) {
	repeatSrt := strings.TrimPrefix(nd.repeat, "m ")
	log.Println("repeatSrt=", repeatSrt)
	isConstainsNumMonth := strings.Contains(repeatSrt, " ")
	log.Println("isConstainsNumMonth=", isConstainsNumMonth)
	monthsSlice := make([]string, 0)
	months := make([]int, 0)
	if isConstainsNumMonth {
		IndexSep := strings.Index(repeatSrt, " ")
		repeatSrtMounth := repeatSrt[IndexSep+1:]
		repeatSrt = repeatSrt[:IndexSep]
		monthsSlice = strings.Split(repeatSrtMounth, ",")

		for _, v := range monthsSlice {
			vi, err := strconv.Atoi(strings.TrimSpace(v))
			if err != nil {
				return "Не корректное значение повтора", err
			}
			if vi < 1 || vi > 12 {
				return "Некорректное значение повтора. Допускается m <через запятую от 1 до 12>",
					fmt.Errorf("Обрати внимание вот сюда: %s", nd.repeat)
			}
			months = append(months, vi)
		}
	}

	monthDaysSlice := strings.Split(repeatSrt, ",")
	log.Println("monthDays=", monthDaysSlice)

	monthDays := make([]int, 0)

	for _, v := range monthDaysSlice {
		vi, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return "Не удалось сконвертировать повтор в число", err
		}
		if vi < -2 || vi > 31 {
			return "Некорректное значение повтора. m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12]",
				fmt.Errorf("Обрати внимание вот сюда: %s", nd.repeat)
		}
		monthDays = append(monthDays, vi)
	}

	log.Println("months=", months)
	log.Println("monthDays=", monthDays)

	nextDates := make([]time.Time, 0)

	if len(months) > 0 {
		for _, m := range months {
			for _, d := range months {
				nextDates = append(nextDates, searchDayOfMonth(now, nd.date, m, d))
			}
		}
	} else if len(monthDays) > 0 {
		for _, d := range monthDays {
			log.Println("d=", d)
			nextDates = append(nextDates, searchDayOfMonth(now, nd.date, 0, d))
		}
	}

	log.Println("nextDates=", nextDates)
	findNearestDate := findNearestDate(now, nd.date, nextDates)
	return findNearestDate.Format("20060102"), nil
	//return "", nil
}

func findNearestDate(now time.Time, date string, dates []time.Time) time.Time {
	if len(dates) == 1 {
		return dates[0]
	}

	var nearestDate time.Time
	dateTodate, err := time.Parse("20060102", date)
	if err != nil {
		fmt.Println(err)
	}
	minDifference := math.MaxInt64

	for _, d := range dates {
		if d.After(now) && d.After(dateTodate) {
			difference := int(d.Sub(now).Seconds())
			if difference < minDifference {
				minDifference = difference
				nearestDate = d
			}
		}
	}
	return nearestDate
}

func searchDayOfMonth(now time.Time, date string, month, repeat int) time.Time {
	maxDate, err := time.Parse("20060102", date)
	if err != nil {
		log.Fatal(err)
	}
	if now.After(maxDate) {
		maxDate = now
	}

	log.Println(maxDate)
	if month == 0 {
		month = int(maxDate.Month())
	}

	searchDay := time.Date(maxDate.Year(), time.Month(month), repeat, 0, 0, 0, 0, time.UTC)
	//log.Println("searchDay=", searchDay)
	daysInThisMonth := daysInMonth(maxDate.Year(), time.Month(month))
	if repeat > daysInThisMonth {
		searchDay = searchDay.AddDate(0, +1, repeat)
	}
	log.Println("searchDay=", searchDay)

	return searchDay
}

func daysInMonth(year int, month time.Month) int {
	// Создаем дату с первым днем следующего месяца
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	// Переходим на предыдущий день, чтобы получить последний день текущего месяца
	lastDayOfMonth := nextMonth.AddDate(0, 0, -1)
	return lastDayOfMonth.Day()
}

func main() {
	now, _ := time.Parse("20060102", "20240126")
	searchWeekday, _ := mmm(now, NextDate{"20240329", "m 10,17 12,8,1", "20240810"})
	fmt.Println(searchWeekday)

	//{"20231106", "m 13", "20240213"},
	//{"20240120", "m 40,11,19", ""},
	//{"20240116", "m 16,5", "20240205"},
	//{"20240126", "m 25,26,7", "20240207"},
	//{"20240409", "m 31", "20240531"},
	//{"20240329", "m 10,17 12,8,1", "20240810"},
	//{"20230311", "m 07,19 05,6", "20240507"},
	//{"20230311", "m 1 1,2", "20240201"},
	//{"20240127", "m -1", "20240131"},
	//{"20240222", "m -2", "20240228"},
	//{"20240222", "m -2,-3", ""},
	//{"20240326", "m -1,-2", "20240330"},
	//{"20240201", "m -1,18", "20240218"},
	//{"20240125", "w 1,2,3", "20240129"},
	//{"20240126", "w 7", "20240128"},
	//{"20230126", "w 4,5", "20240201"},
	//{"20230226", "w 8,4,5", ""},
}
