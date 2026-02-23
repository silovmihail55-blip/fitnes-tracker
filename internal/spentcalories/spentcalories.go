package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	mInKm                      = 1000 // м в км
	minInH                     = 60   // мин в часе
	stepLengthCoefficient      = 0.45 // коэф. длины шага от роста
	walkingCaloriesCoefficient = 0.5  // коэф. калорий при ходьбе
)

// parseTraining — парсит данные тренировки: шаги, тип, длительность
func parseTraining(data string) (int, string, time.Duration, error) {
	spl := strings.Split(data, ",")

	if len(spl) != 3 {
		err := errors.New("Неверный формат данных")
		log.Println(err)
		return 0, "", 0, err
	}
	steps, err := strconv.Atoi(spl[0])
	if err != nil {
		log.Println(err)
		return 0, "", 0, err
	}

	if steps <= 0 {
		err := errors.New("Количество шагов должна быть больше нуля")
		log.Println(err)
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(spl[2])
	if err != nil {
		log.Println(err)
		return 0, "", 0, err
	}

	if duration.Minutes() <= 0 {
		err := errors.New("Длительность тренировки должна быть больше нуля")
		log.Println(err)
		return 0, "", 0, err
	}

	return steps, spl[1], duration, nil
}

// distance — рассчитывает дистанцию в км по шагам и росту
func distance(steps int, height float64) float64 {
	steplen := height * stepLengthCoefficient
	return float64(steps) * steplen / mInKm
}

// meanSpeed — рассчитывает среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Hours() <= 0 {
		return 0
	}

	if steps <= 0 {
		return 0
	}

	if height <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

// TrainingInfo — формирует отчёт по тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, tp, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if weight <= 0 {
		err = errors.New("Вес должен быть больше нуля")
		log.Println(err)
		return "", err
	}

	if height <= 0 {
		err = errors.New("Рост должен быть больше нуля")
		log.Println(err)
		return "", err
	}

	var calories float64

	switch tp {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		err := fmt.Errorf("неизвестный тип тренировки: %s", tp)
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", tp, duration.Hours(), dist, speed, calories), nil
}

// RunningSpentCalories — рассчитывает калории для бега
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		err := errors.New("Вес должен быть больше нуля")
		log.Println(err)
		return 0, err
	}

	if height <= 0 {
		err := errors.New("Рост должен быть больше нуля")
		log.Println(err)
		return 0, err
	}

	if steps <= 0 {
		err := errors.New("Количество шагов должна быть больше нуля")
		log.Println(err)
		return 0, err
	}

	if duration.Minutes() <= 0 {
		err := errors.New("Длительность тренировки должна быть больше нуля")
		log.Println(err)
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	calories := (weight * speed * duration.Minutes()) / minInH

	return calories, nil
}

// WalkingSpentCalories — рассчитывает калории для ходьбы
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		err := errors.New("Вес должен быть больше нуля")
		log.Println(err)
		return 0, err
	}

	if height <= 0 {
		err := errors.New("Рост должен быть больше нуля")
		log.Println(err)
		return 0, err
	}

	if steps <= 0 {
		err := errors.New("Количество шагов не может быть отрицательным")
		log.Println(err)
		return 0, err
	}

	if duration.Minutes() <= 0 {
		err := errors.New("Длительность тренировки должна быть больше нуля")
		log.Println(err)
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	calories := walkingCaloriesCoefficient * (weight * speed * duration.Minutes()) / minInH

	return calories, nil
}
