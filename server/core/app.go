package core

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type App struct {
	Mux         *http.ServeMux
	middlewares []Middleware
}

func NewApp() *App {
	return &App{
		Mux:         http.NewServeMux(),
		middlewares: make([]Middleware, 0),
	}
}

func (s *App) Use(mw ...Middleware) {
	s.middlewares = append(s.middlewares, mw...)
}

// wraps handlers with middlewares
func (s *App) HandleFunc(pattern string, handler http.HandlerFunc) {
	lastMiddleware := len(s.middlewares) - 1

	if lastMiddleware == -1 {
		s.Mux.Handle(pattern, handler)
		return
	}

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		md := s.middlewares[i]
		// encapsulating handlers
		if h, ok := md(handler).(http.HandlerFunc); ok {
			handler = h
		} else {
			panic("middleware should return http.HandlerFunc")
		}
	}

	s.Mux.HandleFunc(pattern, handler)
}
