package db

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

type Task struct {
	ID      string `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

// добавляем задачу в БД и возвращаем её id
func AddTask(task *Task) (string, error) {

	// проверяем инит БД
	if db == nil {
		return "", errors.New("db: не инициализирована")
	}

	// делаем запись в БД
	res, err := db.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return "", err
	}

	// получаем id записи
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	// преобразуем id в строку
	return strconv.FormatInt(id, 10), nil
}

// возвращаем до limit задач, отсортированных по дате
func Tasks(limit int, search string) ([]*Task, error) {

	// проверяем инит БД
	if db == nil {
		return nil, errors.New("db: не инициализирована")
	}

	// если search пустой, то выбираем все задачи
	if search == "" {
		rows, err := db.Query(
			"SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?",
			limit,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanTasks(rows)
	}

	// распознаем search как дату
	t, err := time.Parse("02.01.2006", search)
	if err != nil {

		// если search не дата, то ищем по тексту (title или comment)
		pattern := "%" + search + "%"
		rows, err := db.Query(
			"SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?",
			pattern, pattern, limit,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanTasks(rows)
	}

	// если search это валидная дата, то ищем по точной дате
	rows, err := db.Query(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?",
		t.Format("20060102"),
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// сканируем строки из rows, если их нет, то слайс пустой
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	defer rows.Close()

	// собираем задачи в слайс
	var tasks []*Task
	for rows.Next() {
		var t Task
		// считываем поля из строки
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	// проверяем ошибки
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// возвращаем задачу по id или ошибку
func GetTask(id string) (*Task, error) {

	// проверяем инит БД
	if db == nil {
		return nil, errors.New("db: не инициализирована")
	}

	// выполняем запрос к одной строке
	var t Task
	err := db.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	// если строка не найдена, то возвращаем ошибку
	if err == sql.ErrNoRows {
		return nil, errors.New("задача не найдена")
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// обновляем задачу и возвращаем ошибку, если задача не найдена
func UpdateTask(task *Task) error {

	// проверяем инит БД
	if db == nil {
		return errors.New("db: не инициализирована")
	}

	// парсим id
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("id: неверный формат")
	}

	// обновляем запись
	res, err := db.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, id,
	)
	if err != nil {
		return err
	}

	// проверяем, что строка обновлена
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

// удаляем задачу по id
func DeleteTask(id string) error {

	// проверяем инит БД
	if db == nil {
		return errors.New("db: не инициализирована")
	}

	// реализуем удаление
	res, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	// проверяем, что строка удалена
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

// обновляем только дату задачи
func UpdateDate(next, id string) error {

	// проверяем инит БД
	if db == nil {
		return errors.New("db: не инициализирована")
	}

	// обновляем только дату
	res, err := db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", next, id)
	if err != nil {
		return err
	}

	// проверяем, что строка обновлена
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}
