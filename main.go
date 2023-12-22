package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jtway/bowl-scoreboard-exporter/pkg/bowl_exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := bowl_exporter.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("unable to read in config: %w", err))
	}
	promAddress := ":" + strconv.Itoa(config.Prom.Port)

	bowlExporterClient := bowl_exporter.NewBowlScoreboardClient(config)

	bowlExporterClient.Run()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(promAddress, nil)
}
