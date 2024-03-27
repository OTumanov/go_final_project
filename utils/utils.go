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
		next, err := time.Parse("20060102", date)
		compDate, err := time.Parse("20060102", date)

		if err != nil {
			return "Некорректная дата", fmt.Errorf("Не смог разобрать вот это: %s", date)
		}

		fmt.Println("ДО -- next=", next, "compDate=", compDate)

		for next.Before(compDate) {
			next = next.AddDate(1, 0, 0)
			fmt.Println("После -- next=", next, "compDate=", compDate)
		}
		return next.Format("20060102"), nil
	}

	if typeRepeat == "simple" {
		switch repeat[0] {
		case 'd': //если повтор в днях
			addDays, _ := strconv.Atoi(repeat[2:])
			if addDays < 1 || addDays > 400 {
				return "Некорректное значение повтора. Допускается от 1 до 400 дней", fmt.Errorf("Обрати внимание вот сюда: %s", repeat)
			}
			next, err := time.Parse("20060102", date)
			if err != nil {
				return "Некорректная дата", fmt.Errorf("Не смог разобрать вот это: %s", date)
			}
			for next.Before(now) {
				next = next.AddDate(0, 0, addDays)
			}
			return next.Format("20060102"), nil

		case 'm': //если повтор в месяцах
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

	return "", fmt.Errorf("неправильный формат повтора: %s", repeat)
}
