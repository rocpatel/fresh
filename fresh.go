package fresh

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/julienschmidt/httprouter"
)

type ErrorHandler func(error, *Context) error

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	ctx      context.Context
	params   httprouter.Params
}

func newContext(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	return &Context{
		Response: w,
		Request:  r,
		ctx:      context.Background(),
		params:   params,
	}
}

func (c *Context) Render(component templ.Component) error {
	return component.Render(c.ctx, c.Response)
}

func (c *Context) Set(key string, value any) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

func (c *Context) Get(key string) any {
	return c.ctx.Value(key)
}

func (c *Context) Param(name string) string {
	return c.params.ByName(name)
}

func (c *Context) Query(name string) string {
	return c.Request.URL.Query().Get(name)
}

func (c *Context) FormValue(name string) string {
	return c.Request.FormValue(name)
}

func (c *Context) Redirect(url string, code int) error {
	return nil
}

func (c *Context) JSON(status int, v any) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	return json.NewDecoder(c.Request.Body).Decode(&v)
}

type Plug func(Handler) Handler

type Handler func(c *Context) error

type Fresh struct {
	ErrorHandler ErrorHandler
	router       *httprouter.Router
	plugs        []Plug
}

func New() *Fresh {
	return &Fresh{
		router:       httprouter.New(),
		ErrorHandler: defaultErrorHandler,
	}
}

func (s *Fresh) Plug(plugs ...Plug) {
	s.plugs = append(s.plugs, plugs...)
}

func (s *Fresh) Start(port string) error {
	return http.ListenAndServe(port, s.router)
}

func (s *Fresh) add(method, path string, h Handler, plugs ...Plug) {
	s.router.Handle(method, path, s.makeHTTPRouterHandler(h, plugs...))
}

func (s *Fresh) Get(path string, h Handler, plugs ...Plug) {
	s.add("GET", path, h, plugs...)
}

func (s *Fresh) Post(path string, h Handler, plugs ...Plug) {
	s.add("POST", path, h, plugs...)
}

func (s *Fresh) Put(path string, h Handler, plugs ...Plug) {
	s.add("PUT", path, h, plugs...)
}

func (s *Fresh) Head(path string, h Handler, plugs ...Plug) {
	s.add("HEAD", path, h, plugs...)
}

func (s *Fresh) Options(path string, h Handler, plugs ...Plug) {
	s.add("OPTIONS", path, h, plugs...)
}

func (s *Fresh) makeHTTPRouterHandler(h Handler, plugs ...Plug) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := newContext(w, r, params)

		for i := len(s.plugs) - 1; i >= 0; i-- {
			h = s.plugs[i](h)
		}

		for i := len(plugs) - 1; i >= 0; i-- {
			h = plugs[i](h)
		}

		if err := h(ctx); err != nil {
			// todo: handle the error from teh error handler
			s.ErrorHandler(err, ctx)
		}
	}
}

func defaultErrorHandler(err error, _ *Context) error {
	slog.Error("error", "err", err)
	return nil
}
