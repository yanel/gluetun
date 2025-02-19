package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/updater"
)

func newUpdaterHandler(
	ctx context.Context,
	looper updater.Looper,
	warner warner) http.Handler {
	return &updaterHandler{
		ctx:    ctx,
		looper: looper,
		warner: warner,
	}
}

type updaterHandler struct {
	ctx    context.Context //nolint:containedctx
	looper updater.Looper
	warner warner
}

func (h *updaterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.TrimPrefix(r.RequestURI, "/updater")
	switch r.RequestURI {
	case "/status":
		switch r.Method {
		case http.MethodGet:
			h.getStatus(w)
		case http.MethodPut:
			h.setStatus(w, r)
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func (h *updaterHandler) getStatus(w http.ResponseWriter) {
	status := h.looper.GetStatus()
	encoder := json.NewEncoder(w)
	data := statusWrapper{Status: string(status)}
	if err := encoder.Encode(data); err != nil {
		h.warner.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *updaterHandler) setStatus(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data statusWrapper
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status, err := data.getStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	outcome, err := h.looper.SetStatus(h.ctx, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(outcomeWrapper{Outcome: outcome}); err != nil {
		h.warner.Warn(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
