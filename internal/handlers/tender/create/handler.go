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

//go:generate mockery  --inpackage --name=log --exported --testonly --inpackage-suffix
type log interface {
	Error(msg string, args ...any)
}

//go:generate mockery --inpackage --name=useCasesTenderCreate --exported --testonly --inpackage-suffix
type useCasesTenderCreate interface {
	Create(ctx context.Context, creatorUsername string, organizationId string, saveModel model.Tender) (model.Tender, error)
}

type Handler struct {
	log                  log
	useCasesTenderCreate useCasesTenderCreate
}

func NewHandler(l log, useCasesTenderCreate useCasesTenderCreate) Handler {
	return Handler{
		l, useCasesTenderCreate,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodPost+" /api/tenders/new", h.Create)
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

	var args argCreatTender
	err = json.Unmarshal(b, &args)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tenderCreate model.Tender
	err = json.Unmarshal(b, &tenderCreate)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.useCasesTenderCreate.Create(r.Context(), args.Username, args.OrganizationId, tenderCreate)
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

	tDTO := tenderDTO{resp.Id, resp.Name, resp.Description, resp.ServiceType, resp.Status, resp.Version, resp.Responsible}
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
