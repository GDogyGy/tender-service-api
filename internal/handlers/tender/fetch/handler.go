package fetch

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"TenderServiceApi/internal/handlers/types/convert"
	"TenderServiceApi/internal/handlers/types/transport"
	"TenderServiceApi/internal/model"
)

//go:generate mockery  --inpackage --name=log --exported --testonly --inpackage-suffix
type log interface {
	Error(msg string, args ...any)
}

//go:generate mockery --inpackage --name=useCasesTenderFetch --exported --testonly --inpackage-suffix
type useCasesTenderFetch interface {
	FetchList(ctx context.Context, serviceType string) ([]model.Tender, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Tender, error)
	FetchStatus(ctx context.Context, username string, tenderId string) (model.Tender, error)
}

type Handler struct {
	log         log
	tenderFetch useCasesTenderFetch
}

func NewHandler(l log, t useCasesTenderFetch) Handler {
	return Handler{
		l, t,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /api/tenders", h.FetchList)
	router.HandleFunc(http.MethodGet+" /api/tenders/my", h.FetchListByUser)
	router.HandleFunc(http.MethodGet+" /api/tenders/status", h.FetchStatus)
}

func (h *Handler) FetchList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rq := r.URL.Query()

	param := "servicetype"

	if len(rq) >= 1 && rq.Get(param) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	categories, err := h.tenderFetch.FetchList(r.Context(), rq.Get(param))
	if err != nil {
		h.log.Error("GetTenderList error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tendersTransport []transport.Tender
	for _, v := range categories {
		tendersTransport = append(tendersTransport, convert.TenderModelToTransport(v))
	}

	b, err := json.Marshal(tendersTransport)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(tendersTransport) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) FetchListByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rq := r.URL.Query()

	param := "username"
	if len(rq) >= 1 && rq.Get(param) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tenders, err := h.tenderFetch.FetchListByUser(r.Context(), rq.Get(param))
	if err != nil {
		h.log.Error("FetchListByUser error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tendersTransport []transport.Tender
	for _, v := range tenders {
		tendersTransport = append(tendersTransport, convert.TenderModelToTransport(v))
	}

	b, err := json.Marshal(tendersTransport)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(tendersTransport) == 0 {
		w.WriteHeader(http.StatusNotFound)
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

func (h *Handler) FetchStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rq := r.URL.Query()
	user := "username"
	tenderId := "tenderId"

	if len(rq) > 2 || rq.Get(user) == "" || rq.Get(tenderId) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tender, err := h.tenderFetch.FetchStatus(r.Context(), rq.Get(user), rq.Get(tenderId))
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

	b, err := json.Marshal(convert.TenderModelToTransport(tender).Status)
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
