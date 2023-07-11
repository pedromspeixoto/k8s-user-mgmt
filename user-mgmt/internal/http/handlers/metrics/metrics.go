package metrics

import (
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type MetricServiceHandler interface {
	// Routes creates a REST router for the health service
	Routes() chi.Router
}

type metricServiceHandler struct {
}

func NewMetricServiceHandler() MetricServiceHandler {
	return &metricServiceHandler{}
}

func (h metricServiceHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.GetMetrics())

	return r
}

// Get metrics - Retrieve service metrics
// @Summary Get service metrics.
// @Description This API is used to get the service metrics.
// @Tags metrics
// @Accept  json
// @Produce  json
// @Router /metrics [get]
func (h metricServiceHandler) GetMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}
