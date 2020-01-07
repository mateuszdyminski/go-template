package api

import (
	"encoding/json"
	"net/http"

	_ "github.com/mateuszdyminski/go-template/swagger-docs"
	"github.com/swaggo/swag"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// WriteErrJSON wraps error in JSON structure.
func WriteErrJSON(l *zap.SugaredLogger, w http.ResponseWriter, r *http.Request, err error, httpCode int) {
	// log outgoing errors
	l.With("requestId", GetReqID(r.Context())).Error(err)

	// write error to response
	e := HTTPError{
		HTTPStatusCode:  httpCode,
		Msg:             err.Error(),
		InternalErrCode: -1,
	}

	if err := WriteJSON(w, e, httpCode); err != nil {
		l.Errorw("error while sending err json", "err", err)
	}
}

// WriteJSON writes response to client, response is a struct defining JSON reply.
func WriteJSON(w http.ResponseWriter, data interface{}, httpCode int) error {
	json, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "can't encode JSON")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)

	if _, err := w.Write(json); err != nil {
		return errors.Wrap(err, "can't write bytes to response writer")
	}

	return nil
}

// MustWriteJSON writes response to client, response is a struct defining JSON reply.
func MustWriteJSON(l *zap.SugaredLogger, w http.ResponseWriter, r *http.Request, data interface{}, httpCode int) {
	if err := WriteJSON(w, data, httpCode); err != nil {
		WriteErrJSON(l, w, r, err, http.StatusInternalServerError)
	}
}

func SwaggerHandler(l *zap.SugaredLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if err != nil {
			WriteErrJSON(l, w, r, err, http.StatusInternalServerError)
		}

		if _, err := w.Write([]byte(doc)); err != nil {
			l.Errorw("error while sending err json", "err", err)
		}
	}
}

// HTTPError - general error response for api.
type HTTPError struct {
	HTTPStatusCode  int    `json:"httpStatusCode"`
	Msg             string `json:"msg"`
	InternalErrCode int    `json:"internalErrCode"`
}
