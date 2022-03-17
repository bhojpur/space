package server

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"net/http"

	"github.com/bhojpur/space/pkg/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricDescriptions = map[string]*prometheus.Desc{
		/*
			these metrics are taken from basicStats() / extStats()
			by accessing the map and directly exporting the value found
		*/
		"num_collections":           prometheus.NewDesc("bhojpur_collections", "Total number of collections", nil, nil),
		"pid":                       prometheus.NewDesc("bhojpur_pid", "", nil, nil),
		"aof_size":                  prometheus.NewDesc("bhojpur_aof_size_bytes", "", nil, nil),
		"num_hooks":                 prometheus.NewDesc("bhojpur_hooks", "", nil, nil),
		"in_memory_size":            prometheus.NewDesc("bhojpur_in_memory_size_bytes", "", nil, nil),
		"heap_size":                 prometheus.NewDesc("bhojpur_heap_size_bytes", "", nil, nil),
		"heap_released":             prometheus.NewDesc("bhojpur_memory_reap_released_bytes", "", nil, nil),
		"max_heap_size":             prometheus.NewDesc("bhojpur_memory_max_heap_size_bytes", "", nil, nil),
		"avg_item_size":             prometheus.NewDesc("bhojpur_avg_item_size_bytes", "", nil, nil),
		"pointer_size":              prometheus.NewDesc("bhojpur_pointer_size_bytes", "", nil, nil),
		"cpus":                      prometheus.NewDesc("bhojpur_num_cpus", "", nil, nil),
		"bhojpur_connected_clients": prometheus.NewDesc("bhojpur_connected_clients", "", nil, nil),

		"bhojpur_total_connections_received": prometheus.NewDesc("bhojpur_connections_received_total", "", nil, nil),
		"bhojpur_total_messages_sent":        prometheus.NewDesc("bhojpur_messages_sent_total", "", nil, nil),
		"bhojpur_expired_keys":               prometheus.NewDesc("bhojpur_expired_keys_total", "", nil, nil),

		/*
			these metrics are NOT taken from basicStats() / extStats()
			but are calculated independently
		*/
		"collection_objects": prometheus.NewDesc("bhojpur_collection_objects", "Total number of objects per collection", []string{"col"}, nil),
		"collection_points":  prometheus.NewDesc("bhojpur_collection_points", "Total number of points per collection", []string{"col"}, nil),
		"collection_strings": prometheus.NewDesc("bhojpur_collection_strings", "Total number of strings per collection", []string{"col"}, nil),
		"collection_weight":  prometheus.NewDesc("bhojpur_collection_weight_bytes", "Total weight of collection in bytes", []string{"col"}, nil),
		"server_info":        prometheus.NewDesc("bhojpur_server_info", "Server info", []string{"id", "version"}, nil),
		"replication":        prometheus.NewDesc("bhojpur_replication_info", "Replication info", []string{"role", "following", "caught_up", "caught_up_once"}, nil),
		"start_time":         prometheus.NewDesc("bhojpur_start_time_seconds", "", nil, nil),
	}

	cmdDurations = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "bhojpur_cmd_duration_seconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
	}, []string{"cmd"},
	)
)

func (s *Server) MetricsIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html><head>
<title>Bhojpur Space - Tile Server ` + core.Version + `</title></head>
<body><h1>Bhojpur Space - Tile Server ` + core.Version + `</h1>
<p><a href='/metrics'>Metrics</a></p>
</body></html>`))
}

func (s *Server) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
		prometheus.NewBuildInfoCollector(),
		cmdDurations,
		s,
	)

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func (s *Server) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range metricDescriptions {
		ch <- desc
	}
}

func (s *Server) Collect(ch chan<- prometheus.Metric) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m := make(map[string]interface{})
	s.basicStats(m)
	s.extStats(m)

	for metric, descr := range metricDescriptions {
		if val, ok := m[metric].(int); ok {
			ch <- prometheus.MustNewConstMetric(descr, prometheus.GaugeValue, float64(val))
		} else if val, ok := m[metric].(float64); ok {
			ch <- prometheus.MustNewConstMetric(descr, prometheus.GaugeValue, val)
		}
	}

	ch <- prometheus.MustNewConstMetric(
		metricDescriptions["server_info"],
		prometheus.GaugeValue, 1.0,
		s.config.serverID(), core.Version)

	ch <- prometheus.MustNewConstMetric(
		metricDescriptions["start_time"],
		prometheus.GaugeValue, float64(s.started.Unix()))

	replLbls := []string{"leader", "", "", ""}
	if s.config.followHost() != "" {
		replLbls = []string{"follower",
			fmt.Sprintf("%s:%d", s.config.followHost(), s.config.followPort()),
			fmt.Sprintf("%t", s.fcup), fmt.Sprintf("%t", s.fcuponce)}
	}
	ch <- prometheus.MustNewConstMetric(
		metricDescriptions["replication"],
		prometheus.GaugeValue, 1.0,
		replLbls...)

	/*
		add objects/points/strings stats for each collection
	*/
	s.cols.Ascend(nil, func(v interface{}) bool {
		c := v.(*collectionKeyContainer)
		ch <- prometheus.MustNewConstMetric(
			metricDescriptions["collection_objects"],
			prometheus.GaugeValue,
			float64(c.col.Count()),
			c.key,
		)
		ch <- prometheus.MustNewConstMetric(
			metricDescriptions["collection_points"],
			prometheus.GaugeValue,
			float64(c.col.PointCount()),
			c.key,
		)
		ch <- prometheus.MustNewConstMetric(
			metricDescriptions["collection_strings"],
			prometheus.GaugeValue,
			float64(c.col.StringCount()),
			c.key,
		)
		ch <- prometheus.MustNewConstMetric(
			metricDescriptions["collection_weight"],
			prometheus.GaugeValue,
			float64(c.col.TotalWeight()),
			c.key,
		)
		return true
	})
}
