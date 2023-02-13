package main

import (
	"context"
	"go-jeager-mcsrvs-example/worker/calculator"
	"go-jeager-mcsrvs-example/worker/config"
	"go-jeager-mcsrvs-example/worker/controller"
	"go-jeager-mcsrvs-example/worker/tracer"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

	// Instrument logrus.
	log.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
	)))

	ctx := context.Background()
	config := config.Config{
		Port:        ":8082",
		JeagerURL:   "http://localhost:14268/api/traces",
		ServiceName: "worker",
	}
	// http://localhost:16686

	tp, err := tracer.InitTracer(config.JeagerURL, config.ServiceName)
	if err != nil {
		log.Fatal(err)
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
		log.Fatalln(err)
	}
}
