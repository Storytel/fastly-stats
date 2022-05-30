package fastlystats

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"go.uber.org/zap"
)

var ErrNotFound = errors.New("not found")

const nrEndpoint = "https://metric-api.eu.newrelic.com/metric/v1"

type NewRelicMetricReport struct {
	Metrics []NewRelicMetricDescriptor `json:"metrics"`
}

type NewRelicExporter struct {
	insertKey  string
	ch         <-chan *FastlyMeanStats
	httpClient *http.Client
}

func NewNewRelicExporter(insertKey string, ch <-chan *FastlyMeanStats) (*NewRelicExporter, error) {
	return &NewRelicExporter{
		insertKey: insertKey,
		ch:        ch,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func getNewRelicMetric(name string) (NewRelicMetricDescriptor, error) {
	for _, md := range NRMetricDescriptors {
		if md.Name == name {
			return md, nil
		}
	}

	return NewRelicMetricDescriptor{}, ErrNotFound
}

func (n *NewRelicExporter) buildMetrics(s *FastlyMeanStats) []NewRelicMetricReport {
	metrics := []NewRelicMetricReport{
		{Metrics: make([]NewRelicMetricDescriptor, 0, len(NRMetricDescriptors))},
	}
	t := reflect.TypeOf(*s.Stats)
	v := reflect.ValueOf(*s.Stats)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		name := f.Tag.Get("mapstructure")
		md, err := getNewRelicMetric(name)
		if err != nil {
			if err == ErrNotFound {
				// If we haven't defined the metric, just ignore it
				continue
			}
			continue
		}

		md.Name = fmt.Sprintf("fastly.%s", name)
		md.Value = v.Field(i).Interface()
		md.Timestamp = int64(s.IntervalStart)
		metrics[0].Metrics = append(metrics[0].Metrics, md)
	}

	return metrics
}

func (n *NewRelicExporter) report(ctx context.Context, r []NewRelicMetricReport) error {
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(r); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, nrEndpoint, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Api-Key", n.insertKey)

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute http request: %w", err)
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("invalid response status '%s'", resp.Status)
	}

	return nil
}

func (n *NewRelicExporter) Run(ctx context.Context) {
	l := zap.S()
	for {
		select {
		case s := <-n.ch:
			report := n.buildMetrics(s)
			if err := n.report(ctx, report); err != nil {
				l.Errorf("failed to report to New Relic: %v", err)
				continue
			}
			l.Debugf("successfully reported to new relic")
		case <-ctx.Done():
			return
		}
	}
}
