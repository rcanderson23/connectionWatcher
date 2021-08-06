package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// NewConnections is a counter for the number of observed new connections
	NewConnections = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "proc_net_tcp_new_connections",
			Help: "New connections observed at /proc/net/tcp",
		})
)
