package calculator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrZeroValueFound error = errors.New("zero value found")
)

type CalcEngine struct {
	tracer trace.Tracer
}

func NewEngine() *CalcEngine {
	return &CalcEngine{
		tracer: otel.Tracer("calculator-engine"), // это просто отображается в спане как otel.library.name
	}
}

// SummIntegers генерирует ошибки.
func (engine *CalcEngine) SummIntegers(ctx context.Context, data []int) (int, error) {
	ctx, span := engine.tracer.Start(ctx, "engine SummIntegers")
	defer span.End()

	var result int
	for _, n := range data {
		if err := engine.intergerHandler(ctx, n); err != nil {
			return 0, fmt.Errorf("engine error: %w", err)
		}
		result += n
	}

	return result, nil
}

// MultIntegers спит.
func (engine *CalcEngine) MultIntegers(ctx context.Context, data []int) (int, error) {
	ctx, span := engine.tracer.Start(ctx, "engine MultIntegers")
	defer span.End()

	result := 1
	for _, n := range data {
		if err := engine.intergerHandler(ctx, n); err != nil {
			return 0, fmt.Errorf("engine error: %w", err)
		}
		result *= n
	}

	return result, nil
}

func (engine *CalcEngine) intergerHandler(ctx context.Context, number int) error {
	ctx, span := engine.tracer.Start(ctx, "intergerHandler")
	defer span.End()
	span.SetAttributes(attribute.Int("number", number))

	if number == 0 {
		err := ErrZeroValueFound
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	if number == 1 {
		time.Sleep(3 * time.Second)
	} else {
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}
