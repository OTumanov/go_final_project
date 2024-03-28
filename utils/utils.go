package utils

import (
	"fmt"
	"github.com/OTumanov/go_final_project/settings"
	"log"
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
	typeRepeat, err := isValidRepeat(repeat)

	if err != nil {
		return "Ошибка в разборе повтора", err
	}

	if typeRepeat == "simple" {
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
		}
	}
	if typeRepeat == "adv" {
		switch repeat[0] {
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
		}
		return "", nil
	}

	return "Что-то пошло не так...", fmt.Errorf("Неправильный формат повтора: %s", repeat)
}

func isValidRepeat(repeat string) (string, error) {
	yearRegexp := regexp.MustCompile(`^y$`)
	simpleRegexp := regexp.MustCompile(`^[d]\s\d+$`)
	advRegexp := regexp.MustCompile(`^[d]\s\d+(,\d+)*$`)
	advWOneDayRegexp := regexp.MustCompile(`^w\s\d+$`)
	advWSomeDaysRegexp := regexp.MustCompile(`^w\s\d+(,\d+)+$`)

	if yearRegexp.MatchString(repeat) || simpleRegexp.MatchString(repeat) {
		return "simple", nil
	}
	if advRegexp.MatchString(repeat) || advWOneDayRegexp.MatchString(repeat) || advWSomeDaysRegexp.MatchString(repeat) {
		return "adv", nil
	}

	return "", fmt.Errorf("Неверный формат повтора: %s", repeat)
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
