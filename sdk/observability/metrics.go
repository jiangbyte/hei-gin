package observability

import (
	"database/sql"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"hei-gin/sdk/infra/db"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "hei",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total HTTP requests processed.",
		},
		[]string{"method", "route", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "hei",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)
	httpInflightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "http",
			Name:      "inflight_requests",
			Help:      "Current number of in-flight HTTP requests.",
		},
	)
	httpPanicsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "hei",
			Subsystem: "http",
			Name:      "panics_total",
			Help:      "Total number of recovered HTTP panics.",
		},
	)
	wsConnectionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "hei",
			Subsystem: "ws",
			Name:      "connections_total",
			Help:      "Total accepted WebSocket connections.",
		},
	)
	wsRejectedConnectionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "hei",
			Subsystem: "ws",
			Name:      "rejected_connections_total",
			Help:      "Total rejected WebSocket connections.",
		},
	)
	wsDisconnectedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "hei",
			Subsystem: "ws",
			Name:      "disconnections_total",
			Help:      "Total disconnected WebSocket connections.",
		},
	)
	wsMessagesSentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "hei",
			Subsystem: "ws",
			Name:      "messages_sent_total",
			Help:      "Total WebSocket messages sent.",
		},
		[]string{"channel"},
	)
	wsCurrentConnections = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "ws",
			Name:      "current_connections",
			Help:      "Current number of WebSocket connections.",
		},
		func() float64 {
			return float64(atomic.LoadInt64(&wsCurrentConnectionsValue))
		},
	)
)

var wsCurrentConnectionsValue int64

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpInflightRequests,
		httpPanicsTotal,
		wsConnectionsTotal,
		wsRejectedConnectionsTotal,
		wsDisconnectedTotal,
		wsMessagesSentTotal,
		wsCurrentConnections,
	)

	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "db",
			Name:      "open_connections",
			Help:      "Current number of open database connections.",
		},
		func() float64 { return float64(databaseStats().OpenConnections) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "db",
			Name:      "in_use_connections",
			Help:      "Current number of in-use database connections.",
		},
		func() float64 { return float64(databaseStats().InUse) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "db",
			Name:      "idle_connections",
			Help:      "Current number of idle database connections.",
		},
		func() float64 { return float64(databaseStats().Idle) },
	))
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "hei",
			Subsystem: "db",
			Name:      "wait_count_total",
			Help:      "Total database connection wait count.",
		},
		func() float64 { return float64(databaseStats().WaitCount) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "db",
			Name:      "wait_duration_seconds",
			Help:      "Total time blocked waiting for a new database connection.",
		},
		func() float64 { return databaseStats().WaitDuration.Seconds() },
	))

	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "redis",
			Name:      "pool_hits",
			Help:      "Redis connection pool hits.",
		},
		func() float64 { return float64(redisPoolHits()) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "redis",
			Name:      "pool_misses",
			Help:      "Redis connection pool misses.",
		},
		func() float64 { return float64(redisPoolMisses()) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "redis",
			Name:      "pool_timeouts",
			Help:      "Redis connection pool timeouts.",
		},
		func() float64 { return float64(redisPoolTimeouts()) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "redis",
			Name:      "total_connections",
			Help:      "Current total Redis pool connections.",
		},
		func() float64 { return float64(redisTotalConns()) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "redis",
			Name:      "idle_connections",
			Help:      "Current idle Redis pool connections.",
		},
		func() float64 { return float64(redisIdleConns()) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "hei",
			Subsystem: "redis",
			Name:      "stale_connections",
			Help:      "Current stale Redis pool connections.",
		},
		func() float64 { return float64(redisStaleConns()) },
	))
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func ObserveHTTPRequest(method, route string, status int, seconds float64) {
	statusText := strconv.Itoa(status)
	httpRequestsTotal.WithLabelValues(method, route, statusText).Inc()
	httpRequestDuration.WithLabelValues(method, route, statusText).Observe(seconds)
}

func IncHTTPInflight() {
	httpInflightRequests.Inc()
}

func DecHTTPInflight() {
	httpInflightRequests.Dec()
}

func IncHTTPPanic() {
	httpPanicsTotal.Inc()
}

func IncWSConnection() {
	wsConnectionsTotal.Inc()
	atomic.AddInt64(&wsCurrentConnectionsValue, 1)
}

func DecWSConnection() {
	wsDisconnectedTotal.Inc()
	atomic.AddInt64(&wsCurrentConnectionsValue, -1)
}

func IncWSRejected() {
	wsRejectedConnectionsTotal.Inc()
}

func ObserveWSMessage(channel string, count int) {
	if count <= 0 {
		return
	}
	wsMessagesSentTotal.WithLabelValues(channel).Add(float64(count))
}

func databaseStats() sql.DBStats {
	if db.DB == nil {
		return sql.DBStats{}
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		return sql.DBStats{}
	}
	return sqlDB.Stats()
}

func redisPoolHits() uint32 {
	if db.Redis == nil {
		return 0
	}
	return db.Redis.PoolStats().Hits
}

func redisPoolMisses() uint32 {
	if db.Redis == nil {
		return 0
	}
	return db.Redis.PoolStats().Misses
}

func redisPoolTimeouts() uint32 {
	if db.Redis == nil {
		return 0
	}
	return db.Redis.PoolStats().Timeouts
}

func redisTotalConns() uint32 {
	if db.Redis == nil {
		return 0
	}
	return db.Redis.PoolStats().TotalConns
}

func redisIdleConns() uint32 {
	if db.Redis == nil {
		return 0
	}
	return db.Redis.PoolStats().IdleConns
}

func redisStaleConns() uint32 {
	if db.Redis == nil {
		return 0
	}
	return db.Redis.PoolStats().StaleConns
}
