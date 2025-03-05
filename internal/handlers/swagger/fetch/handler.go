package fetch

import (
	"net/http"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

// TODO: Как отобразить без swagger-ui или как пакеты убрать в external library
func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.yaml")
	})
	router.Handle(http.MethodGet+" /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("swagger-ui/dist"))))
}
