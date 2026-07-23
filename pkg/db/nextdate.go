package db

import (
    "fmt"
    "time"
    "strings"
    "strconv"
)

// формат даты используем как константу
const DateFormat = "20060102"

// отдельная функция для проверки месяца и даты
func isValidMonthAndAfter(date time.Time, validMonths map[int]bool, now time.Time) bool {
    month := int(date.Month())
    return (len(validMonths) == 0 || validMonths[month]) && date.After(now)
}

// NextDate public func. - вычисляет ближайшую дату задачи по правилу repeat.
// now      - текущее время, от которого ищем следующую дату
// dstart   - исходная дата задачи в формате DateFormat
// repeat   - строка с правилом повторения ("d 7", "y", "w 1,3", "m -1,2" и т.п.)
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
    if repeat == "" {
        // Правило не указано - задача на один раз, после выполнения удаляется
        return "", fmt.Errorf("repeat rule is empty")
    }

    // Парсим исходную дату
    date, err := time.Parse(DateFormat, dstart)
    if err != nil {
        return "", fmt.Errorf("invalid dstart: %w", err)
    }

    // Разбираем правило по пробелам
    parts := strings.Split(repeat, " ")

    switch parts[0] {
    case "d":
        // Правило "d - число" — перенос на указанное количество дней
        if len(parts) < 2 {
            return "", fmt.Errorf("invalid d rule")
        }
        interval, err := strconv.Atoi(parts[1])
        if err != nil || interval <= 0 || interval > 400 {
            return "", fmt.Errorf("invalid d interval")
        }
        for {
            date = date.AddDate(0, 0, interval) // сдвигаем на N дней
            if date.After(now) {
                return date.Format(DateFormat), nil
            }
        }

    case "y":
        // Правило "y" — ежегодно
        for {
            date = date.AddDate(1, 0, 0) // сдвигаем на 1 год
            if date.After(now) {
                return date.Format(DateFormat), nil
            }
        }

    case "w":
        // Правило "w - список дней недели"
        // 1 - понедельник ... 7 - воскресенье
        if len(parts) < 2 {
            return "", fmt.Errorf("invalid w rule")
        }
        days := strings.Split(parts[1], ",")
        validDays := make(map[int]bool)
        for _, d := range days {
            n, err := strconv.Atoi(d)
            if err != nil || n < 1 || n > 7 {
                return "", fmt.Errorf("invalid w value: %s", d)
            }
            validDays[n] = true
        }
        for {
            date = date.AddDate(0, 0, 1) // двигаем по одному дню
            weekday := int(date.Weekday())
            if weekday == 0 {
                weekday = 7 // воскресенье - 7
            }
            if validDays[weekday] && date.After(now) {
                return date.Format(DateFormat), nil
            }
        }

    case "m":
        // Правило "m - дни месяца > [месяцы]"
        if len(parts) < 2 {
            return "", fmt.Errorf("invalid m rule")
        }
        days := strings.Split(parts[1], ",")
        validDays := make(map[int]bool)
		negCount := 0
		hasMinusOne := false
		for _, d := range days {
			n, err := strconv.Atoi(d)
			if err != nil || (n == 0 || n < -31 || n > 31) {
				return "", fmt.Errorf("invalid m day: %s", d)
			}
			if n < 0 {
                negCount++
				if n == -1 {
                    hasMinusOne = true
				}
			}
			validDays[n] = true
		}
		// если больше одного отрицательного и среди них нет -1, то считаем правило некорректным
		if negCount > 1 && !hasMinusOne {
			return "", nil
		}
		
        validMonths := make(map[int]bool)
        if len(parts) > 2 {
            months := strings.Split(parts[2], ",")
            for _, m := range months {
                n, err := strconv.Atoi(m)
                if err != nil || n < 1 || n > 12 {
                    return "", fmt.Errorf("invalid m month: %s", m)
                }
                validMonths[n] = true
            }
        }

    for {
        date = date.AddDate(0, 0, 1) // двигаем по одному дню
        day := date.Day()
        lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
            
		// проверка обычных и отрицательных дней
		if validDays[day] && isValidMonthAndAfter(date, validMonths, now) {
            return date.Format(DateFormat), nil
		}

		// проверка отрицательных индексов (-1 = последний день, -2 = предпоследний и т.д.)
		negIndex := day - lastDay - 1
		if validDays[negIndex] && isValidMonthAndAfter(date, validMonths, now) {
            return date.Format(DateFormat), nil
        }
	}
	
	default:
        return "", fmt.Errorf("unknown repeat rule: %s", parts[0])
    }
}