package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	httpRequestTotalMetric           = "http_request_total"  // done
	httpResponseTotalMetric          = "http_response_total" // done
	httpResponseStatusMetric         = "http_response_status_total"
	httpRequestPerRouteMetric        = "http_requests_per_route_total" // done
	httpResponseStatusPerRouteMetric = "http_response_status_per_route_total"
	httpRequestPerConsumerMetric     = "http_requests_per_consumer_total"  // done
	requestDurationMetric            = "http_request_duration_millisecond" // done
	requestTrafficMetric             = "http_request_traffic_total"        // done
	responseTrafficMetric            = "http_response_traffic_total"       // done
)

var (
	totalRequestReceived = prometheus.NewCounter(prometheus.CounterOpts{
		Name: httpRequestTotalMetric,
		Help: "Total number of HTTP requests received",
	})

	totalResponseSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: httpResponseTotalMetric,
		Help: "Total number of HTTP responses sent",
	})

	httpResponseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: httpResponseStatusMetric,
			Help: "Total number of HTTP responses by status code",
		},
		[]string{"status_code"},
	)

	httpRequestsPerRoute = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: httpRequestPerRouteMetric,
			Help: "Total number of HTTP requests per route",
		},
		[]string{"route"},
	)

	httpResponseStatusPerRoute = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: httpResponseStatusPerRouteMetric,
			Help: "Total number of HTTP responses by status code and route",
		},
		[]string{"status_code", "route"},
	)

	httpRequestsPerConsumer = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: httpRequestPerConsumerMetric,
			Help: "Total number of HTTP requests per route",
		},
		[]string{"consumer"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    requestDurationMetric,
			Help:    "Histogram of response latency (millisecond) of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route"},
	)

	requestTraffic = prometheus.NewCounter(prometheus.CounterOpts{
		Name: requestTrafficMetric,
		Help: "Total number of HTTP requests traffic",
	},
	)

	responseTraffic = prometheus.NewCounter(prometheus.CounterOpts{
		Name: responseTrafficMetric,
		Help: "Total number of HTTP responses traffic",
	},
	)
)

func init() {
	prometheus.MustRegister(
		totalRequestReceived,
		totalResponseSent,
		httpResponseStatus,
		httpRequestsPerRoute,
		httpResponseStatusPerRoute,
		httpRequestsPerConsumer,
		requestTraffic,
		responseTraffic,
		requestDuration)
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		consumer := c.GetHeader("X-Consumer-Groups")
		httpRequestsPerConsumer.WithLabelValues(consumer).Inc()

		// Increment the counter for the total request received
		totalRequestReceived.Inc()

		// Increment the counter for the specific route
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		httpRequestsPerRoute.WithLabelValues(route).Inc()

		// Record the request size
		requestTraffic.Add(float64(c.Request.ContentLength))

		c.Next()

		statusCode := c.Writer.Status()

		// Increment the counter for the total response sent
		totalResponseSent.Inc()

		// Increment the counter for the specific route
		httpResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()

		// Increment the counter for the specific route and status code
		httpResponseStatusPerRoute.WithLabelValues(strconv.Itoa(statusCode), route).Inc()

		// Record the response size
		responseTraffic.Add(float64(c.Writer.Size()))

		duration := time.Since(startTime).Milliseconds()

		// Record the request duration
		requestDuration.WithLabelValues(route).Observe(float64(duration))
	}
}

// CORSMiddleware allow all origins (CORS open) access to the API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("User-Name", c.GetHeader("User-Name"))
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
