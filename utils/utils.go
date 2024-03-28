package utils

import (
	"fmt"
	"github.com/OTumanov/go_final_project/settings"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func EnvPORT(key string) string {
	log.Println("Получаем порт из окружения...")
	port := os.Getenv(key)
	if len(port) == 0 {
		log.Println("Порт не задан. Будем использовать 7540")
		port = settings.Port
	} else {
		log.Println("Порт задан -- " + port)
	}
	return ":" + port
}

func EnvDBFILE(key string) string {
	log.Println("Получаем имя БД из окружения...")
	dbName := os.Getenv(key)
	if len(dbName) == 0 {
		log.Println("Имя БД не задано. Будем использовать scheduler.db")
		dbName = settings.DBFile
	} else {
		log.Println("Имя БД задано -- " + dbName)
	}
	return dbName
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	//typeRepeat, err := isValidRepeat(repeat)

	//if err != nil {
	//	return "Ошибка в разборе повтора", err
	//}

	switch repeat[0] {
	case 'd':
		addDays, _ := strconv.Atoi(repeat[2:])
		if addDays < 1 == true || addDays > 400 == true {
			return "Некорректное значение повтора. Допускается от 1 до 400 дней",
				fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
		}

		srtNow := now.Format("20060102")
		compareDate := date

		for srtNow > compareDate || date >= compareDate {
			t, err := time.Parse("20060102", compareDate)
			if err != nil {
				return "Некорректная дата", fmt.Errorf("Не смог разобрать вот это: %s", date)
			}
			compareDate = t.AddDate(0, 0, addDays).Format("20060102")
		}
		return compareDate, nil
	case 'y':
		srtNow := now.Format("20060102")
		compareDate := date

		for srtNow > compareDate || date >= compareDate {
			t, err := time.Parse("20060102", compareDate)
			if err != nil {
				return "Некорректная дата", fmt.Errorf("Не смог разобрать вот это: %s", date)
			}
			compareDate = t.AddDate(1, 0, 0).Format("20060102")
		}
		return compareDate, nil
	case 'w':

		weekdaysStr := strings.TrimPrefix(repeat, "w ")

		//например {"20240126", "w 7", "20240128"}
		if match, _ := regexp.MatchString(`^\d+$`, weekdaysStr); match {
			dayNum, _ := strconv.Atoi(strings.TrimSpace(weekdaysStr))
			if dayNum < 1 || dayNum > 7 {
				return "Некорректное значение повтора. Допускается w <через запятую от 1 до 7>",
					fmt.Errorf("Обрати внимание вот сюда: %s", repeat)

			}
			return NextWeekday(now, date, dayNum)
		}

		//например {"20230126", "w 4,5", "20240201"}, {"20230226", "w 8,4,5", ""}
		weekdaysList := strings.Split(weekdaysStr, ",")
		weekdayMap := make(map[int]bool)

		for _, dayStr := range weekdaysList {
			dayNum, err := strconv.Atoi(strings.TrimSpace(dayStr))
			if err != nil {
				return "Некорректное значение повтора",
					fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
			}
			if dayNum < 1 || dayNum > 7 {
				return "Некорректное значение повтора. Допускается w <через запятую от 1 до 7>",
					fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
			}
			weekdayMap[dayNum] = true
		}
		for i := range weekdayMap {
			findWeekday, err := NextWeekday(now, date, i)
			if err != nil {
				return "Некорректное значение повтора",
					fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
			}
			if findWeekday > date && findWeekday > now.Format("20060102") {
				return findWeekday, nil
			}
		}
		return "Некорректное значение повтора. Допускается w <через запятую от 1 до 7>",
			fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
	case 'm':
		mounthDay, err := searchDayOfMouth(now, date, repeat)
		if err != nil {
			return "Некорректное значение повтора",
				fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
		}
		return mounthDay, err
	}

	return "Что-то пошло не так...", fmt.Errorf("Неправильный формат повтора: %s", repeat)
}

func searchDayOfMouth(now time.Time, date, repeat string) (string, error) {
	monthdaysMap := make(map[int]bool)
	monthsMap := make(map[int]bool)
	dates := make([]time.Time, 0)

	repeatSrt := strings.TrimPrefix(repeat, "m ") //достаем все, после первого пробела -- "m 1,2 1,6" => "1,2 1,6"

	log.Println("repeatSrt=", repeatSrt)

	isNumMonth := strings.Contains(repeatSrt, " ") //если есть пробел, то после него месяцы и их достаем в отдельную мапу
	if isNumMonth {
		log.Println("Есть пробел")
		IndexSep := strings.Index(repeatSrt, " ")
		log.Println("IndexSep=", IndexSep)

		repeatSrtMounth := repeatSrt[IndexSep+1:]
		repeatSrt = repeatSrt[:IndexSep]

		log.Println("repeatSrt=", repeatSrtMounth)
		months := strings.Split(repeatSrtMounth, ",")

		log.Println("months=", months)
		for _, v := range months {
			vi, err := strconv.Atoi(strings.TrimSpace(v))
			log.Println("vi=", vi)
			if err != nil {
				log.Println("err=", err)
				return "Не корректное значение повтора", err
			}
			if vi < 1 || vi > 12 {
				log.Println("err=", err)
				return "Некорректное значение повтора. Допускается m <через запятую от 1 до 12>",
					fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
			}
			monthsMap[vi] = true

			log.Println("Записали в мапу месяц")
		}

		log.Println("monthsMap=", monthsMap)
	}

	monthdays := strings.Split(repeatSrt, ",")

	log.Println("monthdays=", monthdays)
	for _, v := range monthdays {
		vi, err := strconv.Atoi(strings.TrimSpace(v))
		log.Println("vi=", vi)
		if err != nil {
			log.Println("err=", err)
			return "", err
		}
		if vi < -2 || vi > 31 {
			log.Println("err=", err)
			return "Некорректное значение повтора. m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12]",
				fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
		}
		monthdaysMap[vi] = true

		log.Println("Записали в мапу день месяца")
	}

	log.Println("monthdaysMap=", monthdaysMap)

	log.Println("monthsMap=", len(monthsMap))

	log.Println("monthdaysMap=", len(monthdaysMap))
	if len(monthsMap) > 0 {
		for m, _ := range monthsMap {
			for d, _ := range monthdaysMap {
				log.Println("m=", m)
				log.Println("d=", d)
				dates = append(dates, FindDayOfMonthDayWithMount(now, date, m, d))
			}
		}
	} else if len(monthdaysMap) > 0 {
		for d, _ := range monthdaysMap {
			log.Println("d=", d)
			dates = append(dates, FindDayOfMonthOneDate(now, date, d))
		}
	}

	// Проходим по каждому месяцу в указанном периоде
	//day := FindDayOfMonth(now, date, repeat)
	//dates := map[time.Time]bool{day: true}

	log.Println("dates=", dates)
	findNearestDate := findNearestDate(now, date, dates)
	return findNearestDate.Format("20060102"), nil
}

func isValidRepeat(repeat string) (string, error) {
	yearRegexp := regexp.MustCompile(`^y$`)
	simpleRegexp := regexp.MustCompile(`^[d]\s\d+$`)
	advRegexp := regexp.MustCompile(`^[d]\s\d+(,\d+)*$`)
	advWMOneDayRegexp := regexp.MustCompile(`^(w|m)\s\d+$`)
	advMWSomeDaysRegexp := regexp.MustCompile(`^(w|m)\s\d+(,\d+)+$`)
	advMSomeDaysRegexp := regexp.MustCompile(`^m\s(((-?\d+)(,\s?-?\d+)*\s?)+\s?)*$
`)

	if yearRegexp.MatchString(repeat) ||
		simpleRegexp.MatchString(repeat) {
		return "simple", nil
	}
	if advRegexp.MatchString(repeat) ||
		advWMOneDayRegexp.MatchString(repeat) ||
		advMWSomeDaysRegexp.MatchString(repeat) ||
		advMSomeDaysRegexp.MatchString(repeat) {
		return "adv", nil
	}

	return "", fmt.Errorf("Неверный формат повтора: %s", repeat)
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

func FindDayOfMonthOneDate(now time.Time, date string, repeat int) time.Time {
	strNow := now.Format("20060102")
	strSearchDay := date

	if strSearchDay <= strNow && strSearchDay <= date {
		log.Println("Повторять каждое число месяца: ", repeat)

		//searchDay, err := time.Parse("20060102", date)
		searchDay := now
		log.Println("Парсим date в searchDay=", searchDay)
		//if err != nil {
		//	fmt.Println(err)
		//}
		if repeat > 0 {
			searchDay = time.Date(searchDay.Year(), searchDay.Month(), repeat, 0, 0, 0, 0, time.UTC)
			log.Println("устанавливаем год и месяц от date в searchDay и пишем туда же дату повтора -- ", searchDay)
			searchDay = searchDay.AddDate(0, 1, 0)
			log.Println("Добавляем 1 месяц в searchDay=", searchDay)
		}
		if repeat < 0 {
			searchDay = searchDay.AddDate(0, 1, 0)
			log.Println("Добавляем 1 месяц в searchDay=", searchDay)
			searchDay = time.Date(now.Year(), searchDay.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, repeat)
			log.Println("устанавливаем год и месяц от date в searchDay и пишем туда же дату повтора -- ", searchDay)
		}
		strSearchDay = searchDay.Format("20060102")
		log.Println("strSearchDay=", strSearchDay)
	}

	searchDay, err := time.Parse("20060102", strSearchDay)

	if err != nil {
		fmt.Println(err)
	}

	return searchDay
}

func FindDayOfMonthDayWithMount(now time.Time, date string, month, repeat int) time.Time {
	strNow := now.Format("20060102")
	strSearchDay := date

	for strSearchDay <= strNow || strSearchDay <= date {
		log.Println("Повторять каждое число месяца: ", repeat)

		searchDay, err := time.Parse("20060102", date)
		log.Println("Парсим date в searchDay=", searchDay, " и начинаем искать от нее")
		if err != nil {
			fmt.Println(err)
		}

		log.Println("Добавляем 1 месяц в searchDay=", searchDay)
		if repeat > 0 {
			searchDay = time.Date(searchDay.Year(), time.Month(month), repeat, 0, 0, 0, 0, time.UTC)
			log.Println("устанавливаем год и месяц от date в searchDay и пишем туда же дату повтора -- ", searchDay)
		}
		if repeat < 0 {
			searchDay = searchDay.AddDate(0, 1, 0)
			searchDay = time.Date(searchDay.Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, repeat)
			log.Println("устанавливаем год и месяц от date в searchDay и пишем туда же дату повтора -- ", searchDay)
		}
		strSearchDay = searchDay.Format("20060102")
		log.Println("strSearchDay=", strSearchDay)
	}

	searchDay, err := time.Parse("20060102", strSearchDay)

	if err != nil {
		fmt.Println(err)
	}

	return searchDay
}

func NextWeekday(now time.Time, date string, weekday int) (string, error) {
	eventDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("некорректная дата создания события: %v", err)
	}
	currentWeekday := int(now.Weekday())
	daysUntilWeekday := (weekday - currentWeekday + 7) % 7
	nextWeekday := now.AddDate(0, 0, daysUntilWeekday)

	if nextWeekday.Before(eventDate) {
		nextWeekday = eventDate.AddDate(0, 0, (7-currentWeekday+weekday)%7)
	}

	return nextWeekday.Format("20060102"), nil
}
