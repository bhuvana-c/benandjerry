package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/zalora/benandjerry/httputil"
	"github.com/zalora/benandjerry/model"
)

type IceCreamHandler struct {
	IceCreamStore model.IceCreamStore
}

func (h *IceCreamHandler) AddRoutes(r *httprouter.Router) {
	r.POST("/icecreams/create", httputil.ToHandle(h.Create))
	r.GET("/icecreams/list", httputil.ToHandle(h.List))
	r.GET("/icecreams/show/:name", httputil.ToHandle(h.Get))
	r.POST("/icecreams/update", httputil.ToHandle(h.Update))
	r.DELETE("/icecreams/delete/:name", httputil.ToHandle(h.Delete))

}

func WriteJson(status int, response interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logrus.WithError(err).Error("error while writing the response")
	}
}

func (h *IceCreamHandler) List(w http.ResponseWriter, r *http.Request, p httprouter.Params) *httputil.HandlerError {
	iceCreams, err := h.IceCreamStore.List()
	if err != nil {
		return httputil.NewUnexpectedError("error while listing ice creams", err)
	}
	var response = struct {
		IceCreams []*model.IceCream `json:"IceCreams"`
	}{iceCreams}
	WriteJson(http.StatusOK, response, w)
	return nil
}

func (h *IceCreamHandler) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) *httputil.HandlerError {
	name := p.ByName("name")
	if len(name) == 0 {
		return httputil.NewFormatError("'name' is invalid", errors.New("'name' is invalid"))
	}
	iceCream, err := h.IceCreamStore.Get(name)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); ok {
			return httputil.NewNotFoundError(fmt.Sprintf("'%s' not found", name), errors.New("not_found"))
		}
		return httputil.NewUnexpectedError("error while getting ice cream", err)
	}
	WriteJson(http.StatusOK, iceCream, w)
	return nil
}

func (h *IceCreamHandler) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) *httputil.HandlerError {
	iceCream := &model.IceCream{}
	err := json.NewDecoder(r.Body).Decode(iceCream)
	if err != nil {
		return httputil.NewFormatError("new ice cream is invalid", errors.New("new ice cream is invalid"))
	}
	err = h.IceCreamStore.Update(iceCream)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); ok {
			return httputil.NewNotFoundError(fmt.Sprintf("'%s' not found", iceCream.Name), errors.New("not_found"))
		}
		return httputil.NewUnexpectedError("error while getting ice cream", err)
	}
	WriteJson(http.StatusOK, "{}", w)
	return nil
}

func (h *IceCreamHandler) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) *httputil.HandlerError {
	newIceCream := &model.IceCream{}
	err := json.NewDecoder(r.Body).Decode(&newIceCream)
	if err != nil {
		return httputil.NewFormatError("body is invalid", err)
	}
	err = h.IceCreamStore.Create(newIceCream)
	if err != nil {
		return httputil.NewUnexpectedError("error while adding a new ice cream", err)
	}
	response := struct {
		Name string `json:"name"`
	}{newIceCream.Name}
	WriteJson(http.StatusOK, response, w)
	return nil
}
func (h *IceCreamHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) *httputil.HandlerError {
	name := p.ByName("name")
	if len(name) == 0 {
		return httputil.NewFormatError("'name' is invalid", errors.New("'name' is invalid"))
	}
	err := h.IceCreamStore.Delete(name)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); ok {
			return httputil.NewNotFoundError(fmt.Sprintf("'%s' not found", name), errors.New("not_found"))
		}
		return httputil.NewUnexpectedError("error while deleting ice creams", err)
	}
	WriteJson(http.StatusOK, "{}", w)
	return nil
}
