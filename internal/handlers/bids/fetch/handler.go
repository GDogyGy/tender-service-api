package fetch

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"TenderServiceApi/internal/model"
)

//go:generate mockery  --inpackage --name=log --exported --testonly --inpackage-suffix
type log interface {
	Error(msg string, args ...any)
}

//go:generate mockery --inpackage --name=useCaseBidsFetch --exported --testonly --inpackage-suffix
type useCaseBidsFetch interface {
	FetchListByTender(ctx context.Context, username string, tenderId string) ([]model.Bids, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Bids, error)
}

type Handler struct {
	log       log
	bidsFetch useCaseBidsFetch
}

func NewHandler(l log, t useCaseBidsFetch) Handler {
	return Handler{
		l, t,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /api/bids/{tenderId}/list", h.FetchListByTender)
	router.HandleFunc(http.MethodGet+" /api/bids/my", h.FetchListByUser)
}

var tenderIdRegexp = regexp.MustCompile(`/api/bids/(.*)/list`)

func (h *Handler) FetchListByTender(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rq := r.URL.Query()
	user := "username"

	if rq.Get(user) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := tenderIdRegexp.FindStringSubmatch(r.RequestURI)
	if id == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bids, err := h.bidsFetch.FetchListByTender(r.Context(), rq.Get(user), id[1])
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error("FetchListByTender error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, model.NotFound) {
		h.log.Error("Bids FetchListByTender error: " + err.Error())
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		h.log.Error("FetchListByTender error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(bids) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	b, err := json.Marshal(bids)
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

	tenders, err := h.bidsFetch.FetchListByUser(r.Context(), rq.Get(param))
	if err != nil {
		h.log.Error("FetchListByUser error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(tenders)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(tenders) == 0 {
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
