package tender

import (
	"context"
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

//go:generate mockery --inpackage --name=service --exported --testonly --inpackage-suffix
type service interface {
	FetchList(ctx context.Context, serviceType string) ([]model.Tender, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Tender, error)
	CheckResponsibleTender(ctx context.Context, username string, tenderId string) (bool, error)
	FetchTenderById(ctx context.Context, tenderId string) (model.Tender, error)
	CreateTender(ctx context.Context, saveModel model.Tender) (model.Tender, error)
}

//go:generate mockery --inpackage --name=organizationResponsibleFacade --exported --testonly --inpackage-suffix
type organizationResponsibleFacade interface {
	Fetch(ctx context.Context, args []byte) (model.OrganizationResponsible, error)
}

type Handler struct {
	log                           log
	service                       service
	organizationResponsibleFacade organizationResponsibleFacade
}

func NewHandler(l log, s service, f organizationResponsibleFacade) Handler {
	return Handler{
		l, s, f,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /api/tenders", h.FetchList)
	router.HandleFunc(http.MethodGet+" /api/tenders/my", h.FetchListByUser)
	router.HandleFunc(http.MethodGet+" /api/tenders/status", h.FetchTenderStatus)
	router.HandleFunc(http.MethodPost+" /api/tenders/new", h.CreateTender)
}

func (h *Handler) FetchList(w http.ResponseWriter, r *http.Request) {
	//log := h.log.With( // TODO: сделать обертку
	//	[]string{"op", "handlers.organizationResponsible.GetList"},
	//)
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

	categories, err := h.service.FetchList(r.Context(), rq.Get(param))

	if err != nil {
		h.log.Error("GetTenderList error: %op" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Возможно ли это протестировать? По идее в
	b, err := json.Marshal(categories)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(categories) > 0 {
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
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

	tenders, err := h.service.FetchListByUser(r.Context(), rq.Get(param))

	if err != nil {
		h.log.Error("FetchListByUser error: %op" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Тот же вопрос с тестами чтобы coverage добить
	b, err := json.Marshal(tenders)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(tenders) > 0 {
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}

// TODO: в тз не особо описана ручка /tenders/status поэтому предположим, что status по id тендера только пользователю ответственному за тендер
func (h *Handler) FetchTenderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rq := r.URL.Query()
	user := "username"
	tenderId := "tenderId"

	if len(rq) >= 2 && rq.Get(user) == "" || rq.Get(tenderId) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.service.CheckResponsibleTender(r.Context(), rq.Get(user), rq.Get(tenderId))

	if errors.Is(err, model.NotFindResponsibleTender) {
		h.log.Error("FetchTenderStatus error: %op" + err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err != nil {
		h.log.Error("FetchTenderStatus error: %op" + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tender, err := h.service.FetchTenderById(r.Context(), rq.Get(tenderId))

	if err != nil {
		h.log.Error("FetchTenderStatus error: %op" + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(tender)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

func (h *Handler) CreateTender(w http.ResponseWriter, r *http.Request) {
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

	if len(b) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: тут прежде нужно проверить user принадлежит ли организации и составить DTO с responsible для tenders и сохранить
	//_, _, err := h.organizationResponsibleFacade.Fetch(r.Context(), b)
	//
	//if err != nil {
	//	h.log.Error(err.Error())
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	organizationResponsible, err := h.organizationResponsibleFacade.Fetch(r.Context(), b)
	if err != nil {
		h.log.Error(err.Error())
		// TODO: http.StatusForbidden потому что у username нет прав создавать тендер от этой органицации
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var tenderDTO model.Tender
	err = json.Unmarshal(b, &tenderDTO)
	tenderDTO.Responsible = organizationResponsible.Id

	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: в тз ответ на роут с маленькой буквы все параметры, возможно нужен какой-то converter?
	tender, err := h.service.CreateTender(r.Context(), tenderDTO)

	if err != nil {
		h.log.Error("CreateTender error: %op" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err = json.Marshal(tender)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}
