package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	apiCreateEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_create_event_count",
		Help:        "API create event",
		ConstLabels: prometheus.Labels{"api": "create"},
	})

	apiDeleteEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_delete_event_count",
		Help:        "API delete event",
		ConstLabels: prometheus.Labels{"api": "delete"},
	})

	apiUpdateEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_update_event_count",
		Help:        "API update event",
		ConstLabels: prometheus.Labels{"api": "update"},
	})

	apiListEventsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_list_events_count",
		Help:        "API list event",
		ConstLabels: prometheus.Labels{"api": "list"},
	})

	apiGetEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_get_event_count",
		Help:        "API get event",
		ConstLabels: prometheus.Labels{"api": "get"},
	})

	apiCreateEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_create_event_error_count",
		Help:        "API create event error",
		ConstLabels: prometheus.Labels{"api": "create"},
	})

	apiDeleteEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_delete_event_error_count",
		Help:        "API delete event error",
		ConstLabels: prometheus.Labels{"api": "delete"},
	})

	apiUpdateEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_update_event_error_count",
		Help:        "API update event error",
		ConstLabels: prometheus.Labels{"api": "update"},
	})

	apiListEventsErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_list_events_error_count",
		Help:        "API list event error",
		ConstLabels: prometheus.Labels{"api": "list"},
	})

	apiGetEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "api_get_event_error_count",
		Help:        "API get event error",
		ConstLabels: prometheus.Labels{"api": "get"},
	})
)

func init() {
	prometheus.MustRegister(apiCreateEventCounter)
	prometheus.MustRegister(apiGetEventCounter)
	prometheus.MustRegister(apiDeleteEventCounter)
	prometheus.MustRegister(apiUpdateEventCounter)
	prometheus.MustRegister(apiListEventsCounter)
	prometheus.MustRegister(apiCreateEventErrorCounter)
	prometheus.MustRegister(apiGetEventErrorCounter)
	prometheus.MustRegister(apiDeleteEventErrorCounter)
	prometheus.MustRegister(apiUpdateEventErrorCounter)
	prometheus.MustRegister(apiListEventsErrorCounter)
}
