package daysteps

import (
	"errors"
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// Функция принимает строку с данными, которая содержит количество шагов и продолжительность прогулки в формате 3h50m(3 часа 50 минут)
func parsePackage(data string) (int, time.Duration, error) {
	// разделяем строку по запятой
	parts := strings.Split(data, ",")
	// проверка длины строки
	if len(parts) != 2 {
		return 0, 0, errors.New("Ошибка: строка не содержит ровно два элемента")
	}

	// преобразование первого элемента слайса в int
	steps, err := strconv.Atoi(parts[0]) // преобразование строки в int
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, errors.New("Ошибка: количество шагов должно быть больше 0")
	}

	// преобразование второго элемента слайса в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, errors.New("Ошибка: время должно быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка: ", err)
		return ""
	}

	// проверка на количество шагов
	if steps <= 0 {
		log.Println("Ошибка: количество шагов должно быть больше 0")
		return ""
	}

	// вычисление дистанции в метрах
	distanceM := stepLength * float64(steps)
	// дистанция в километрах
	distanceKm := distanceM / mInKm

	// вычисление калорий
	calories, _ := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка расчета калорий:", err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
