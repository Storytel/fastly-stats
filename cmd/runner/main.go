package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	fastlystackdriver "github.com/Storytel/fastly-stackdriver-exporter"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

var rebuildMetricDescriptors bool
var outputJson bool
var googleCloudProject string

func main() {
	flag.BoolVar(&outputJson, "output-json", false, "Whether output should be JSON encoded")
	flag.BoolVar(&rebuildMetricDescriptors, "rebuild-metric-descriptors", false, "Re-build all metric descriptors and exit")
	flag.StringVar(&googleCloudProject, "project", "", "The Google Cloud Project to delete metrics from")
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	if outputJson {
		logger, _ = zap.NewProduction(zap.IncreaseLevel(zap.InfoLevel))
	}
	zap.ReplaceGlobals(logger)

	ll := logger.Sugar()

	if googleCloudProject == "" {
		ll.Fatal("Specify Google Cloud Project with the -project flag.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		closeSignal := <-ch

		ll.Infof("Received close signal (%v), terminating gracefully.", closeSignal)
		cancel()

		closeSignal = <-ch
		ll.Infof("Received second close (%v) signal, exiting forcefully (2).", closeSignal)
		os.Exit(2)
	}()

	if rebuildMetricDescriptors {
		fastlystackdriver.SetupMetricDescriptors(ctx, googleCloudProject)
		return
	}

	cfg := &fastlystackdriver.Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		ll.Fatal(err)
	}

	if cfg.FastlyAPIKey == "" {
		ll.Fatal("Fastly API key missing, set env FASTLY_API_KEY")
	}

	if cfg.FastlyService == "" {
		ll.Fatal("Fastly Service is missing, set env FASTLY_SERVICE")
	}

	ch := make(chan *fastlystackdriver.FastlyMeanStats)

	provider, err := fastlystackdriver.NewFastlyStatsProvider(cfg.FastlyService, cfg.FastlyAPIKey, ch)
	if err != nil {
		ll.Fatal(err)
	}

	consumer, err := fastlystackdriver.NewStackdriverExporter(googleCloudProject, ch)
	if err != nil {
		ll.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer cancel()
		provider.Run(ctx)
	}()

	go func() {
		defer wg.Done()
		defer cancel()
		consumer.Run(ctx)
	}()

	wg.Wait()
}
