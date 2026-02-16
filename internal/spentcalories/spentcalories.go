package spentcalories

import (
    "errors"
    "fmt"
    "log"
    "math"
    "strconv"
    "strings"
    "time"
)

// Константы
const (
    stepLength = 0.65 // Средняя длина шага в метрах
    mInKm      = 1000 // Кол-во метров в километра
    kWalk      = 0.03 // Коэффициент для ходьбы
    kRun       = 0.06 // Коэффициент для бега
)

// parseTraining парсит строку с данными о тренировке
func parseTraining(data string) (steps int, activityType string, duration time.Duration, err error) {
    parts := strings.Split(strings.TrimSpace(data), ",")

    if len(parts) != 3 {
        err = fmt.Errorf("неверный формат данных")
        return
    }

    stepsStr := parts[0]
    activityType = parts[1]
    durationStr := parts[2]

    steps, err = strconv.Atoi(stepsStr)
    if err != nil || steps <= 0 {
        err = fmt.Errorf("ошибка преобразования количества шагов")
        return
    }

    duration, err = time.ParseDuration(durationStr)
    if err != nil {
        err = fmt.Errorf("невозможно распарсить длительность")
        return
    }

    if duration <= 0 {
        err = fmt.Errorf("недопустимая продолжительность тренировки")
        return
    }

    return
}

// distance вычисляет пройденное расстояние в километрах
func distance(steps int, height float64) float64 {
    stepLengthWithHeight := height * 0.45
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
    return math.Round(distanceKM*100/hours) / 100
}

// spentCalories вычисляет общие затраты калорий
func spentCalories(activityType string, steps int, weight, height float64, duration time.Duration) (float64, error) {
    if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
        return 0, fmt.Errorf("некорректные входные данные")
    }

    distanceKM := distance(steps, height)
    var coefficient float64

    switch activityType {
    case "Бег":
        coefficient = kRun
    case "Ходьба":
        coefficient = kWalk
    default:
        return 0, fmt.Errorf("неизвестный тип тренировки")
    }

    // Основной расчет калорий
    calories := (weight * distanceKM * coefficient) + (float64(steps) * 0.05)

    return calories, nil
}

// TrainingInfo формирует отчет о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
    steps, activityType, duration, err := parseTraining(data)
    if err != nil {
        log.Printf("Ошибка парсинга данных: %v", err)
        return "", err
    }

    if steps <= 0 || duration <= 0 || weight <= 0 || height <= 0 {
        return "", fmt.Errorf("некорректные входные данные")
    }

    distanceKM := distance(steps, height)
    speed := meanSpeed(steps, height, duration)

    calories, err := spentCalories(activityType, steps, weight, height, duration)
    if err != nil {
        log.Printf("Ошибка расчета калорий: %v", err)
        return "", err
    }

    report := fmt.Sprintf(
        "Тип тренировки: %s\n"+
            "Длительность: %.2f ч.\n"+
            "Дистанция: %.2f км.\n"+
            "Скорость: %.2f км/ч\n"+
            "Сожгли калорий: %.2f",
        activityType,
        duration.Hours(),
        distanceKM,
        speed,
        calories,
    )

    return report, nil
}