package fetch

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.Handle("/metrics", promhttp.Handler())
}
