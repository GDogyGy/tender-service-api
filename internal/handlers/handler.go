package handlers

import (
	"net/http"
)

// Handler interface TODO: обсудить куда положить, хотел сюда чтобы в др сущностях переиспользовать
type Handler interface {
	Register(route *http.ServeMux)
}
