package server

import "github.com/go-chi/chi/v5"

func (s *Server) RegisterRoutes() {
	s.Router.Get("/healthz", s.HealthHandler.Healthz)
	s.Router.Get("/readyz", s.HealthHandler.Readyz)

	s.Router.Route("/api/v1", func(r chi.Router) {
		// Phase 4: analysis endpoints
		// Phase 5: history and refinement endpoints
	})
}
