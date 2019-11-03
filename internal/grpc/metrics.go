package grpc

import "github.com/prometheus/client_golang/prometheus"

var apiCreateEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_create_event_count",
	Help:        "API create event",
	ConstLabels: prometheus.Labels{"api": "create"},
})

var apiDeleteEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_delete_event_count",
	Help:        "API delete event",
	ConstLabels: prometheus.Labels{"api": "delete"},
})

var apiUpdateEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_update_event_count",
	Help:        "API update event",
	ConstLabels: prometheus.Labels{"api": "update"},
})

var apiListEventsCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_list_events_count",
	Help:        "API list event",
	ConstLabels: prometheus.Labels{"api": "list"},
})

var apiGetEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_get_event_count",
	Help:        "API get event",
	ConstLabels: prometheus.Labels{"api": "get"},
})

var apiCreateEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_create_event_error_count",
	Help:        "API create event error",
	ConstLabels: prometheus.Labels{"api": "create"},
})

var apiDeleteEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_delete_event_error_count",
	Help:        "API delete event error",
	ConstLabels: prometheus.Labels{"api": "delete"},
})

var apiUpdateEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_update_event_error_count",
	Help:        "API update event error",
	ConstLabels: prometheus.Labels{"api": "update"},
})

var apiListEventsErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_list_events_error_count",
	Help:        "API list event error",
	ConstLabels: prometheus.Labels{"api": "list"},
})

var apiGetEventErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name:        "api_get_event_error_count",
	Help:        "API get event error",
	ConstLabels: prometheus.Labels{"api": "get"},
})

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
