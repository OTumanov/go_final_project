package utils

import (
	"fmt"
	"github.com/OTumanov/go_final_project/model"
	"log"
	"strings"
	"time"
)

const (
	TITLE_NOT_SET = "Заголовок не может быть пустым!"
)

func checkTask(m *model.Task) (model.Task, error) {
	if strings.TrimSpace(m.Title) == "" {
		log.Println("Заголовок не может быть пустым!")
		return model.Task{}, fmt.Errorf(TITLE_NOT_SET)
	}

	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())

	if m.Date == "" {
		log.Println("Дата не может быть пустой!")
		m.Date = now.Format("20060102")
	}

	_, err := time.Parse("20060102", m.Date)
	if err != nil {
		log.Println("Неверная дата!")
		return model.Task{}, fmt.Errorf("Не могу преобразовать дату!")
	}

	if m.Date < time.Now().Format("20060102") {
		log.Println("Дата не может быть раньше сегодняшней!")
		if m.Repeat == "" {
			m.Date, err = NextDateSearch(time.Now(), m.Date, m.Repeat)
			if err != nil {
				return model.Task{}, err
			}
			log.Println("Новая дата: " + m.Date)
		} else {
			m.Date, err = NextDateSearch(time.Now(), m.Date, m.Repeat)
			if err != nil {
				return model.Task{}, err
			}
			log.Println("Новая дата: " + m.Date)
		}

	}

	return *m, nil
}
