package update

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

//go:generate mockery  --inpackage --name=log --exported --testonly --inpackage-suffix
type log interface {
	Error(msg string, args ...any)
}

//go:generate mockery --inpackage --name=useCaseTenderEdite --exported --testonly --inpackage-suffix
type useCaseTenderEdite interface {
	Edite(ctx context.Context, id string, username string, tenderNew model.Tender) (model.Tender, error)
	Rollback(ctx context.Context, id string, username string, version string) (model.Tender, error)
}

type Handler struct {
	log         log
	tenderEdite useCaseTenderEdite
}

var tenderEditeRegexp = regexp.MustCompile(`/tenders/(.*)/edit`)
var tenderIdRollbackRegexp = regexp.MustCompile(`/tenders/(.*)/rollback/`)
var tenderVersionRollbackRegexp = regexp.MustCompile(`/rollback/(.*)\?`)

func NewHandler(l log, t useCaseTenderEdite) Handler {
	return Handler{
		l, t,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodPatch+" /api/tenders/{id}/edit", h.Edite)
	router.HandleFunc(http.MethodPut+" /api/tenders/{id}/rollback/{version}", h.Rollback)
}

// Edite TODO: Почему в две базы кладем? потому что Тендеры могут создавать только пользователи от имени своей организации. А этих тендеров может быть несколько от одного человека и как понять какой тендер откатывать а какой не трогать?
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

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := tenderEditeRegexp.FindStringSubmatch(r.RequestURI)
	if id == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tenderNew model.Tender
	err = json.Unmarshal(b, &tenderNew)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := h.tenderEdite.Edite(r.Context(), id[1], rq.Get(user), tenderNew)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Rollback TODO: После отката, считается новой правкой с увеличением версии.
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

	id := tenderIdRollbackRegexp.FindStringSubmatch(r.RequestURI)
	if id == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	version := tenderVersionRollbackRegexp.FindStringSubmatch(r.RequestURI)
	if version == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tender, err := h.tenderEdite.Rollback(r.Context(), id[1], rq.Get(user), version[1])
	// TODO: Валидно так? errors.Is(err, sql.ErrNoRows) || errors.Is(err, model.NotFindResponsibleTender)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, model.NotFindResponsibleTender) {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tDTO := tenderDTO{tender.Id, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.Version, tender.Responsible}
	b, err := json.Marshal(tDTO)

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
