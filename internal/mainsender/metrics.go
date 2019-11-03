package mainsender

import "github.com/prometheus/client_golang/prometheus"

var senderEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "sender_event_count",
	Help: "Send event",
})

var senderEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "sender_event_error_count",
	Help: "Send error event",
})

func init() {
	prometheus.MustRegister(senderEventCounter)
	prometheus.MustRegister(senderEventErrorCounter)
}
