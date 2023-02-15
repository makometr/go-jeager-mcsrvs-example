package calculator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrZeroValueFound error = errors.New("zero value found")
)

type ErrArraySizeInvalid struct {
	Size int
}

func (err ErrArraySizeInvalid) Error() string {
	return fmt.Sprintf("size = %d > 5", err.Size)
}

type CalcEngine struct {
	tracer trace.Tracer
}

func NewEngine() *CalcEngine {
	return &CalcEngine{
		tracer: otel.Tracer("calculator-engine"), // это просто отображается в спане как otel.library.name
	}
}

// SummIntegers генерирует ошибки, если есть нули или длина больше 5.
func (engine *CalcEngine) SummIntegers(ctx context.Context, data []int) (int, error) {
	ctx, span := engine.tracer.Start(ctx, "engine SummIntegers")
	defer span.End()

	if len(data) > 5 {
		logrus.WithContext(ctx).Error(ErrArraySizeInvalid{Size: len(data)})
		return 0, &ErrArraySizeInvalid{Size: len(data)}
	}

	var result int
	for _, n := range data {
		if err := engine.intergerHandler(ctx, n); err != nil {
			return 0, fmt.Errorf("engine error: %w", err)
		}
		result += n
	}

	return result, nil
}

func (engine *CalcEngine) intergerHandler(ctx context.Context, number int) error {
	ctx, span := engine.tracer.Start(ctx, "intergerHandler")
	defer span.End()
	span.SetAttributes(attribute.Int("number", number))

	if number == 0 {
		logrus.WithContext(ctx).Error(ErrZeroValueFound)
		return ErrZeroValueFound
	}
	if number == 1 {
		time.Sleep(3 * time.Second)
	} else {
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}
