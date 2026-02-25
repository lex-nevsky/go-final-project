package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// комментирую подробно и для себя в том числе xD

// now - время, от которого ищется ближайшая дата
// dstart - исходное время в формате 20060102, от которого начинается отсчёт повторений
// repeat - правило повторения
// на выходе функции - cтрока с датой в формате "20060102" или ошибка
func NextDate(now time.Time, dstart, repeat string) (string, error) {

	// проверяем подключено ли правило
	if repeat == "" {
		return "", errors.New("repeat: строка пуста")
	}

	// парсим начальную дату dstart
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("dstart: %w", err)
	}

	// разбираем строку правила
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("repeat: строка пуста")
	}
	rule := parts[0]

	// переменные для хранения настроек правила
	var intervalDays int
	var weekDays, monthDays, monthNum map[int]bool

	// парсим разные данные в зависимости от типа правила
	switch rule {

	case "d": // правило переноса на заданное кол-во дней

		// должно быть минимум два слова в правиле
		if len(parts) < 2 {
			return "", errors.New("d: не указан интервал")
		}

		// превращаем второе слово в число
		days, _ := strconv.Atoi(parts[1])

		// вносим ограниечение согласно ТЗ
		if days < 1 || days > 400 {
			return "", errors.New("d: превышен максимальный интервал")
		}
		// сохраняем число для сдвига даты
		intervalDays = days

	case "y": // правило переноса на год, парсить не надо

	case "w": // правило переноса на дни недели

		// должно быть минимум два слова в правиле
		if len(parts) < 2 {
			return "", errors.New("w: не указаны дни недели")
		}

		// ключ - порядковый номер дня недели
		// значение - результат проверки соответствия дня недели с его порядковым номером
		weekDays = make(map[int]bool)

		// разделяем текст по запятым, пустые части пропускаем
		for _, p := range strings.Split(parts[1], ",") {
			if p == "" {
				continue
			}

			d, _ := strconv.Atoi(p)

			// проверяем, что день недели от 1 (понедельник) до 7 (воскресенье)
			if d < 1 || d > 7 {
				return "", errors.New("w: недопустимый день недели")
			}
			// день подходит
			weekDays[d] = true
		}

	case "m": // правило переноса на указанные дни месяца

		// должно быть минимум два слова в правиле
		if len(parts) < 2 {
			return "", errors.New("m: не указаны дни месяца")
		}

		// ключ - указанный день месяца
		// значение - результат проверки соответствия числа в месяце (согласно ТЗ)
		monthDays = make(map[int]bool)

		// разделяем текст по запятым, пустые части пропускаем
		for _, p := range strings.Split(parts[1], ",") {
			if p == "" {
				continue
			}

			// сперва проверяем -1 и -2 (последний и предпоследний день месяца)
			if strings.HasPrefix(p, "-") {
				val, _ := strconv.Atoi(p)

				// проверяем, что нашлись только -1 и -2
				if val < -2 || val > -1 {
					return "", errors.New("m: недопустимый день месяца")
				}

				// день месяца подходит
				monthDays[val] = true

			} else {

				// проверяем дни месяца от 1 до 31
				n, _ := strconv.Atoi(p)
				if n < 1 || n > 31 {
					return "", errors.New("m: недопустимый день месяца")
				}

				// день месяца подходит
				monthDays[n] = true
			}
		}

		// парсим сами месяцы, если они идут третьим параметром

		monthNum = make(map[int]bool)

		// проверяем, что месяц присутствует в данных
		if len(parts) >= 3 {

			// разделяем текст по запятым, пустые части пропускаем
			for _, p := range strings.Split(parts[2], ",") {
				if p == "" {
					continue
				}
				mon, _ := strconv.Atoi(p)

				// проверяем, что месяцы указаны от 1 (январь) до 12 (декабрь)
				if mon < 1 || mon > 12 {
					return "", errors.New("m: недопустимый месяц")
				}
				monthNum[mon] = true
			}
		} else {

			// если конкретные месяцы не указаны,то берем все 12
			for i := 1; i <= 12; i++ {
				monthNum[i] = true
			}
		}

	default:

		// если правило не "d", "y", "w" или "m", то данные с ошибкой
		return "", fmt.Errorf("repeat: неподдерживаемый формат %q", rule)
	}

	// основной цикл, идём от стартовой даты вперёд, пока не найдём подходящую дату
	current := start
	for {

		// сдвигаем дату на один шаг вперёд
		switch rule {

		case "d": // для "d" добавляем intervalDays дней

			current = current.AddDate(0, 0, intervalDays)

		case "y": // для "y" добавляем 1 год

			current = current.AddDate(1, 0, 0)

		default: // для "w" и "m" сдвигаем по одному дню

			current = current.AddDate(0, 0, 1)
		}

		// проверяем, дата строго после now
		// если нет, то продолжаем цикл, сдвигая ещё
		if !current.After(now) {
			continue //
		}

		// проверяем, подходит ли текущая дата под правило
		switch rule {

		case "d", "y":

			// для "d" и "y" любая дата после now подходит, просто возвращаем её
			return current.Format(dateFormat), nil

		case "w":
			// для "w" проверяем день недели
			// в Go 0=воскресенье и 6=суббота, подгоняем под ТЗ, где 1=понедельник и 7=воскресенье
			wday := int(current.Weekday())
			if wday == 0 {
				wday = 7
			}
			// если этот день есть в списке weekDays, то подходит
			if weekDays[wday] {
				return current.Format(dateFormat), nil
			}

		case "m":
			// для "m" проверяем сначала месяц
			// если месяц не подходит, то пропускаем эту дату
			mon := int(current.Month())
			if !monthNum[mon] {
				continue
			}

			// вычисляем, сколько дней в текущем месяце
			daysInMonth := time.Date(current.Year(), time.Month(mon+1), 0, 0, 0, 0, 0, time.UTC).Day()

			// проверяем день месяца
			for k := range monthDays {
				if k > 0 {
					// обычный день (1..31)
					if current.Day() == k {
						return current.Format(dateFormat), nil
					}
				} else {
					// -1 последний день месяца (daysInMonth)
					// -2 предпоследний день месяца (daysInMonth - 1)
					if current.Day() == daysInMonth+k+1 {
						return current.Format(dateFormat), nil
					}
				}
			}
		}
	}
}
