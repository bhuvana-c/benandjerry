package httputil

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

var httpStatusCodes = map[int]string{
	http.StatusInternalServerError: "internal_server_error",
	http.StatusNotFound:            "not_found",
	http.StatusBadRequest:          "bad_request",
}

type HandlerFunc func(http.ResponseWriter, *http.Request, httprouter.Params) *HandlerError

type HandlerError struct {
	HttpStatusCode int    `json:"HttpStatusCode"`
	Err            string `json:"Error"`
	Message        string `json:"Message"`
}

func NewHandlerError(statusCode int, err error, message string) *HandlerError {
	return &HandlerError{
		HttpStatusCode: statusCode,
		Err:            err.Error(),
		Message:        message,
	}
}

func NewNotFoundError(message string, err error) *HandlerError {
	return NewHandlerError(http.StatusNotFound, err, message)
}
func NewUnexpectedError(message string, err error) *HandlerError {
	return NewHandlerError(http.StatusInternalServerError, err, message)
}
func NewFormatError(message string, err error) *HandlerError {
	return NewHandlerError(http.StatusBadRequest, err, message)
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func ToHandle(h HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}
		handlerErr := h(w, r, p)
		//To Respond with error.
		if handlerErr != nil {
			WriteError(handlerErr, r, w)
		}
	}
}

func WriteError(handlerErr *HandlerError, r *http.Request, w http.ResponseWriter) {
	logFields := map[string]interface{}{
		"httpStatusCode": httpStatusCodes[handlerErr.HttpStatusCode],
		"error":          handlerErr.Err,
		"method":         r.Method,
	}
	log.WithFields(logFields).Error("request failed")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(handlerErr.HttpStatusCode)
	err := json.NewEncoder(w).Encode(handlerErr)
	if err != nil {
		log.Fatal("serializing http error failed: ", err)
	}
}
