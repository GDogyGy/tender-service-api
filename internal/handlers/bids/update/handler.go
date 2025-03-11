package update

import (
	"TenderServiceApi/internal/handlers/types/convert"
	"TenderServiceApi/internal/handlers/types/transport"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"TenderServiceApi/internal/model"
)

type log interface {
	Error(msg string, args ...any)
}

type useCaseBidEdit interface {
	Edit(ctx context.Context, id string, username string, bidNew model.Bids) (model.Bids, error)
	Rollback(ctx context.Context, id string, username string, version string) (model.Bids, error)
	Status(ctx context.Context, username string, bidID string, status string) (model.Bids, error)
}

type useCaseBidDecision interface {
	SubmitDecision(ctx context.Context, username string, bidID string, decision string, organizationID string) (model.Bids, error)
}

type Handler struct {
	log          log
	bidsEdite    useCaseBidEdit
	bidsDecision useCaseBidDecision
}

var bidEditRegexp = regexp.MustCompile(`/bids/(.*)/edit`)
var bidIdRollbackRegexp = regexp.MustCompile(`/bids/(.*)/rollback/`)
var bidVersionRollbackRegexp = regexp.MustCompile(`/rollback/(.*)\?`)
var bidIdStatusRegexp = regexp.MustCompile(`/bids/(.*)/status\?`)
var bidEditDecisionRegexp = regexp.MustCompile(`/api/bids/(.*)/submit_decision\?`)

func NewHandler(l log, b useCaseBidEdit, d useCaseBidDecision) Handler {
	return Handler{
		l, b, d,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodPatch+" /api/bids/{bidId}/edit", h.Edit)
	router.HandleFunc(http.MethodPut+" /api/bids/{bidId}/rollback/{version}", h.Rollback)
	router.HandleFunc(http.MethodPut+" /api/bids/{bidId}/status", h.Status)
	router.HandleFunc(http.MethodPut+" /api/bids/{bidId}/submit_decision", h.SubmitDecision)
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

	id := bidEditRegexp.FindStringSubmatch(r.RequestURI)
	if id == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var bid transport.Bid
	err = json.Unmarshal(b, &bid)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := h.bidsEdite.Edit(r.Context(), id[1], rq.Get(user), convert.BidsTransportToModel(bid))
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bid = convert.BidsModelToTransport(resp)
	b, err = json.Marshal(bid)
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

	id := bidIdRollbackRegexp.FindStringSubmatch(r.RequestURI)
	if id == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	version := bidVersionRollbackRegexp.FindStringSubmatch(r.RequestURI)
	if version == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bids, err := h.bidsEdite.Rollback(r.Context(), id[1], rq.Get(user), version[1])
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, model.NotFindResponsible) || errors.Is(err, model.AlreadyVotedResponsible) {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(convert.BidsModelToTransport(bids))
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
	status := "status"

	if len(rq) > 3 || rq.Get(user) == "" || rq.Get(status) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bidID := bidIdStatusRegexp.FindStringSubmatch(r.RequestURI)
	if bidID == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bid, err := h.bidsEdite.Status(r.Context(), rq.Get(user), bidID[1], rq.Get(status))
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error("FetchBidStatus error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error("FetchBidStatus error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error("FetchBidsStatus error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(convert.BidsModelToTransport(bid))
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

func (h *Handler) SubmitDecision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rq := r.URL.Query()
	user := "username"
	decision := "decision"
	organizationID := "organizationId"

	if len(rq) > 3 || rq.Get(user) == "" || rq.Get(decision) == "" || rq.Get(organizationID) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := map[string]bool{
		"REJECTED": false,
		"APPROVED": true,
	}
	_, ok := m[rq.Get(decision)]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bidID := bidEditDecisionRegexp.FindStringSubmatch(r.RequestURI)
	if bidID == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bid, err := h.bidsDecision.SubmitDecision(r.Context(), rq.Get(user), bidID[1], rq.Get(decision), rq.Get(organizationID))
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error("SubmitDecision error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error("SubmitDecision error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error("SubmitDecision error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(convert.BidsModelToTransport(bid))
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
