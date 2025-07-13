package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining возвращает количество шагов, вид активности, её продолжительность и ошибку в случае её возникновения
func parseTraining(data string) (int, string, time.Duration, error) {

	if len(data) == 0 {
		return 0, "", 0, errors.New("отсутствуют данные о прогулке")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("некорректные данные о прогулке")
	}

	stepsTmp := parts[0]
	steps, err := strconv.Atoi(stepsTmp)
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("некорректные данные о количестве пройденных шагов")
	}

	activity := parts[1]

	durationTmp := parts[2]
	duration, err := time.ParseDuration(durationTmp)
	if err != nil || duration <= 0 {
		return 0, "", 0, errors.New("некорректные данные о продолжительности прогулки")
	}

	return steps, activity, duration, nil

}

// distance возвращает дистанцию в километрах
func distance(steps int, height float64) float64 {

	lengthPace := height * stepLengthCoefficient
	distance := (lengthPace * float64(steps)) / mInKm
	return distance
}

// meanSpeed возвращает среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}
	pathLength := distance(steps, height)
	meanSpeed := pathLength / duration.Hours()
	return meanSpeed
}

// TrainingInfo возвращает строку с информацией о тренировке.
//
// Пример возвращаемой строки:
//
//	Тип тренировки: Бег
//	Длительность: 0.75 ч.
//	Дистанция: 10.00 км.
//	Скорость: 13.34 км/ч
//	Сожгли калорий: 18621.75
func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var calories float64

	switch activity {
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
		return "", errors.New("неизвестный тип тренировки")
	}

	distanceKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		distanceKm,
		speed,
		calories,
	)

	return result, nil

}

// RunningSpentCalories возвращает количество калорий, потраченных при беге, и ошибку в случае её возникновения
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные данные")
	}
	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * speed * minutes) / minInH
	return calories, nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе, и ошибку в случае её возникновения
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные данные")
	}
	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * speed * minutes) / minInH
	calories *= walkingCaloriesCoefficient
	return calories, nil

}
