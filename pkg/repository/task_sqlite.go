package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/OTumanov/go_final_project/pkg/model"

	"github.com/jmoiron/sqlx"
)

const (
	ERR_INVALID_REPEAT_VALUE                 = "Некорректное значение повтора."
	ERR_REPEAT                               = "Обрати внимание вот сюда: %v."
	ALLOWED_DAYS_RANGE                       = "Допускается от 1 до 400 дней."
	ALLOWED_WEEKDAYS_RANGE                   = "Допускается w <через запятую от 1 до 7>."
	ALLOWED_DAYS_AND_MONTH_RANGE             = "m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12]"
	INVALID_DATE_MESSAGE                     = "Некорректная дата."
	UNABLE_TO_CONVERT_REPEAT_TO_NUMBER_ERROR = "Не удалось сконвертировать повтор в число."
	ERROR_PARSE_DATE                         = "Не смог разобрать вот это: %s."
	MONTH_NUMBER_RANGE                       = "Допускается m <через запятую от 1 до 12>."
	WRONG_DATE                               = "некорректная дата создания события: %v"
	WRONG_REPEAT                             = "Проблемы с форматом повтора, погляди сюда повтор: %v"
	TASK_NAME_REQUIRED_ERROR_MESSAGE         = "Название задачи не может быть пустым."
	TASK_NOT_FOUND_ERROR_MESSAGE             = "Задача не найдена"
	TASK_ID_REQUIRED_ERROR_MESSAGE           = "Не указан идентификатор"
	INVALID_TASK_ID_ERROR_MESSAGE            = "Некорректный идентификатор"
	PREFIX_REPEAT_M_                         = "m "
	PREFIX_REPEAT_D_                         = "d "
	PREFIX_REPEAT_W_                         = "w "
	PREFIX_DAY                               = 'd'
	PREFIX_WEEK                              = 'w'
	PREFIX_MONTH                             = 'm'
	PREFIX_YEAR                              = 'y'
	SEPARATOR_SPACE                          = " "
	SEPARATOR_COMMA                          = ","
	DATE_FORMAT_YYYYMMDD                     = `20060102`
	DATE_FORMAT_YYYY_MM_DD                   = `02.01.2006`
	FIRST_DAY                                = 1
	MINUS_ONE_DAY                            = -1
	ADDING_ONE_MOUNTH                        = 1
	MIN_REPEAT_INTERVAL_DAYS                 = 1
	MAX_REPEAT_INTERVAL_DAY                  = 31
	MIN_MINUS_REPEAT_INTERVAL_DAY            = -2
	MAX_REPEAT_INTERVAL_DAYS                 = 400
	MIN_MONTHS                               = 1
	MAX_MONTHS                               = 12
	MIN_WEEK                                 = 1
	MAX_WEEK                                 = 7
	ONE_WEEK
	LIMIT_TASKS = 25
)

type TodoTaskSqlite struct {
	db *sqlx.DB
}

func NewTodoTaskSqlite(db *sqlx.DB) *TodoTaskSqlite {
	return &TodoTaskSqlite{db: db}
}
func (t *TodoTaskSqlite) NextDate(nd model.NextDate) (string, error) {
	if nd.Repeat == "" {
		return "", fmt.Errorf(WRONG_REPEAT, nd.Repeat)
	}

	if !regexp.MustCompile(`^([wdm]\s.*|y)?$`).MatchString(nd.Repeat) {
		return "", fmt.Errorf(WRONG_REPEAT, nd.Repeat)
	}

	switch nd.Repeat[0] {
	case PREFIX_DAY:
		repeatIntervalDays, err := findRepeatIntervalDays(nd)
		return repeatIntervalDays, err
	case PREFIX_YEAR:
		repeatIntervalYears, err := findRepeatIntervalYears(nd)
		return repeatIntervalYears, err
	case PREFIX_WEEK:
		repeatIntervalWeeks, err := findRepeatIntervalWeeks(nd)
		return repeatIntervalWeeks, err
	case PREFIX_MONTH:
		repeatIntervalMonths, err := findRepeatIntervalMonths(nd)
		return repeatIntervalMonths, err
	default:
		return "", nil
	}
}
func (t *TodoTaskSqlite) CreateTask(task model.Task) (int64, error) {
	err := t.checkTask(&task)
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf("INSERT INTO %s (title, comment, date, repeat) VALUES ($1, $2, $3, $4) RETURNING id", taskTable)
	row := t.db.QueryRow(query, task.Title, task.Comment, task.Date, task.Repeat)

	var id int64
	if err = row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
func (t *TodoTaskSqlite) GetTasks(search string) (model.ListTasks, error) {
	var tasks []model.Task
	var query string

	switch typeSearch(search) {
	case 0:
		query = fmt.Sprintf("SELECT * FROM %s ORDER BY date LIMIT ?", taskTable)
		err := t.db.Select(&tasks, query, LIMIT_TASKS)
		if err != nil {
			return model.ListTasks{}, err
		}
	case 1:
		s, _ := time.Parse(DATE_FORMAT_YYYY_MM_DD, search)
		st := s.Format(DATE_FORMAT_YYYYMMDD)
		query = fmt.Sprintf("SELECT * FROM %s WHERE date = ? ORDER BY date LIMIT ?", taskTable)
		err := t.db.Select(&tasks, query, st, LIMIT_TASKS)
		if err != nil {
			return model.ListTasks{}, err
		}
	case 2:
		searchQuery := fmt.Sprintf("%%%s%%", search)
		query := `SELECT * FROM scheduler WHERE LOWER(title) LIKE ? OR LOWER(comment) LIKE ? ORDER BY date LIMIT ?`
		rows, err := t.db.Queryx(query, searchQuery, searchQuery, LIMIT_TASKS)
		if err != nil {
			return model.ListTasks{}, err
		}
		for rows.Next() {
			var task model.Task
			err := rows.StructScan(&task)
			if err != nil {
				return model.ListTasks{}, err
			}
			tasks = append(tasks, task)
		}
	}

	if len(tasks) == 0 {
		return model.ListTasks{Tasks: []model.Task{}}, nil
	}
	return model.ListTasks{Tasks: tasks}, nil
}
func (t *TodoTaskSqlite) GetTaskById(id string) (model.Task, error) {
	if id == "" {
		return model.Task{}, fmt.Errorf(TASK_ID_REQUIRED_ERROR_MESSAGE)
	}
	if _, err := strconv.Atoi(id); err != nil {
		return model.Task{}, fmt.Errorf(INVALID_TASK_ID_ERROR_MESSAGE)
	}
	var task model.Task
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", taskTable)
	err := t.db.Get(&task, query, id)
	if err != nil {
		return model.Task{}, fmt.Errorf(TASK_NOT_FOUND_ERROR_MESSAGE)
	}
	return task, err
}
func (t *TodoTaskSqlite) UpdateTask(task model.Task) error {
	err := t.checkTask(&task)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET title = ?, comment = ?, date = ?, repeat = ? WHERE id = ?", taskTable)
	_, err = t.db.Exec(query, task.Title, task.Comment, task.Date, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf(TASK_NOT_FOUND_ERROR_MESSAGE)
	}
	return nil
}
func (t *TodoTaskSqlite) DeleteTask(id string) error {
	_, err := t.GetTaskById(id)
	if err != nil {
		return err
	}
	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE id = ?", taskTable)
	_, err = t.db.Exec(queryDelete, id)
	if err != nil {
		return err
	}
	return nil
}
func (t *TodoTaskSqlite) TaskDone(id string) error {
	task, err := t.GetTaskById(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		queryDeleteTask := fmt.Sprintf("DELETE FROM %s WHERE id = ?", taskTable)
		logrus.Println(queryDeleteTask)
		t.db.Exec(queryDeleteTask, id)
		return nil
	}

	nd := model.NextDate{
		Date:   task.Date,
		Now:    time.Now().Format(DATE_FORMAT_YYYYMMDD),
		Repeat: task.Repeat,
	}

	newDate, err := t.NextDate(nd)
	if err != nil {
		return err
	}

	task.Date = newDate
	queryUpdateTask := fmt.Sprintf("UPDATE %s SET date = ? WHERE id = ?", taskTable)
	logrus.Println(queryUpdateTask)
	_, err = t.db.Exec(queryUpdateTask, task.Date, id)
	if err != nil {
		return err
	}
	return nil

}
func (t *TodoTaskSqlite) checkTask(task *model.Task) error {
	if task.Title == "" {
		return fmt.Errorf(TASK_NAME_REQUIRED_ERROR_MESSAGE)
	}

	if !regexp.MustCompile(`^([wdm]\s.*|y)?$`).MatchString(task.Repeat) {
		return fmt.Errorf(WRONG_REPEAT, task.Repeat)
	}

	now := time.Now().Format(DATE_FORMAT_YYYYMMDD)

	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse(DATE_FORMAT_YYYYMMDD, task.Date)
	if err != nil {
		return fmt.Errorf(INVALID_DATE_MESSAGE)
	}

	if task.Date < now {
		if task.Repeat == "" {
			task.Date = now
		}
		if task.Repeat != "" {
			nd := model.NextDate{
				Date:   task.Date,
				Now:    now,
				Repeat: task.Repeat,
			}
			task.Date, err = t.NextDate(nd)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func findRepeatIntervalDays(nd model.NextDate) (string, error) {
	now, err := timeNow(nd)
	if err != nil {
		return "", err
	}
	stringRepeatIntervalDays := strings.TrimPrefix(nd.Repeat, PREFIX_REPEAT_D_)
	repeatIntervalDays, err := strconv.Atoi(stringRepeatIntervalDays)
	if err != nil {
		return ERR_INVALID_REPEAT_VALUE, fmt.Errorf(ERR_REPEAT, nd.Repeat)
	}
	if repeatIntervalDays < MIN_REPEAT_INTERVAL_DAYS || repeatIntervalDays > MAX_REPEAT_INTERVAL_DAYS {
		return ERR_INVALID_REPEAT_VALUE + ALLOWED_DAYS_RANGE, fmt.Errorf(ERR_REPEAT, nd.Repeat)
	}
	searchDate := nd.Date

	for searchDate <= now.Format(DATE_FORMAT_YYYYMMDD) || searchDate <= nd.Date {
		d, err := time.Parse(DATE_FORMAT_YYYYMMDD, searchDate)
		if err != nil {
			return INVALID_DATE_MESSAGE, fmt.Errorf(ERROR_PARSE_DATE, nd.Date)
		}
		searchDate = d.AddDate(0, 0, repeatIntervalDays).Format(DATE_FORMAT_YYYYMMDD)
	}
	return searchDate, nil
}
func findRepeatIntervalMonths(nd model.NextDate) (string, error) {
	now, err := timeNow(nd)
	if err != nil {
		return "", err
	}
	repeatSrt := strings.TrimPrefix(nd.Repeat, PREFIX_REPEAT_M_)
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
					fmt.Errorf(ERR_REPEAT, nd.Repeat)
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
			return ERR_INVALID_REPEAT_VALUE + ALLOWED_DAYS_AND_MONTH_RANGE, fmt.Errorf(ERR_REPEAT, nd.Repeat)
		}
		monthDays = append(monthDays, vi)
	}
	nextDates := make([]time.Time, 0)

	if len(months) > 0 {
		for i := 0; i < len(months); i++ {
			m := months[i]
			for j := 0; j < len(monthDays); j++ {
				d := monthDays[j]
				nd := findDayOfMonth(now, nd.Date, m, d)
				nextDates = append(nextDates, nd)
			}
		}
	} else if len(monthDays) > 0 {
		for _, d := range monthDays {
			nextDates = append(nextDates, findDayOfMonth(now, nd.Date, 0, d))
		}
	}

	findNearestDate := findNearestDate(now, nd.Date, nextDates)
	return findNearestDate.Format(DATE_FORMAT_YYYYMMDD), nil
}
func findRepeatIntervalWeeks(nd model.NextDate) (string, error) {
	now, err := timeNow(nd)
	if err != nil {
		return "", err
	}

	weekdayNumber := strings.TrimPrefix(nd.Repeat, PREFIX_REPEAT_W_)
	weekDaysSlice := strings.Split(weekdayNumber, SEPARATOR_COMMA)

	for _, v := range weekDaysSlice {
		vi, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return UNABLE_TO_CONVERT_REPEAT_TO_NUMBER_ERROR, err
		}
		if vi < MIN_WEEK || vi > MAX_WEEK {
			return ERR_INVALID_REPEAT_VALUE + ALLOWED_WEEKDAYS_RANGE, fmt.Errorf(ERR_REPEAT, nd.Repeat)
		}
		findWeekday, err := findNextWeekDay(now, nd.Date, vi)
		if err != nil {
			return ERR_INVALID_REPEAT_VALUE, fmt.Errorf(ERR_REPEAT, nd.Repeat)
		}
		if findWeekday > nd.Date && findWeekday > now.Format(DATE_FORMAT_YYYYMMDD) {
			return findWeekday, nil
		}
	}
	return ERR_INVALID_REPEAT_VALUE + ALLOWED_WEEKDAYS_RANGE, fmt.Errorf(ERR_REPEAT, nd.Repeat)
}
func findRepeatIntervalYears(nd model.NextDate) (string, error) {
	now, err := timeNow(nd)
	if err != nil {
		return "", err
	}

	formattedNow := now.Format(DATE_FORMAT_YYYYMMDD)
	searchDate := nd.Date

	for searchDate <= formattedNow || searchDate <= nd.Date {
		d, err := time.Parse(DATE_FORMAT_YYYYMMDD, searchDate)
		if err != nil {
			return INVALID_DATE_MESSAGE, fmt.Errorf(ERR_REPEAT, nd.Date)
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
func timeNow(nd model.NextDate) (time.Time, error) {
	var now time.Time
	if nd.Now == "" {
		now = time.Now()
	}
	now, err := time.Parse(DATE_FORMAT_YYYYMMDD, nd.Now)
	if err != nil {
		return time.Time{}, fmt.Errorf(ERROR_PARSE_DATE, nd.Now)
	}
	return now, nil
}
func typeSearch(str string) int {
	if str == "" {
		return 0
	}
	_, err := time.Parse(DATE_FORMAT_YYYY_MM_DD, str)
	if err == nil {
		return 1
	}
	return 2
}
