package main

import (
	"context"
	"go-jeager-mcsrvs-example/proxy/config"
	"go-jeager-mcsrvs-example/proxy/controller"
	"go-jeager-mcsrvs-example/proxy/tracer"
	"go-jeager-mcsrvs-example/proxy/workerClient"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)

	// Instrument logrus.
	logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	)))
	// https://github.com/uptrace/opentelemetry-go-extra/tree/main/otellogrus

	ctx := context.Background()
	config := config.Config{
		Port:        ":8081",
		ServiceName: "proxy-service",

		WorkerHost: "localhost",
		WorkerPort: 8082,

		JeagerURL: "http://localhost:14268/api/traces",
	}

	// Jeager UI: http://localhost:16686
	tp, err := tracer.InitTracer(config.JeagerURL, config.ServiceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer tracer.Stop(ctx, tp)

	workerClient := workerClient.NewWorkerClient(config.WorkerHost, config.WorkerPort)

	controller := controller.New(config, workerClient)

	router := gin.Default()
	router.Use(otelgin.Middleware("proxy-server")) // внутри спана помечается как net.host.name
	router.POST("/summ", controller.ProxySummHandler)

	server := &http.Server{
		Addr:    config.Port,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		logrus.Fatalln(err)
	}
}
