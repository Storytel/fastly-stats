package fastlystackdriver

import (
	"context"
	"errors"
	"math"
	"reflect"
	"time"

	"github.com/fastly/go-fastly/v3/fastly"
	"go.uber.org/zap"
)

// Must not be shorter than the quota for creating time series.
// See https://cloud.google.com/monitoring/quotas
// At the time of writing this is 1 point per 10 seconds, so must be more than 10 seconds
//
// The shorter this is, the higher resolution there will be on datapoints, but
// also increases resource usage and billing
const pollInterval = 15 * time.Second

type FastlyMeanStats struct {
	IntervalStart uint64
	IntervalEnd   uint64
	Stats         *fastly.Stats
}

type FastlyStatsProvider struct {
	fastlyClient *fastly.RTSClient
	service      string
	ch           chan<- *FastlyMeanStats

	timestamp uint64
}

func NewFastlyStatsProvider(service, apiKey string, ch chan<- *FastlyMeanStats) (*FastlyStatsProvider, error) {
	fastlyClient, err := fastly.NewRealtimeStatsClientForEndpoint(apiKey, fastly.DefaultRealtimeStatsEndpoint)
	if err != nil {
		return nil, err
	}

	return &FastlyStatsProvider{
		fastlyClient: fastlyClient,
		service:      service,
		ch:           ch,
	}, nil
}

func (f *FastlyStatsProvider) Run(ctx context.Context) {
	ll := zap.S()
	ll.Infof("starting fastly stats provider")
	for {
		start := time.Now()
		if err := f.next(ctx, f.ch); err != nil && !errors.Is(err, context.Canceled) {
			ll.Warnf("failed to get stats, retrying in a while: %v", err)
		}
		dur := time.Since(start)

		ll.Debugf("getting and reporting stats took %v - sleeping for %v", dur, pollInterval-dur)

		select {
		case <-time.After(pollInterval - dur):
		case <-ctx.Done():
			return
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func (f *FastlyStatsProvider) mean(list []*fastly.RealtimeData) *FastlyMeanStats {
	stats := &fastly.Stats{}
	n := uint64(len(list))

	var min uint64 = math.MaxUint64
	var max uint64 = 0

	refStats := reflect.ValueOf(stats)

	for _, rtdata := range list {
		vs := reflect.ValueOf(rtdata.Aggregated)
		for i := 0; i < vs.Elem().NumField(); i++ {
			sf := vs.Elem().Field(i)
			df := refStats.Elem().Field(i)

			switch sf.Kind() {
			case reflect.Uint64:
				df.SetUint(sf.Uint() + df.Uint())
			case reflect.Float64:
				df.SetFloat(sf.Float() + df.Float())
			}
		}

		if rtdata.Recorded < min {
			min = rtdata.Recorded
		}

		if rtdata.Recorded > max {
			max = rtdata.Recorded
		}
	}
	for i := 0; i < refStats.Elem().NumField(); i++ {
		f := refStats.Elem().Field(i)
		switch f.Kind() {
		case reflect.Uint64:
			f.SetUint(uint64(math.Round(float64(f.Uint()) / float64(n))))
		case reflect.Float64:
			f.SetFloat(f.Float() / float64(n))
		}
	}

	return &FastlyMeanStats{
		IntervalStart: min,
		IntervalEnd:   max,
		Stats:         stats,
	}
}

func (s *FastlyStatsProvider) next(ctx context.Context, ch chan<- *FastlyMeanStats) error {
	ll := zap.S()
	req := &fastly.GetRealtimeStatsInput{
		ServiceID: s.service,
		Timestamp: s.timestamp,
	}

	ll.Debugf("getting realtime stats")
	resp, err := s.fastlyClient.GetRealtimeStats(req)
	if err != nil {
		return err
	}

	ll.Debugf("got %d seconds worth of value", len(resp.Data))

	meanStats := s.mean(resp.Data)
	s.timestamp = resp.Timestamp

	select {
	case ch <- meanStats:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
