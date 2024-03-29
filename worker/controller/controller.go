package controller

import (
	"context"
	"go-jeager-mcsrvs-example/worker/config"
	"go-jeager-mcsrvs-example/worker/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Calculator interface {
	SummIntegers(ctx context.Context, data []int) (int, error)
}

type Controller struct {
	calc   Calculator
	tracer trace.Tracer
}

func New(cfg config.Config, calculator Calculator) *Controller {
	return &Controller{
		calc:   calculator,
		tracer: otel.Tracer("http-server-worker"), // это просто отображается в спане как otel.library.name
	}
}

func (server *Controller) SummHandler(c *gin.Context) {
	ctx, span := server.tracer.Start(c.Request.Context(), "server SummHandler")
	defer span.End()

	var data models.Request
	if err := c.BindJSON(&data); err != nil {
		logrus.WithContext(ctx).WithError(err).Error("parse failed")
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	span.SetAttributes(attribute.IntSlice("data", data.Numbers))

	if data.Numbers == nil {
		c.XML(400, gin.H{
			"msg": "nastay, privet",
		})
		return
	}

	db.setUser(id)\
	selct from user where id = $id

	result, err := server.calc.SummIntegers(ctx, data.Numbers)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"result": result,
	})
}
