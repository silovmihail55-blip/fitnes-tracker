package spentcalories

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Константы
const (
	stepLengthCoefficient   = 0.45 // Коэффициент для расчёта длины шага от роста
	mInKm                 = 1000  // Метров в километре
	minInH                = 60    // Минут в часе
	walkingCaloriesCoefficient = 0.8 // Корректирующий коэффициент для ходьбы
)

// parseTraining парсит строку с данными о тренировке
func parseTraining(data string) (steps int, activityType string, duration time.Duration, err error) {
	parts := strings.Split(strings.TrimSpace(data), ",")

	if len(parts) != 3 {
		err = errors.New("неверный формат данных")
		return
	}

	stepsStr := parts[0]
	activityType = parts[1]
	durationStr := parts[2]

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
		err = errors.New("недопустимая продолжительность тренировки")
		return
	}

	return
}

// distance вычисляет пройденное расстояние в километрах
func distance(steps int, height float64) float64 {
	stepLengthWithHeight := height * stepLengthCoefficient
	totalDistance := float64(steps) * stepLengthWithHeight
	return totalDistance / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceKM := distance(steps, height)
	hours := duration.Hours()

	if hours == 0 {
		return 0
	}

	speed := distanceKM / hours
	return math.Round(speed*100) / 100
}

// RunningSpentCalories рассчитывает калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные данные")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH
	return math.Round(calories*100) / 100, nil
}

// WalkingSpentCalories рассчитывает калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные данные")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	baseCalories := (weight * speed * durationInMinutes) / minInH
	calories := baseCalories * walkingCaloriesCoefficient

	return math.Round(calories*100) / 100, nil
}

// TrainingInfo формирует подробный отчёт о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	distanceKM := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64
	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	report := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activityType,
		duration.Hours(),
		distanceKM,
		speed,
		calories,
	)

	return report, nil
}
