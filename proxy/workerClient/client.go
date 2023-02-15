package workerClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type WorkerClient struct {
	URL    string
	tracer trace.Tracer
}

func NewWorkerClient(host string, port int) *WorkerClient {
	return &WorkerClient{
		URL:    fmt.Sprintf("http://%s:%d", host, port),
		tracer: otel.Tracer("worker-client"), // это просто отображается в спане как otel.library.name
	}
}

func (client *WorkerClient) SummIntegers(ctx context.Context, data []int) (int, error) {
	ctx, span := client.tracer.Start(ctx, "WorkerClient SummIntegers")
	defer span.End()

	body := SummRequest{
		Numbers: data,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return 0, err
	}
	bodyReader := bytes.NewReader(bodyBytes)

	req, err := http.NewRequest(http.MethodPost, client.URL+"/summ", bodyReader)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("client: error making http request")
		return 0, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("client: parse result error")
	}
	fmt.Printf("client: response body: %s\n", string(resBody))

	return 10, nil
}
