package metrics

import (
	"context"
	"log"
	"time"

	"github.com/influxdata/telegraf/agent"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/models"
	output_http "github.com/influxdata/telegraf/plugins/outputs/http"
	"github.com/influxdata/telegraf/plugins/serializers/prometheusremotewrite"
)

type PromConfig struct {
	PrometheusRemoteWriteURI      string
	PrometheusRemoteWriteUsername string
	PrometheusRemoteWritePassword string
}

func isConfigValid(config PromConfig) bool {
	return config.PrometheusRemoteWritePassword != "" && config.PrometheusRemoteWriteURI != "" && config.PrometheusRemoteWriteUsername != ""
}

func InitMetrics(ctx context.Context, promConfig PromConfig) *Metrics {
	inputPlugin := GetMetricsHandler()

	if !isConfigValid(promConfig) {
		log.Println("Config is not set for metrics")
		return inputPlugin
	}

	outputPlugin := &output_http.HTTP{
		URL:      promConfig.PrometheusRemoteWriteURI,
		Username: promConfig.PrometheusRemoteWriteUsername,
		Password: promConfig.PrometheusRemoteWritePassword,
	}

	err := outputPlugin.Connect()
	if err != nil {
		log.Fatalf("Can not init metrics: %v", err)
	} else {
		log.Print("Plugin connected")
	}

	s, err := prometheusremotewrite.NewSerializer(prometheusremotewrite.FormatConfig{
		MetricSortOrder: prometheusremotewrite.SortMetrics,
		StringHandling:  prometheusremotewrite.StringAsLabel,
	})
	if err != nil {
		log.Fatalf("Can not init metrics: %v", err)
	}
	outputPlugin.SetSerializer(s)

	config := config.NewConfig()

	duration, err := time.ParseDuration("10s")
	if err != nil {
		log.Fatalf("Can not init metrics: %v", err)
	}

	runningInput := models.NewRunningInput(inputPlugin, &models.InputConfig{Interval: duration})
	runningInput.Init()

	runningOutput := models.NewRunningOutput(outputPlugin, &models.OutputConfig{}, 1000, 10000)
	runningOutput.Init()
	config.InputFilters = append(config.InputFilters, "metrics")
	config.Inputs = append(config.Inputs, runningInput)
	config.Outputs = append(config.Outputs, runningOutput)

	telegrafClient, err := agent.NewAgent(config)
	if err != nil {
		log.Fatalf("Can not init metrics: %v", err)
	}
	go telegrafClient.Run(ctx)
	log.Print("Run telegraf client")

	return inputPlugin
}
