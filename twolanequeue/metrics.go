package twolanequeue

import (
	"sync"

	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

var (
	twoLaneQueueDepth = metrics.NewGaugeVec(
		&metrics.GaugeOpts{
			Name: "two_lane_queue_depth",
			Help: "number of items in the two lane queue",
		},
		[]string{"name", "lane"},
	)
)

var registerMetricOnce sync.Once

func RegisterMetrics() {
	registerMetricOnce.Do(func() {
		legacyregistry.MustRegister(twoLaneQueueDepth)
	})
}
