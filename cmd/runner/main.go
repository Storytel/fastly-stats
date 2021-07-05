package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	fastlystats "github.com/Storytel/fastly-stackdriver-exporter"
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

func multiplexChannel(ctx context.Context, ch <-chan *fastlystats.FastlyMeanStats, consumers []chan *fastlystats.FastlyMeanStats) {
	for {
		select {
		case <-ctx.Done():
			return
		case stats := <-ch:
			for _, c := range consumers {
				c <- stats
			}
		}
	}
}

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
		fastlystats.SetupMetricDescriptors(ctx, googleCloudProject)
		return
	}

	cfg := &fastlystats.Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		ll.Fatal(err)
	}

	if cfg.FastlyAPIKey == "" {
		ll.Fatal("Fastly API key missing, set env FASTLY_API_KEY")
	}

	if cfg.FastlyService == "" {
		ll.Fatal("Fastly Service is missing, set env FASTLY_SERVICE")
	}

	if cfg.NewRelicInsertKey == "" {
		ll.Fatal("New Relic insert key is missing, set env NEWRELIC_INSERT_KEY")
	}

	ch := make(chan *fastlystats.FastlyMeanStats)

	provider, err := fastlystats.NewFastlyStatsProvider(cfg.FastlyService, cfg.FastlyAPIKey, ch)
	if err != nil {
		ll.Fatal(err)
	}

	consumers := []chan *fastlystats.FastlyMeanStats{
		make(chan *fastlystats.FastlyMeanStats, 1024),
		make(chan *fastlystats.FastlyMeanStats, 1024),
	}

	go multiplexChannel(ctx, ch, consumers)

	wg := sync.WaitGroup{}
	wg.Add(len(consumers) + 1)

	go func() {
		defer wg.Done()
		defer cancel()
		provider.Run(ctx)
	}()

	go func() {
		defer wg.Done()
		// defer cancel()
		// consumer, err := fastlystats.NewStackdriverExporter(googleCloudProject, consumers[0])
		// if err != nil {
		// 	ll.Fatal(err)
		// }
		// consumer.Run(ctx)
	}()

	go func() {
		defer wg.Done()
		defer cancel()
		consumer, err := fastlystats.NewNewRelicExporter(cfg.NewRelicInsertKey, consumers[1])
		if err != nil {
			ll.Fatal(err)
		}
		consumer.Run(ctx)
	}()

	wg.Wait()
}
