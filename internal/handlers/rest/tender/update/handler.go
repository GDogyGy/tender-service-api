package update

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"TenderServiceApi/internal/model"
)

type log interface {
	Error(msg string, args ...any)
}

type useCaseTenderEdite interface {
	Edit(ctx context.Context, id string, username string, tenderNew model.Tender) (model.Tender, error)
	Rollback(ctx context.Context, id string, username string, version string) (model.Tender, error)
	Status(ctx context.Context, username string, tenderId string, status string) (model.Tender, error)
}

type prometheusMiddleware interface {
	PrometheusMiddleware(handlerName string, next http.HandlerFunc) http.Handler
}

type producer interface {
	SendEvent(eventType string, model interface{}) error
	SendAsync(eventType string, model interface{})
}

type Handler struct {
	log         log
	prometheus  prometheusMiddleware
	producer    producer
	tenderEdite useCaseTenderEdite
}

var tenderEditeRegexp = regexp.MustCompile(`/tenders/(.*)/edit`)
var tenderIdRollbackRegexp = regexp.MustCompile(`/tenders/(.*)/rollback/`)
var tenderVersionRollbackRegexp = regexp.MustCompile(`/rollback/(.*)\?`)

func NewHandler(l log, pm prometheusMiddleware, p producer, t useCaseTenderEdite) Handler {
	return Handler{
		l, pm, p, t,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.Handle(http.MethodPatch+" /api/tenders/{id}/edit", h.prometheus.PrometheusMiddleware("tenderEdit", h.Edit))
	router.Handle(http.MethodPut+" /api/tenders/{id}/rollback/{version}", h.prometheus.PrometheusMiddleware("tenderRollback", h.Rollback))
	router.Handle(http.MethodPut+" /api/tenders/status", h.prometheus.PrometheusMiddleware("tenderStatus", h.Status))
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rq := r.URL.Query()
	user := "username"

	if rq.Get(user) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := tenderEditeRegexp.FindStringSubmatch(r.RequestURI)
	if id == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tenderRequest transport.TenderEditRequest
	err = json.Unmarshal(b, &tenderRequest)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := h.tenderEdite.Edit(r.Context(), id[1], rq.Get(user), convert.TenderReqEditTransportToModel(tenderRequest))
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err = json.Marshal(convert.TenderModelToTransport(resp))
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Rollback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rq := r.URL.Query()
	user := "username"

	if rq.Get(user) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := tenderIdRollbackRegexp.FindStringSubmatch(r.RequestURI)
	if id == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	version := tenderVersionRollbackRegexp.FindStringSubmatch(r.RequestURI)
	if version == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tender, err := h.tenderEdite.Rollback(r.Context(), id[1], rq.Get(user), version[1])
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, model.NotFindResponsible) {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(convert.TenderModelToTransport(tender))
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rq := r.URL.Query()
	user := "username"
	tenderId := "tenderId"
	status := "status"

	if len(rq) > 3 || rq.Get(user) == "" || rq.Get(tenderId) == "" || rq.Get(status) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if transport.PutStatusTenderParamsStatus(rq.Get(status)) != transport.PutStatusTenderParamsStatusCREATED && transport.PutStatusTenderParamsStatus(rq.Get(status)) != transport.PutStatusTenderParamsStatusCLOSED && transport.PutStatusTenderParamsStatus(rq.Get(status)) != transport.PutStatusTenderParamsStatusPUBLISHED {
		h.log.Error(model.BadStatus.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tender, err := h.tenderEdite.Status(r.Context(), rq.Get(user), rq.Get(tenderId), rq.Get(status))
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error("FetchTenderStatus error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error("FetchTenderStatus error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error("FetchTenderStatus error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.producer.SendEvent("tender_published", tender)
	if err != nil {
		h.log.Error("Error kafka handler" + err.Error())
	}

	b, err := json.Marshal(convert.TenderModelToTransport(tender))
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
