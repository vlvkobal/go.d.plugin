// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorBlockedTask = "PendingTasks"
	metricBlockedType    = "org_apache_cassandra_metrics_threadpools_count"
)

func doCollectBlockedTask(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, metricBlockedType, collectorBlockedTask, false)
	return enabled && success
}

func collectBlockedTask(pms prometheus.Metrics) *blockedTask {
	if !doCollectBlockedTask(pms) {
		return nil
	}

	var bt blockedTask
	collectBlockedTaskByType(&bt, pms)

	return &bt
}

func collectBlockedTaskByType(bt *blockedTask, pms prometheus.Metrics) {
	var values blockedTask
	for _, pm := range pms.FindByName(metricBlockedType) {
		metricScope := pm.Labels.Get("scope")
		metricName := pm.Labels.Get("name")
		if metricName == "CurrentlyBlockedTasks" {
			assignBlockedTaskMetric(&values, metricScope, pm.Value)
		}
	}
	bt.task = values.task
}

func assignBlockedTaskMetric(pt *blockedTask, scope string, value float64) {
	switch scope {
	default:
	case "CounterMutationStage":
		pt.task += int64(value)
	case "MutationStage":
		pt.task += int64(value)
	case "ReadRepairStage":
		pt.task += int64(value)
	case "ReadStage":
		pt.task += int64(value)
	case "RequestResponseStage":
		pt.task += int64(value)
	}
}