package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Константы
const (
	stepLength = 0.65 // Средняя длина шага в метрах
	mInKm      = 1000 // Метров в километре
)

// parsePackage парсит строку с данными о прогулке
func parsePackage(data string) (steps int, duration time.Duration, err error) {
	parts := strings.Split(strings.TrimSpace(data), ",")

	if len(parts) != 2 {
		err = errors.New("неверный формат данных")
		return
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	steps, err = strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		err = errors.New("ошибка преобразования количества шагов")
		return
	}

	duration, err = time.ParseDuration(durationStr)
	if err != nil {
		err = errors.New("невозможно распарсить длительность")
		return
	}

	if duration <= 0 {
		err = errors.New("недопустимая продолжительность прогулки")
		return
	}

	return
}

// DayActionInfo выводит информацию о прогрессе пользователя
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)

	return result
}
