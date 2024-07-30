package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VladislavLisovenko/antibruteforce/internal/app"
	"github.com/VladislavLisovenko/antibruteforce/internal/logger"
)

const (
	LoginField    = "login"
	PasswordField = "password"
	IPField       = "ip"
)

var (
	pageNotFoundText    = []byte("Page not found")
	errMethodNotAllowed = errors.New("method not allowed")
)

type Handler struct {
	app    app.Application
	logger logger.Logger
}

type result struct {
	Ok  bool  `json:"ok"`
	Err error `json:"error,omitempty"`
}

func NewHandlers(a app.Application, l logger.Logger) Handler {
	return Handler{app: a, logger: l}
}

func (h *Handler) Handlers(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res result

		switch r.URL.Path {
		case "/check":
			res = h.check(r)
		case "/reset":
			res = h.reset(r)
		case "/whiteList":
			res = h.whiteList(ctx, r)
		case "/blackList":
			res = h.blackList(ctx, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write(pageNotFoundText)
			if err != nil {
				h.logger.Error(err)
			}
			return
		}

		body, err := json.Marshal(res)
		if err != nil {
			h.logger.Error(err)
			res.Err = err
		}

		if res.Err != nil {
			if errors.Is(res.Err, errMethodNotAllowed) {
				w.WriteHeader(http.StatusMethodNotAllowed)
			} else {
				h.logger.Error(res.Err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			if _, err := w.Write([]byte(res.Err.Error())); err != nil {
				h.logger.Error(err)
			}
		}
		if _, err := w.Write(body); err != nil {
			h.logger.Error(err)
		}
	}
}

func (h *Handler) check(r *http.Request) result {
	vs := r.URL.Query()
	ok := h.app.CheckAuth(vs.Get(LoginField), vs.Get(PasswordField), vs.Get(IPField))
	return result{Ok: ok}
}

func (h *Handler) reset(r *http.Request) result {
	vs := r.URL.Query()
	h.app.ResetAuth(vs.Get(LoginField), vs.Get(IPField))
	return result{Ok: true}
}

func (h *Handler) whiteList(ctx context.Context, r *http.Request) result {
	vs := r.URL.Query()
	switch r.Method {
	case http.MethodPost:
		err := h.app.AddToWhiteList(ctx, vs.Get(IPField))
		if err != nil {
			return result{Ok: false, Err: err}
		}
	case http.MethodDelete:
		err := h.app.DeleteFromWhiteList(ctx, vs.Get(IPField))
		if err != nil {
			return result{Ok: false, Err: err}
		}
	default:
		return result{Ok: false, Err: errMethodNotAllowed}
	}

	return result{Ok: true}
}

func (h *Handler) blackList(ctx context.Context, r *http.Request) result {
	vs := r.URL.Query()
	switch r.Method {
	case http.MethodPost:
		err := h.app.AddToBlackList(ctx, vs.Get(IPField))
		if err != nil {
			return result{Ok: false, Err: err}
		}
	case http.MethodDelete:
		err := h.app.DeleteFromBlackList(ctx, vs.Get(IPField))
		if err != nil {
			return result{Ok: false, Err: err}
		}
	default:
		return result{Ok: false, Err: errMethodNotAllowed}
	}

	return result{Ok: true}
}
