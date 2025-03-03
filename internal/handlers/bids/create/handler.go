package create

import (
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
	defer r.Body.Close()

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var args argCreatBids
	err = json.Unmarshal(b, &args)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var bidsCreate model.Bids
	err = json.Unmarshal(b, &bidsCreate)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.useCaseBidsCreate.Create(r.Context(), args.Username, args.OrganizationId, bidsCreate)
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

	tDTO := bidsDTO{resp.Id, resp.Name, resp.Description, resp.Status, resp.TenderId, resp.Version, resp.Responsible}
	b, err = json.Marshal(tDTO)
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

	var bidFeedbackCreate model.BidFeedback
	err = json.Unmarshal(b, &bidFeedbackCreate)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.useCaseBidFeedbackCreate.Create(r.Context(), rq.Get(username), bidID[1], bidFeedbackCreate)
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

	tDTO := bidsFeedbackDTO{resp.Id, resp.BidID, resp.Description, resp.Responsible, resp.CreatedAt}
	b, err = json.Marshal(tDTO)
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
