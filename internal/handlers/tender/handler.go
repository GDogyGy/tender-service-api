package tender

import (
	"context"
	"encoding/json"
	"net/http"

	"TenderServiceApi/internal/model"
)

// TODO: генерирует не экспортируемый мок (приходится переписывать)
//
//go:generate mockery --dir=./ --name=log
type log interface {
	Error(msg string, args ...any)
}

// TODO: генерирует не экспортируемый мок (приходится переписывать) Спросить, как тут быть
//
//go:generate mockery --dir=./ --name=service
type service interface {
	FetchList(ctx context.Context, serviceType string) ([]model.Tender, error)
}

type Handler struct {
	log     log
	service service
}

func NewHandler(l log, s service) Handler {
	return Handler{
		l, s,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /api/tenders", h.FetchList)
	router.HandleFunc(http.MethodGet+" /api/tenders/my", h.FetchList)
	router.HandleFunc(http.MethodGet+" /api/tenders/new", h.Create)
}

func (h *Handler) FetchList(w http.ResponseWriter, r *http.Request) {
	//log := h.log.With( // TODO: обсудить как сделать или как в таких случаях делают
	//	[]string{"op", "handlers.tender.GetList"},
	//)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rq := r.URL.Query()

	param := "servicetype"

	if len(rq) > 1 || len(rq) == 1 && rq.Get(param) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	categories, err := h.service.FetchList(r.Context(), rq.Get(param))

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte(fmt.Sprint(r.URL.Query())))
	//w.WriteHeader(http.StatusCreated)
}
