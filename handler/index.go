package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IndexHandler struct{}

func (h *IndexHandler) AddRoutes(r *httprouter.Router) {
	r.GET("/", h.Home)
}

func (h *IndexHandler) Home(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write([]byte(""))
}
