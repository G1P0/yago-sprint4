package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage возвращает количество шагов, продолжительность прогулки и ошибку в случае её возникновения
func parsePackage(data string) (int, time.Duration, error) {

	if len(data) == 0 {
		return 0, 0, errors.New("отсутствуют данные о прогулке")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("некорректные данные о прогулке")
	}

	stepsTmp := parts[0]
	durationTmp := parts[1]

	steps, err := strconv.Atoi(stepsTmp)
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("некорректные данные о количестве пройденных шагов")
	}

	duration, err := time.ParseDuration(durationTmp)
	if err != nil || duration <= 0 {
		return 0, 0, errors.New("некорректные данные о продолжительности прогулки")
	}

	return steps, duration, nil

}

// DayActionInfo возвращает количество шагов, дистанцию в километрах и количество потраченных калорий
func DayActionInfo(data string, weight, height float64) string {

	if len(data) == 0 {
		log.Println("отсутствуют данные о прогулке")
		return ""
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if weight <= 0 || height <= 0 {
		log.Println("вес и рост должны быть больше 0")
		return ""
	}

	path := (float64(steps) * stepLength) / mInKm
	cals, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, path, cals)

}
