package fastlystats

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// timeSeriesBatchSize is the number of time series to send in one request.
// TimeSeries will be batched and this number of time series will be sent with
// each request.
//
// This is limited by a quota, see https://cloud.google.com/monitoring/quotas
// This cannot be higher than 200 as of writing this.
const timeSeriesBatchSize = 200

// maxReportTimeout is the maximum time to wait for a batch to reach `timeSeriesBatchSize`
// entries before sending a batch with fewer time series.
const maxReportTimeout = 2 * time.Second

type StackdriverExporter struct {
	metricClient *monitoring.MetricClient

	ch                 <-chan *FastlyMeanStats
	timeSeriesCh       chan *monitoringpb.TimeSeries
	googleCloudProject string
	nodeID             string
}

func SetupMetricDescriptors(ctx context.Context, googleCloudProject string) {
	ll := zap.S()
	metricClient, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		ll.Fatal(err)
	}

	ll.Infof("Setting up metric descriptors")
	for _, m := range MetricDescriptors {
		name := fmt.Sprintf("projects/%s/metricDescriptors/%s", googleCloudProject, m.Type)

		ll.Infof("Recreating metric '%s'", m.Type)
		err = metricClient.DeleteMetricDescriptor(ctx, &monitoringpb.DeleteMetricDescriptorRequest{
			Name: name,
		})
		if status.Code(err) != codes.OK && status.Code(err) != codes.NotFound {
			ll.Warnf("Failed to delete metric (will attempt to create anyway): %v", err)
		}

		_, err := metricClient.CreateMetricDescriptor(ctx, &monitoringpb.CreateMetricDescriptorRequest{
			Name:             fmt.Sprintf("projects/%s", googleCloudProject),
			MetricDescriptor: m,
		})
		if err != nil {
			ll.Warnf("Failed to create metric '%s': %v", m.Type, err)
		}
	}
}

func NewStackdriverExporter(project string, nodeID string, ch <-chan *FastlyMeanStats) (*StackdriverExporter, error) {
	metricClient, err := monitoring.NewMetricClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &StackdriverExporter{
		metricClient:       metricClient,
		ch:                 ch,
		timeSeriesCh:       make(chan *monitoringpb.TimeSeries, timeSeriesBatchSize),
		googleCloudProject: project,
		nodeID:             nodeID,
	}, nil
}

func (s *StackdriverExporter) Run(ctx context.Context) {
	ll := zap.S()
	ll.Infof("starting stackdriver exporter to project %s", s.googleCloudProject)

	tasks := 2
	wg := sync.WaitGroup{}
	wg.Add(tasks)

	ictx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		defer wg.Done()
		s.timeSeriesWorker(ictx)
	}()

	go func() {
		defer cancel()
		defer wg.Done()
		s.timeSeriesReporter(ictx)
	}()

	wg.Wait()
}

func (s *StackdriverExporter) timeSeries(stats *FastlyMeanStats) []*monitoringpb.TimeSeries {
	var result []*monitoringpb.TimeSeries

	t := reflect.TypeOf(*stats.Stats)
	v := reflect.ValueOf(*stats.Stats)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		metricName := field.Tag.Get("mapstructure")

		metricKind := metric.MetricDescriptor_GAUGE

		valueType := metric.MetricDescriptor_VALUE_TYPE_UNSPECIFIED
		switch t.Field(i).Type.Kind() {
		case reflect.Uint64:
			valueType = metric.MetricDescriptor_INT64
		case reflect.Float64:
			valueType = metric.MetricDescriptor_DOUBLE
		}

		var getValue valuer
		switch valueType {
		case metric.MetricDescriptor_INT64:
			getValue = int64Valuer
		case metric.MetricDescriptor_DOUBLE:
			getValue = doubleValuer
		default:
			zap.S().Debugf("no valuer for type %v: skipping metric '%s'", valueType, metricName)
			continue
		}

		ts := &monitoringpb.TimeSeries{
			Metric: &metric.Metric{
				Type: fmt.Sprintf("custom.googleapis.com/fastly/%s", metricName),
			},
			MetricKind: metricKind,
			ValueType:  valueType,
			Points:     []*monitoringpb.Point{{Value: getValue(v.Field(i))}},
		}

		result = append(result, ts)
	}

	return result
}

func (s *StackdriverExporter) sendTimeSeries(
	ctx context.Context,
	intervalStart, intervalEnd uint64,
	resource *monitoredres.MonitoredResource,
	timeSeries []*monitoringpb.TimeSeries,
) error {
	for i, ts := range timeSeries {

		// Set some common values for all time series
		timeSeries[i].Resource = resource
		for j := range ts.Points {
			timeSeries[i].Points[j].Interval = &monitoringpb.TimeInterval{
				// StartTime is set to End Time, because it's not supported that these
				// differ yet (it gives an API error).
				StartTime: timestamppb.New(time.Unix(int64(intervalEnd), 0)),
				EndTime:   timestamppb.New(time.Unix(int64(intervalEnd), 0)),
			}
		}

		// Send it on the channel for batching and reporting
		select {
		case s.timeSeriesCh <- timeSeries[i]:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (s *StackdriverExporter) timeSeriesWorker(ctx context.Context) {
	ll := zap.S().With("type", "producer")

	ll.Infof("started worker")
	for {
		select {
		case meanStats, ok := <-s.ch:
			if !ok {
				ll.Infof("channel closed, not exporting any more stats")
				return
			}

			monitoredResource := &monitoredres.MonitoredResource{
				Type: "generic_node",
				Labels: map[string]string{
					"project_id": s.googleCloudProject,
					"location":   "global",
					"namespace":  "fastly",
					"node_id":    s.nodeID,
				},
			}

			if err := s.sendTimeSeries(ctx, meanStats.IntervalStart, meanStats.IntervalEnd, monitoredResource, s.timeSeries(meanStats)); err != nil {
				zap.S().Warnf("failed to send time series: %v", err)
			}

		case <-ctx.Done():
			ll.Infof("context done, exiting")
			return
		}
	}
}

func (s *StackdriverExporter) timeSeriesReporter(ctx context.Context) {
	ll := zap.S().With("type", "reporter")

	batch := make([]*monitoringpb.TimeSeries, 0, timeSeriesBatchSize)
	var timeoutCh <-chan time.Time

	report := func() {
		if err := s.reportBatch(ctx, batch); err != nil {
			ll.Warnf("failed to report timeseries: %v", err)
		}
		batch = batch[:0]
		timeoutCh = nil
	}

	for {
		select {
		case <-timeoutCh:
			report()

		case ts, ok := <-s.timeSeriesCh:
			if !ok {
				ll.Infof("time series chan closed, exiting")
				return
			}

			batch = append(batch, ts)
			if timeoutCh == nil {
				timeoutCh = time.After(maxReportTimeout)
			}

			if len(batch) == cap(batch) {
				report()
			}

		case <-ctx.Done():
			ll.Infof("context done, exiting")
			return
		}
	}
}

func (s *StackdriverExporter) reportBatch(ctx context.Context, batch []*monitoringpb.TimeSeries) error {
	if len(batch) == 0 {
		return nil
	}

	start := time.Now()
	err := s.metricClient.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name:       fmt.Sprintf("projects/%s", s.googleCloudProject),
		TimeSeries: batch,
	})
	if err != nil {
		return fmt.Errorf("failed to create %d time series: %v", len(batch), err)
	}

	zap.S().Debugf("successfully reported %d time series in %v", len(batch), time.Since(start))

	return nil
}
