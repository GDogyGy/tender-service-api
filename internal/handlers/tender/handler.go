package tender

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"TenderServiceApi/internal/model"
)

//go:generate mockery  --inpackage --name=log --exported --testonly --inpackage-suffix
type log interface {
	Error(msg string, args ...any)
}

//go:generate mockery --inpackage --name=serviceTender --exported --testonly --inpackage-suffix
type serviceTender interface {
	FetchList(ctx context.Context, serviceType string) ([]model.Tender, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Tender, error)
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
	FetchById(ctx context.Context, tenderId string) (model.Tender, error)
	Create(ctx context.Context, saveModel model.Tender) (model.Tender, error)
	Edite(ctx context.Context, tenderNew model.Tender, tender model.Tender) (model.Tender, error)
}

//go:generate mockery --inpackage --name=serviceOrganization --exported --testonly --inpackage-suffix
type serviceOrganization interface {
	FetchById(ctx context.Context, id string) (model.Organization, error)
	CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error)
	FetchRelationsById(ctx context.Context, id string) (model.OrganizationResponsible, error)
}

type Handler struct {
	log          log
	tender       serviceTender
	organization serviceOrganization
}

func NewHandler(l log, t serviceTender, o serviceOrganization) Handler {
	return Handler{
		l, t, o,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /api/tenders", h.FetchList)
	router.HandleFunc(http.MethodGet+" /api/tenders/my", h.FetchListByUser)
	router.HandleFunc(http.MethodGet+" /api/tenders/status", h.FetchStatus)
	router.HandleFunc(http.MethodPost+" /api/tenders/new", h.Create)
	router.HandleFunc(http.MethodPatch+" /api/tenders/{id}/edit", h.Edite)
	router.HandleFunc(http.MethodPut+" /api/tenders/{id}/rollback/{version}", h.Rollback)
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

	categories, err := h.tender.FetchList(r.Context(), rq.Get(param))

	if err != nil {
		h.log.Error("GetTenderList error: %op" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	tenders, err := h.tender.FetchListByUser(r.Context(), rq.Get(param))

	if err != nil {
		h.log.Error("FetchListByUser error: %op" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	_, err := h.tender.CheckResponsible(r.Context(), rq.Get(user), rq.Get(tenderId))

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
	tender, err := h.tender.FetchById(r.Context(), rq.Get(tenderId))

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

	if len(b) <= 0 {
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
	organizationResponsible, err := h.organization.CheckResponsible(r.Context(), args.Username, args.OrganizationId)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var tender model.Tender
	err = json.Unmarshal(b, &tender)

	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tender.Responsible = organizationResponsible.Id

	resp, err := h.tender.Create(r.Context(), tender)

	if err != nil {
		h.log.Error("CreateTender error: %op" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tDTO := tenderDTO{resp.Id, resp.Name, resp.Description, resp.ServiceType, resp.Status, resp.Responsible}
	b, err = json.Marshal(tDTO)

	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

func (h *Handler) Edite(w http.ResponseWriter, r *http.Request) {
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

	if len(b) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	re, err := regexp.Compile("/tenders/(.*)/edit")
	id := re.FindStringSubmatch(r.RequestURI)[1]
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tender, err := h.tender.FetchById(r.Context(), id)

	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	organization, err := h.organization.FetchRelationsById(r.Context(), tender.Responsible)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = h.organization.CheckResponsible(r.Context(), rq.Get(user), organization.OrganizationId)

	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var tenderNew model.Tender

	tenderNew = tender
	err = json.Unmarshal(b, &tenderNew)

	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.tender.Edite(r.Context(), tenderNew, tender)

	if err != nil {
		h.log.Error("Edite Tender error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tDTO := tenderDTO{resp.Id, resp.Name, resp.Description, resp.ServiceType, resp.Status, resp.Responsible}
	b, err = json.Marshal(tDTO)

	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

func (h *Handler) Rollback(w http.ResponseWriter, r *http.Request) {
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

	_, err := io.ReadAll(r.Body)

	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
