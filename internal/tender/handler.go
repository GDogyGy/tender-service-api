package tender

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"TenderServiceApi/internal/handlers"
	"TenderServiceApi/internal/storage/postgres"
)

var _ handlers.Handler = &handler{}

// handler struct TODO: вот тут убери зависимости в интерфейсы пожалуйста и напиши тесты потом
type handler struct {
	log     *slog.Logger
	storage *postgres.Storage
}

func NewHandler(log *slog.Logger, db *postgres.Storage) handlers.Handler {
	return &handler{
		log:     log,
		storage: db,
	}
}

func (h *handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /api/tender", h.GetList)
	router.HandleFunc(http.MethodGet+" /api/tender/my", h.GetList)
	router.HandleFunc(http.MethodGet+" /api/tender/new", h.Create)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.tender.GetList"
	log := h.log.With(
		slog.String("op", op),
	)
	rq := r.URL.Query()
	param := "servicetype"

	if len(rq) > 1 || len(rq) == 1 && rq.Get(param) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//categories, err := GetTenderList(h.storage.Db, rq.Get(param))
	services := tenderServices{h.log, h.storage.Db}
	categories, err := services.GetTenderList(rq.Get(param))

	if err != nil {
		log.Error("GetTenderList error: %op" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(categories)
	if err != nil {
		log.Error(err.Error())
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

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint(r.URL.Query())))
	w.WriteHeader(http.StatusCreated)
}
