package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/igor-trentini/diffable/backend/internal/handler"
	"github.com/igor-trentini/diffable/backend/internal/middleware"
	"github.com/igor-trentini/diffable/backend/internal/service"
)

type Server struct {
	Router          *chi.Mux
	HealthHandler   *handler.HealthHandler
	AnalysisHandler *handler.AnalysisHandler
}

func New(db handler.DBPinger, frontendURL string, analysisSvc service.AnalysisService) *Server {
	s := &Server{
		Router:          chi.NewRouter(),
		HealthHandler:   handler.NewHealthHandler(db),
		AnalysisHandler: handler.NewAnalysisHandler(analysisSvc),
	}

	s.Router.Use(chimw.RequestID)
	s.Router.Use(chimw.RealIP)
	s.Router.Use(middleware.Logging)
	s.Router.Use(chimw.Recoverer)
	s.Router.Use(middleware.CORS(frontendURL))
	s.Router.Use(chimw.Timeout(30 * time.Second))

	s.RegisterRoutes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
