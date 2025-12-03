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

func parseTraining(data string) (int, string, time.Duration, error) {
	// разделяем строку по запятой
	parts := strings.Split(data, ",")
	// проверка длины строки
	if len(parts) != 3 {
		return 0, "", 0, errors.New("string does not contain exactly 3 elements")
	}

	// преобразование первого элемента слайса в int
	steps, err := strconv.Atoi(parts[0]) // преобразование строки в int
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("number of steps must be a positive number")
	}

	// преобразование третьего элемента слайса в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("time must be a positive number")
	}

	activityType := parts[1]

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// длина шага
	lengthStep := height * stepLengthCoefficient
	// пройденный путь
	pathTraveled := float64(steps) * lengthStep
	// дистанция в километрах
	return pathTraveled / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// проверка продолжительности прогулки
	if duration <= 0 {
		return 0
	}

	// вычисление дистанции (шаг и рост пользователя)
	distanceKm := distance(steps, height)

	// вычисление средней скорости
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return distanceKm / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// проверка входных параметров на корректность
	if steps <= 0 {
		return 0, errors.New("steps must be greater than 0")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, errors.New("height must be greater than 0")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be greater than 0")
	}

	// расчет средней скорости
	speedAverage := meanSpeed(steps, height, duration)

	// расчет калорий
	minutes := duration.Minutes()
	calories := (weight * speedAverage * minutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// проверка входных параметров на корректность
	if steps <= 0 {
		return 0, errors.New("steps must be greater than 0")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, errors.New("height must be greater than 0")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be greater than 0")
	}

	// расчет средней скорости
	speedAverage := meanSpeed(steps, height, duration)

	// расчет калорий
	minutes := duration.Minutes()
	calories := (weight * speedAverage * minutes) / minInH

	// калории для ходьбы
	coefficientActivity := calories * walkingCaloriesCoefficient
	return coefficientActivity, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if steps <= 0 {
		return "", err
	}

	distans := distance(steps, height)
	speedAverage := meanSpeed(steps, height, duration)
	var calories float64

	// подсчет параметров для каждого вида тренирвоки
	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, float64(duration.Hours()), distans, speedAverage, calories)
	return result, nil
}
