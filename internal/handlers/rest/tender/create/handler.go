package create

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"TenderServiceApi/internal/model"
)

type log interface {
	Error(msg string, args ...any)
}

type useCasesTenderCreate interface {
	Create(ctx context.Context, creatorUsername string, organizationId string, saveModel model.Tender) (model.Tender, error)
}

type producer interface {
	SendEvent(eventType string, model interface{}) error
	SendAsync(eventType string, model interface{})
}

type Handler struct {
	log                  log
	producer             producer
	useCasesTenderCreate useCasesTenderCreate
}

func NewHandler(l log, p producer, useCasesTenderCreate useCasesTenderCreate) Handler {
	return Handler{
		l, p, useCasesTenderCreate,
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

	defer func() { _ = r.Body.Close() }()

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tenderRequest transport.TenderCreateRequest
	err = json.Unmarshal(b, &tenderRequest)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if tenderRequest.Status != transport.TenderCreateRequestStatusCREATED {
		h.log.Error(model.BadStatus.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.useCasesTenderCreate.Create(r.Context(), tenderRequest.CreatorUsername, tenderRequest.OrganizationId, convert.TenderReqCreateTransportToModel(tenderRequest))
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

	tender := convert.TenderModelToTransport(resp)

	err = h.producer.SendEvent("tender_created", tender)
	if err != nil {
		h.log.Error("Error kafka handler" + err.Error())
	}

	b, err = json.Marshal(tender)
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
