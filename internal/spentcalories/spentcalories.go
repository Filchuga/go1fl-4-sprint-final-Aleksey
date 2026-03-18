package spentcalories

import (
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
	// Делим строку на слайс строк
	split := strings.Split(data, ",")
	if len(split) == 3 {
		steps, err := strconv.Atoi(split[0]) // Преобразуем число шагов в int
		if err != nil {
			return 0, "", 0, fmt.Errorf("ошибка преобразования количества шагов: %v", err)
		}
		duration, err := time.ParseDuration(split[2]) // Преобразуем продолжительность прогулки в time.Duration
		if err != nil {
			return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности прогулки: %v", err)
		}
		// Проверяем введенные данные
		if steps <= 0 || duration <= 0 {
			return 0, "", 0, fmt.Errorf("количество шагов и продолжительность должно быть больше 0, введено %d, %v", steps, duration)
		}
		return steps, split[1], duration, nil
	} else {
		return 0, "", 0, fmt.Errorf("ошибка: введите данные с тремя значениями, введено %d", len(split))
	}
}

func distance(steps int, height float64) float64 {
	// Расчитываем длину шага
	lengthSteps := float64(height) * stepLengthCoefficient
	distance := float64(steps) * lengthSteps // Расчитываем пройденую дистанцию в метрах
	distanceKm := distance / mInKm           // Переводим метры в километры
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем продолжительность прогулки
	if duration <= 0 {
		return 0.0
	}
	distanceKm := distance(steps, height) // Расчитываем пройденную дистанцию
	// Расчитываем среднюю скорость
	hours := duration.Hours()
	averageSpeed := distanceKm / float64(hours)
	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Делим строку на слайс строк
	steps, tipeTr, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var distanceKm float64
	var averageSpeed float64
	var calories float64
	switch tipeTr {
	case "Ходьба":
		distanceKm = distance(steps, height)
		averageSpeed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Бег":
		distanceKm = distance(steps, height)
		averageSpeed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		errMes := fmt.Errorf("неизвестный тип тренировки %s", tipeTr)
		log.Println(errMes)
		return "", errMes
	}
	resultString := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", tipeTr, duration.Hours(), distanceKm, averageSpeed, calories)
	return resultString, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем входные данные
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0.0, fmt.Errorf("введеные значения должны быть больше 0, у вас:\n%d шаги\n%.2f вес,\n%.2f рост,\n%v продолжительность прогулки", steps, weight, height, duration)
	}
	// Расчитываем количество калорий
	averageSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	caloriesInM := weight * averageSpeed * minutes
	caloriesInH := caloriesInM / minInH
	return caloriesInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем входные данные
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0.0, fmt.Errorf("введеные значения должны быть больше 0, у вас:\n%d шаги\n%.2f вес,\n%.2f рост,\n%v продолжительность прогулки", steps, weight, height, duration)
	}
	// Расчитываем количество калорий
	averageSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	caloriesInM := weight * averageSpeed * minutes
	caloriesInH := caloriesInM / minInH
	caloriesCoeff := caloriesInH * walkingCaloriesCoefficient
	return caloriesCoeff, nil
}
