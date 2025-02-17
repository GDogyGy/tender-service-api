package create

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"TenderServiceApi/internal/model"
)

type log interface {
	Error(msg string, args ...any)
}

type useCaseBidsCreate interface {
	Create(ctx context.Context, creatorUsername string, organizationId string, saveModel model.Bids) (model.Bids, error)
}

type Handler struct {
	log               log
	useCaseBidsCreate useCaseBidsCreate
}

func NewHandler(l log, useCaseBidsCreate useCaseBidsCreate) Handler {
	return Handler{
		l, useCaseBidsCreate,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodPost+" /api/bids/new", h.Create)
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
