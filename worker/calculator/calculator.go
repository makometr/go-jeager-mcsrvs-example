package calculator

import (
	"context"
	"fmt"
	"time"
)

type CalcEngine struct {
}

func NewEngine() *CalcEngine {
	return &CalcEngine{}
}

// SummIntegers генерирует ошибки.
func (ce *CalcEngine) SummIntegers(ctx context.Context, data []int) (int, error) {
	var result int
	for _, n := range data {
		if n == 0 {
			return 0, fmt.Errorf("zero value found")
		}
		result += n
	}

	return result, nil
}

// MultIntegers спит.
func (ce *CalcEngine) MultIntegers(ctx context.Context, data []int) (int, error) {
	result := 1
	for _, n := range data {
		if n == 1 {
			time.Sleep(time.Second * 3)
		}
		result *= n
	}

	return result, nil
}
