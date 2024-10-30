package handlers

import (
	"net/http"
)

type Handler interface {
	Register(route *http.ServeMux)
}
