package main

import (
	"context"
	"fmt"
	"go-jeager-mcsrvs-example/worker/calculator"
	"go-jeager-mcsrvs-example/worker/config"
	"go-jeager-mcsrvs-example/worker/controller"
	"go-jeager-mcsrvs-example/worker/tracer"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	ctx := context.Background()
	config := config.Config{
		Port:        ":8082",
		JeagerURL:   "http://localhost:14268/api/traces",
		ServiceName: "worker",
	}

	tp, err := tracer.InitTracer(config.JeagerURL, config.ServiceName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tracer.Stop(ctx, tp)

	calcEngine := calculator.NewEngine()

	controller := controller.New(config, calcEngine)

	router := gin.Default()
	router.Use(otelgin.Middleware("worker-http-server")) // внутри спана помечается как net.host.name
	router.POST("/summ", controller.SummHandler)
	router.POST("/multi", controller.MultiHandler)

	server := &http.Server{
		Addr:    config.Port,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
