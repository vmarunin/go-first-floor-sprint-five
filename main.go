package main

import (
	"fmt"
	"math"
	"time"
)

const (
	MInKm      = 1000
	MinInHours = 60
	LenStep    = 0.65
)

type Training struct {
	TrainingType string
	Action       int
	LenStep      float64
	Duration     time.Duration
	Weight       float64
}

func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

func (t Training) meanSpeed() float64 {
	return t.distance() / t.Duration.Hours()
}

func (t Training) Calories() float64 {
	return 0 // Базовая реализация, будет переопределена в подклассах
}

type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

const (
	CaloriesMeanSpeedMultiplier = 18
	CaloriesMeanSpeedShift      = 1.79
)

type Running struct {
	Training
}

func (r Running) Calories() float64 {
	return ((CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Hours() * MinInHours)
}

const (
	CaloriesWeightMultiplier      = 0.035
	CaloriesSpeedHeightMultiplier = 0.029
	KmHInMsec                     = 0.278
)

type Walking struct {
	Training
	Height float64
}

func (w Walking) Calories() float64 {
	speedInMPerS := w.meanSpeed() * KmHInMsec
	return ((CaloriesWeightMultiplier*w.Weight + (math.Pow(speedInMPerS, 2)/(w.Height/100))*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours)
}

const (
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2
)

type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

func (s Swimming) meanSpeed() float64 {
	return float64(s.LengthPool*s.CountPool) / MInKm / s.Duration.Hours()
}

func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	updateCalories := training.Calories()
	info.Calories = updateCalories
	return info.String()
}

func main() {
	swimming := Swimming{
		Training:   Training{"Плавание", 2000, SwimmingLenStep, 90 * time.Minute, 85},
		LengthPool: 50,
		CountPool:  40,
	}
	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{"Ходьба", 20000, LenStep, 3*time.Hour + 45*time.Minute, 85},
		Height:   185,
	}
	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{"Бег", 5000, LenStep, 30 * time.Minute, 85},
	}
	fmt.Println(ReadData(running))
}
