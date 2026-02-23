package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	sc "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65 // длина шага (м)
	mInKm      = 1000 // м в км
)

// parsePackage — парсит строку на шаги и длительность
func parsePackage(data string) (int, time.Duration, error) {
	spl := strings.Split(data, ",")
	var step int
	var err error
	var duration time.Duration

	if len(spl) == 2 {
		// Парсим шаги
		step, err = strconv.Atoi(spl[0])
		if err != nil {
			log.Println(err)
			return 0, 0, err
		}
		if step <= 0 {
			err := errors.New("Количество шагов должна быть больше нуля")
			log.Println(err)
			return 0, 0, err
		}

		// Парсим длительность
		duration, err = time.ParseDuration(spl[1])
		if err != nil {
			log.Println(err)
			return 0, 0, err
		}
		if duration <= 0 {
			err := errors.New("Длительность должна быть больше нуля")
			log.Println(err)
			return 0, 0, err
		}

		return step, duration, nil
	}

	err = errors.New("Неверный формат данных")
	log.Println(err)
	return 0, 0, err
}

// DayActionInfo — формирует отчёт о активности
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		fmt.Println("Ошибка DayActionInfo: ", err)
		return ""
	}

	// Расчёт дистанции в км
	distance := float64(steps) * stepLength / mInKm

	// Расчёт сожжённых калорий
	wspent, err := sc.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		fmt.Println("Ошибка: ", err)
		return ""
	}

	// Форматирование итогового сообщения
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, wspent)
}
