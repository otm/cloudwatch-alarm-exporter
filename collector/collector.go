package collector

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "cloudwatch"
	subsystem = "alarm"
	collector = namespace + subsystem
)

type alarmFetcher interface {
	Alarms() ([]*cloudwatch.MetricAlarm, error)
}

// CloudWatchAlarms implements the Prometheus Collector interface
type CloudWatchAlarms struct {
	desc                *prometheus.Desc
	alarmFetcher        alarmFetcher
	scrapeStatus        *prometheus.GaugeVec
	totalScrapeDuration prometheus.Summary
	totalFetchedAlarms  prometheus.Counter
	metrics             []prometheus.Metric
}

// New returns a CloudWatchAlarms instance
func New(alarmFetcher alarmFetcher) *CloudWatchAlarms {
	return &CloudWatchAlarms{
		totalFetchedAlarms: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: collector,
				Name:      "total_fetched_cloudwatchalarms",
				Help:      "Alarms fetched from CloudWatch.",
			},
		),
		alarmFetcher: alarmFetcher,
		totalScrapeDuration: prometheus.NewSummary(
			prometheus.SummaryOpts{
				Namespace: collector,
				Name:      "total_scrape_duration",
				Help:      "The time the scrape of CloudWatch alarms took in seconds.",
			},
		),
		scrapeStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: collector,
				Name:      "up",
				Help:      "Status of the current scrape",
			},
			[]string{"collector"},
		),
	}
}

// Collect implements the Prometheus Collecter interface
func (c *CloudWatchAlarms) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	var alarmCount int
	alarms, err := c.alarmFetcher.Alarms()

	if err != nil {
		c.scrapeStatus.WithLabelValues(collector).Set(0)
		c.scrapeStatus.Collect(ch)

		elapsed := time.Since(start)
		c.totalScrapeDuration.Observe(float64(elapsed.Seconds()))
		ch <- c.totalScrapeDuration

		log.Printf("unable to fetch alarms: %v", err)
		return
	}
	c.metrics = nil

	scrapeStatus := true
	for _, alarm := range alarms {
		alarmCount++

		value, err := translateCloudWatchState(*alarm.StateValue)
		if err != nil {
			scrapeStatus = false
			continue
		}

		metric := prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "state"),
				"Cloudwatch alarm state: 0=OK, 1=Alarm, -1=InsufficientData, -1000=InternalCollectorError",
				[]string{"name", "namespace"},
				nil,
			),
			prometheus.GaugeValue,
			value,
			*alarm.AlarmName, *alarm.Namespace,
		)
		ch <- metric
		c.metrics = append(c.metrics, metric)
	}

	c.scrapeStatus.WithLabelValues(collector).Set(1)
	if !scrapeStatus {
		c.scrapeStatus.WithLabelValues(collector).Set(0)
	}
	c.scrapeStatus.Collect(ch)

	c.totalFetchedAlarms.Add(float64(alarmCount))

	elapsed := time.Since(start)
	c.totalScrapeDuration.Observe(float64(elapsed.Seconds()))
	ch <- c.totalScrapeDuration
	log.Println(elapsed.Seconds())
}

// Describe implements the Prometheus Describer interface
func (c *CloudWatchAlarms) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.Desc()
	}
	c.scrapeStatus.Describe(ch)
	c.totalScrapeDuration.Describe(ch)
	c.totalFetchedAlarms.Describe(ch)
}

func translateCloudWatchState(state string) (float64, error) {
	switch state {
	case cloudwatch.StateValueOk:
		return float64(0), nil
	case cloudwatch.StateValueAlarm:
		return float64(1), nil
	case cloudwatch.StateValueInsufficientData:
		return float64(-1), nil
	default:
		return 0, fmt.Errorf("unknown cloudwatch state: %s", state)
	}
}
