package metrics

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Metrics struct {
	totalRequestsCount int64
	getPlaylistCount   int64
	getTrackCount      int64
}

var metrics *Metrics

func GetMetricsHandler() *Metrics {
	if metrics == nil {
		metrics = &Metrics{}
	}

	return metrics
}

type MetricEvent int64

const (
	GetTrackEvent MetricEvent = iota
	GetPlaylistEvent
)

func (c *Metrics) CountRequests(event MetricEvent) {
	switch event {
	case GetTrackEvent:
		c.getTrackCount++
	case GetPlaylistEvent:
		c.getPlaylistCount++
	}
	c.totalRequestsCount++
}

func (c *Metrics) Description() string {
	return "Provides metrics about bot statistic"
}

func (c *Metrics) SampleConfig() string {
	return ""
}

func (c *Metrics) reset() {
	c.totalRequestsCount = 0
	c.getPlaylistCount = 0
	c.getTrackCount = 0
}

func (c *Metrics) isEmpty() bool {
	return c.getPlaylistCount == 0 && c.getTrackCount == 0 && c.totalRequestsCount == 0
}

func (c *Metrics) Gather(acc telegraf.Accumulator) error {
	if c.isEmpty() {
		return nil
	}

	acc.AddCounter(
		"metrics",
		map[string]interface{}{"total_count": c.totalRequestsCount},
		map[string]string{"name": "total_count"},
	)
	acc.AddCounter(
		"metrics",
		map[string]interface{}{"get_playlist_count": c.getPlaylistCount},
		map[string]string{"name": "get_playlist_count"},
	)
	acc.AddCounter(
		"metrics",
		map[string]interface{}{"get_track_count": c.getTrackCount},
		map[string]string{"name": "get_track_count"},
	)
	defer c.reset()
	return nil
}

func init() {
	inputs.Add("metrics", func() telegraf.Input {
		return GetMetricsHandler()
	})
}
