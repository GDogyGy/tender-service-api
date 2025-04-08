package create

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"TenderServiceApi/internal/handlers/types/convert"
	"TenderServiceApi/internal/handlers/types/transport"
	"TenderServiceApi/internal/model"
)

type log interface {
	Error(msg string, args ...any)
}

type useCaseBidsCreate interface {
	Create(ctx context.Context, creatorUsername string, organizationId string, saveModel model.Bids) (model.Bids, error)
}

type useCaseBidFeedbackCreate interface {
	Create(ctx context.Context, creatorUsername string, tenderID string, saveModel model.BidFeedback) (model.BidFeedback, error)
}

type Handler struct {
	log                      log
	useCaseBidsCreate        useCaseBidsCreate
	useCaseBidFeedbackCreate useCaseBidFeedbackCreate
}

func NewHandler(l log, useCaseBidsCreate useCaseBidsCreate, useCaseBidFeedbackCreate useCaseBidFeedbackCreate) Handler {
	return Handler{
		l, useCaseBidsCreate, useCaseBidFeedbackCreate,
	}
}

var bidFeedbackID = regexp.MustCompile(`/api/bids/(.*)/feedback\?`)

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodPost+" /api/bids/new", h.Create)
	router.HandleFunc(http.MethodPut+" /api/bids/{bidId}/feedback", h.Feedback)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		r.Body.Close()
	}()

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var bidCreateRequest transport.BidCreateRequest
	err = json.Unmarshal(b, &bidCreateRequest)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if bidCreateRequest.Status != transport.BidCreateRequestStatusCREATED {
		h.log.Error(model.BadStatus.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.useCaseBidsCreate.Create(r.Context(), bidCreateRequest.CreatorUsername, bidCreateRequest.OrganizationId, convert.BidsReqCreateTransportToModel(bidCreateRequest))
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bid := convert.BidsModelToTransport(resp)
	b, err = json.Marshal(bid)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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

func (h *Handler) Feedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rq := r.URL.Query()
	username := "username"

	if rq.Get(username) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bidID := bidFeedbackID.FindStringSubmatch(r.RequestURI)
	if bidID == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var feedback transport.Feedback
	err = json.Unmarshal(b, &feedback)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.useCaseBidFeedbackCreate.Create(r.Context(), rq.Get(username), bidID[1], convert.BidFeedbackTransportToModel(feedback))
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

	feedback = convert.BidFeedbackModelToTransport(resp)
	b, err = json.Marshal(feedback)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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
