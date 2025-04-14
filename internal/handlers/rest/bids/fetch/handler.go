package fetch

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/handlers/rest/types/transport"
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
	FetchStatus(ctx context.Context, username string, bidsId string) (model.Bids, error)
}

//go:generate mockery --inpackage --name=useCaseBidFeedbackFetch --exported --testonly --inpackage-suffix
type useCaseBidFeedbackFetch interface {
	FetchReviews(ctx context.Context, username string, tenderID string, authorUsername string, organizationID string) ([]model.BidFeedback, error)
}

type Handler struct {
	log              log
	bidsFetch        useCaseBidsFetch
	bidFeedbackFetch useCaseBidFeedbackFetch
}

func NewHandler(l log, t useCaseBidsFetch, f useCaseBidFeedbackFetch) Handler {
	return Handler{
		l, t, f,
	}
}

var tenderIdReviewRegexp = regexp.MustCompile(`/api/bids/(.*)/reviews\?`)
var tenderIdRegexp = regexp.MustCompile(`/api/bids/(.*)/list`)
var bidIdStatusRegexp = regexp.MustCompile(`/api/bids/(.*)/status\?`)

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /api/bids/{tenderID}/list", h.FetchListByTender)
	router.HandleFunc(http.MethodGet+" /api/bids/my", h.FetchListByUser)
	router.HandleFunc(http.MethodGet+" /api/bids/{bidID}/status", h.FetchStatus)
	router.HandleFunc(http.MethodGet+" /api/bids/{tenderID}/reviews", h.FetchReviews)
}

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

	var bidsTransport []transport.Bid
	for _, v := range bids {
		bidsTransport = append(bidsTransport, convert.BidsModelToTransport(v))
	}

	b, err := json.Marshal(bidsTransport)
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

	bids, err := h.bidsFetch.FetchListByUser(r.Context(), rq.Get(param))
	if err != nil {
		h.log.Error("FetchListByUser error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var bidsTransport []transport.Bid
	for _, v := range bids {
		bidsTransport = append(bidsTransport, convert.BidsModelToTransport(v))
	}

	b, err := json.Marshal(bidsTransport)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(bidsTransport) == 0 {
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

	if len(rq) > 1 || rq.Get(user) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bidID := bidIdStatusRegexp.FindStringSubmatch(r.RequestURI)
	if bidID == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bid, err := h.bidsFetch.FetchStatus(r.Context(), rq.Get(user), bidID[1])
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error("FetchBidsStatus error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error("FetchBidsStatus error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error("FetchBidsStatus error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(convert.BidsModelToTransport(bid).Status)
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

func (h *Handler) FetchReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rq := r.URL.Query()
	username := "username"
	organizationID := "organizationId"
	authorUser := "authorUsername"
	if len(rq) < 3 || rq.Get(username) == "" || rq.Get(organizationID) == "" || rq.Get(authorUser) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tenderID := tenderIdReviewRegexp.FindStringSubmatch(r.RequestURI)
	if tenderID == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	feedback, err := h.bidFeedbackFetch.FetchReviews(r.Context(), rq.Get(username), tenderID[1], rq.Get(authorUser), rq.Get(organizationID))
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error("FetchReviews error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error("FetchReviews error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error("FetchReviews error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var feedbackTransport []transport.Feedback
	for _, v := range feedback {
		feedbackTransport = append(feedbackTransport, convert.BidFeedbackModelToTransport(v))
	}

	b, err := json.Marshal(feedbackTransport)
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
