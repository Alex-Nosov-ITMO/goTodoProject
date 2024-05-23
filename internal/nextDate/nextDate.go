package nextDate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Количество дней в месяце
var (
	dayM = map[int]int{
		1:  31,
		2:  29,
		3:  31,
		4:  30,
		5:  31,
		6:  30,
		7:  31,
		8:  31,
		9:  30,
		10: 31,
		11: 30,
		12: 31,
	}
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	// Парсим дату из строки в объект типа time.Time
	t, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("невозможно преобразовать дату: %s. Ошибка: %v", date, err)
	}

	// Проверяем формат правила повторения
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	// Разбиваем правило повторения на части
	parts := strings.Split(repeat, " ")

	// Обрабатываем базовые правила
	switch parts[0] {
	// Ежедневно(интервалы)
	case "d":

		if len(parts) != 2 {
			return "", errors.New("неверное количество аргументов правила повторения")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		// Проверка на корректность формата правила повторения
		if days < 1 || days > 400 {
			return "", fmt.Errorf("некорректный формат правила повторения: %s, должно быть целое число от 1 до 400", repeat)
		}

		// Прибавляем к дате нужное количество дней
		t = t.AddDate(0, 0, days)

		// Проверяем, больше ли дата, чем сегодняшняя
		for checkTime(t, now) {
			t = t.AddDate(0, 0, days)
		}

	// Ежегодно
	case "y":

		if len(parts) != 1 {
			return "", errors.New("неверное количество аргументов правила повторения")
		}

		t = t.AddDate(1, 0, 0)

		// Проверяем, больше ли дата, чем сегодняшняя
		for checkTime(t, now) {
			t = t.AddDate(1, 0, 0)
		}

	// Еженедельно
	case "w":

		if len(parts) != 2 {
			return "", errors.New("неверное количество аргументов правила повторения")
		}

		// Разделяем дни недели
		daysString := strings.Split(parts[1], ",")

		// Проверка на корректность формата правила повторения
		days := make([]int, 0)
		for _, v := range daysString {
			d, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			if d < 1 || d > 7 {
				return "", fmt.Errorf("некорректный формат правила повторения: %s, должно быть целое число от 1 до 7", repeat)
			}

			if !contains(&days, d) {
				days = append(days, d)
			} else {
				return "", fmt.Errorf("некорректный формат правила повторения: %s, введены дублирующиеся дни недели", repeat)
			}
		}

		// Прибавляем к дате нужное количество дней
		for {
			t = t.AddDate(0, 0, 1)

			// Переводим день недели в индекс
			checkIndex := int(t.Weekday())
			if checkIndex == 0 {
				checkIndex = 7
			}

			// Проверяем, день недели входит в правило
			if !contains(&days, checkIndex) {
				continue
			}

			// Проверяем, больше ли дата, чем сегодняшняя
			if !checkTime(t, now) {
				break
			}
		}

	// Ежемесячно
	case "m":
		if len(parts) != 2 && len(parts) != 3 {
			return "", errors.New("неверное количество аргументов правила повторения")
		}

		var monthsString, daysMonthString []string

		// Проверка на корректность формата правила повторения
		daysMonthString = strings.Split(parts[1], ",")
		daysMonth := make([]int, 0)
		for _, v := range daysMonthString {
			d, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			if !((d >= 1 && d <= 31) || (d == -1) || (d == -2)) {
				return "", fmt.Errorf("некорректный формат правила повторения: %s, должно быть целое число от 1 до 31 или -1, -2", repeat)
			}

			if !contains(&daysMonth, d) {
				daysMonth = append(daysMonth, d)
			} else {
				return "", fmt.Errorf("некорректный формат правила повторения: %s, введены дублирующиеся дни месяца", repeat)
			}
		}

		// Проверка на корректность формата правила повторения
		months := make([]int, 0)
		if len(parts) == 3 {
			monthsString = strings.Split(parts[2], ",")
			for _, v := range monthsString {
				m, err := strconv.Atoi(v)
				if err != nil {
					return "", err
				}
				if m < 1 || m > 12 {
					return "", fmt.Errorf("некорректный формат правила повторения: %s, должно быть целое число месяцев от 1 до 12", repeat)
				}

				if !contains(&months, m) {
					months = append(months, m)
				} else {
					return "", fmt.Errorf("некорректный формат правила повторения: %s, введены дублирующиеся месяцы", repeat)
				}
			}

			// Обработка случая, когда введенных дней нет в месяце (Например m 31 6,8)
			ok := false
		loof:
			for _, d := range months {
				for _, m := range daysMonth {
					if m < dayM[d] {
						ok = true
						break loof
					}
				}
			}

			if !ok {
				return "", fmt.Errorf("некорректный формат правила повторения: %s, не существует такого дня в введенных месяцах", repeat)
			}
		}

		// Прибавляем к дате нужное количество дней
		for {
			t = t.AddDate(0, 0, 1)

			// Проверяем год на високосность
			if !((t.Year()%4 == 0 && t.Year()%100 != 0) || t.Year()%400 == 0) {
				dayM[2] = 28
			} else {
				dayM[2] = 29
			}

			// Проверяем, день месяца входит ли в правило
			if !contains(&daysMonth, t.Day()) {
				if contains(&daysMonth, -2) && contains(&daysMonth, -1) {
					if t.Day() != dayM[int(t.Month())]-1 && t.Day() != dayM[int(t.Month())] {
						continue
					}
				} else if contains(&daysMonth, -2) {
					if t.Day() != dayM[int(t.Month())]-1 {
						continue
					}

				} else if contains(&daysMonth, -1) {
					if t.Day() != dayM[int(t.Month())] {
						continue
					}
				} else {
					continue
				}
			}

			// Проверяем, месяц входит ли в правило
			if len(parts) == 3 {
				if !contains(&months, int(t.Month())) {
					continue
				}
			}

			// Проверяем, больше ли дата, чем сегодняшняя
			if !checkTime(t, now) {
				break
			}
		}

	// Некорректный формат правила повторения
	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}

	return t.Format("20060102"), nil
}

// Проверка на наличие элемента в срезе
func contains(arr *[]int, target int) bool {
	for _, num := range *arr {
		if num == target {
			return true
		}
	}
	return false
}

// Проверка на меньше ли дата, чем сегодняшняя
func checkTime(t, now time.Time) bool {
	return t.Before(now) || t.Equal(now)
}
