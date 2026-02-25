package db

import (
	"database/sql"

	_ "modernc.org/sqlite" // драйвер
)

var db *sql.DB

// создаем правило генерации таблицы scheduler и индекс по date
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8) NOT NULL DEFAULT '',
    title   VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT,
    repeat  VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// открываем БД по пути dbFile и создаём таблицу и индекс, если их нет
func Init(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if _, err = db.Exec(schema); err != nil {
		db.Close()
		db = nil
		return err
	}

	return nil
}
