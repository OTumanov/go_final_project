package utils

import (
	"fmt"
	"github.com/OTumanov/go_final_project/settings"
	"log"
	"os"
	"regexp"
	"strconv"
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

	if typeRepeat == "year" {
		srtNow := now.Format("20060102")
		nextDate := date

		for srtNow > nextDate || date >= nextDate {
			t, err := time.Parse("20060102", nextDate)
			if err != nil {
				return "Некорректная дата", fmt.Errorf("Не смог разобрать вот это: %s", date)
			}
			nextDate = t.AddDate(1, 0, 0).Format("20060102")
		}
		return nextDate, nil
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
		case 'm':
			addMonths, _ := strconv.Atoi(repeat[2:])
			if addMonths < 1 || addMonths > 12 {
				return "Некорректное значение повтора. Допускается от 1 до 12 месяцев", fmt.Errorf("некорректное значение повтора: %s", repeat)
			}
			next, err := time.Parse("20060102", date)
			if err != nil {
				return "Некорректная дата", fmt.Errorf("Не смог разобрать вот это: %s", date)
			}
			for next.Before(now) {
				next = next.AddDate(0, addMonths, 0)
			}
			return next.Format("20060102"), nil
		}
	}
	if typeRepeat == "adv" {
		return now.Format("20060102"), nil
	}

	return "Что-то пошло не так...", fmt.Errorf("Неправильный формат повтора: %s", repeat)
}

func isValidRepeat(repeat string) (string, error) {
	year := regexp.MustCompile(`^y$`)
	simple := regexp.MustCompile(`^[wdm]\s\d+$`)
	adv := regexp.MustCompile(`^[wdm]\s\d+(,\d+)*$`)

	if year.MatchString(repeat) {
		return "year", nil
	}
	if simple.MatchString(repeat) {
		return "simple", nil
	}

	if adv.MatchString(repeat) {
		return "adv", nil
	}

	return "", fmt.Errorf("Неверный формат повтора: %s", repeat)
}
