package xhttpmiddleware

import (
	"net/http"
)

type XHTTPMethodOverrideHandler struct {
	SubHandler http.Handler
	logger     Logger
}

func NewXHTTPMethodOverrideHandler(subHandler http.Handler, logger Logger) (xhttp *XHTTPMethodOverrideHandler) {
	xhttp = &XHTTPMethodOverrideHandler{
		SubHandler: subHandler,
		logger:     logger,
	}
	return
}

func (xhm XHTTPMethodOverrideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subHandler := xhm.SubHandler
	if nil == subHandler {
		if xhm.logger != nil {
			xhm.logger.Error("subhandler is not passed")
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if xHttpHeader := r.Header.Get("X-HTTP-Method-Override"); xHttpHeader != "" {
		r.Method = xHttpHeader
		r.Header.Del("X-HTTP-Method-Override")
	}

	xhm.SubHandler.ServeHTTP(w, r)
}
