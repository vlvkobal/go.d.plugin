// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"errors"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

const (
	metricCollectorSuccess = "org_apache_cassandra_metrics_clientrequest_count"
)

func isValidCassandraMetrics(pms prometheus.Metrics) bool {
	return pms.FindByName(metricCollectorSuccess).Len() > 0
}

func (c *Cassandra) collect() (map[string]int64, error) {
	pms, err := c.prom.Scrape()
	if err != nil {
		return nil, err
	}

	if !isValidCassandraMetrics(pms) {
		return nil, errors.New("collected metrics aren't cassandra metrics")
	}

	mx := collect(pms)

	return stm.ToMap(mx), nil
}

func collect(pms prometheus.Metrics) *metrics {
	mx := metrics{
		ioThroughput: collectThroughput(pms),
		ioLatency:    collectLatency(pms),
		hit:          collectCache(pms),
		hd:           collectDisk(pms),
		gcCount:      collectGC(pms, metricGCCount),
		gcTime:       collectGC(pms, metricGCTime),
		etimeout:     collectRequestError(pms, collectorTimeout),
		eunavailable: collectRequestError(pms, collectorUnavailable),
		pTask:        collectPendingTask(pms),
		bTask:        collectBlockedTask(pms),
	}

	return &mx
}

func checkCollector(pms prometheus.Metrics, metric string, testValue string, testScope bool) (enabled, success bool) {
	for _, pm := range pms.FindByName(metric) {
		metricScope := pm.Labels.Get("scope")
		metricName := pm.Labels.Get("name")
		metricGC := pm.Labels.Get("gc")
		// FOr some metrics we need to verify scope, for others name, so we test both
		if metricName == testValue && !testScope {
			return true, true
		} else if metricScope == testValue && testScope {
			return true, true
		} else if metricGC == testValue {
			return true, true
		}
	}
	return false, false
}