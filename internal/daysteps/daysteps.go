package daysteps

import (
    "fmt"
    "strconv"
    "strings"
    "time"
)

const (
    stepLength = 0.65 // Длина шага в метрах
    mInKm     = 1000  // Метров в одном километре
)

// Функция parsePackage парсит строку с данными о шагах и длительности прогулки
func parsePackage(data string) (int, time.Duration, error) {
    parts := strings.Split(strings.TrimSpace(data), ",")

    if len(parts) != 2 {
        return 0, 0, fmt.Errorf("неверный формат данных")
    }

    stepsStr := parts[0]
    durationStr := parts[1]

    steps, err := strconv.Atoi(stepsStr)
    if err != nil || steps <= 0 {
        return 0, 0, fmt.Errorf("ошибка преобразования количества шагов")
    }

    duration, err := time.ParseDuration(durationStr)
    if err != nil {
        return 0, 0, fmt.Errorf("невозможно распарсить длительность")
    }

    if duration <= 0 { // Добавляем проверку на недопустимую продолжительность
        return 0, 0, fmt.Errorf("недопустимая продолжительность прогулки")
    }

    return steps, duration, nil
}

// Функция DayActionInfo возвращает информацию о действии в течение дня
func DayActionInfo(data string, weight, height float64) string {
    steps, duration, err := parsePackage(data)
    if err != nil {
        return fmt.Sprintf("Ошибка: %v", err)
    }

    if steps <= 0 {
        return "" // Возвращаем пустую строку, если количество шагов меньше нуля
    }

    // Расстояние в метрах
    distanceMeters := float64(steps) * stepLength

    // Переводим дистанцию в километры
    distanceKilometers := distanceMeters / mInKm

    // Формула для расчета калорий:
    spentCalories := (weight * distanceKilometers * 0.01) + (float64(steps) * 0.05)

    // Формируем итоговую строку
    result := fmt.Sprintf(
        "Количество шагов: %d.\n"+
            "Дистанция составила %.2f км.\n"+
            "Выпали %.2f ккал.",
        steps,
        distanceKilometers,
        spentCalories,
    )

    return result
}