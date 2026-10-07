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

// Виды тренировок.
const (
	trainingWalking = "Ходьба"
	trainingRunning = "Бег"
)

// parseTraining разбирает строку формата "3456,Ходьба,3h00m".
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных: ожидается \"шаги,вид активности,продолжительность\"")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, parts[1], duration, nil
}

// distance возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

// meanSpeed возвращает среднюю скорость в км/ч.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// TrainingInfo возвращает информацию о тренировке.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	switch activity {
	case trainingRunning:
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case trainingWalking:
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), dist, speed, calories), nil
}

// validateParams проверяет входные параметры для расчёта калорий.
func validateParams(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return errors.New("количество шагов должно быть больше 0")
	case weight <= 0:
		return errors.New("вес должен быть больше 0")
	case height <= 0:
		return errors.New("рост должен быть больше 0")
	case duration <= 0:
		return errors.New("продолжительность должна быть больше 0")
	}
	return nil
}

// RunningSpentCalories возвращает количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	return weight * speed * duration.Minutes() / minInH, nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	calories := weight * speed * duration.Minutes() / minInH
	return calories * walkingCaloriesCoefficient, nil
}
