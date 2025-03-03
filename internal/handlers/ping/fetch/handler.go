package fetch

import (
	"net/http"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /api/ping", h.Ping)
}

// TODO: Обратите внимание, что успешное выполнение запроса GET /api/ping обязательно для начала тестирования приложения.
// решил с помощью makefile может есть другое решение?
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}
