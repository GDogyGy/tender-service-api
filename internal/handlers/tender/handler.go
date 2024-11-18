package tender

import (
	"encoding/json"
	"net/http"

	"TenderServiceApi/internal/model"
)

type log interface {
	Error(msg string, args ...any)
	//With(args ...any) log // TODO тут не получилось сделать не импортируемым: обсудить. При *log не считается что удовлетворяет With(args ...any) *Logger
}

type service interface {
	GetTenderList(serviceType string) ([]model.Tender, error)
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
	router.HandleFunc(http.MethodGet+" /api/tender", h.GetList)
	router.HandleFunc(http.MethodGet+" /api/tender/my", h.GetList)
	router.HandleFunc(http.MethodGet+" /api/tender/new", h.Create)
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	//log := h.log.With( // TODO: обсудить как сделать или как в таких случаях делают
	//	[]string{"op", "handlers.tender.GetList"},
	//)
	rq := r.URL.Query()

	param := "servicetype"

	if len(rq) > 1 || len(rq) == 1 && rq.Get(param) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	categories, err := h.service.GetTenderList(rq.Get(param))

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
