// pkg/gateway/server.go
package gateway

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	h http.Handler
}

func New() *Server {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// 可以加一些基础路由
	r.Get("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return &Server{h: r}
}

func (s *Server) Handler() http.Handler {
	return s.h
}
