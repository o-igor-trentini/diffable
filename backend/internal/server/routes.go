package server

import "github.com/go-chi/chi/v5"

func (s *Server) RegisterRoutes() {
	s.Router.Get("/healthz", s.HealthHandler.Healthz)
	s.Router.Get("/readyz", s.HealthHandler.Readyz)

	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Post("/analyses/commit", s.AnalysisHandler.AnalyzeCommit)
		r.Post("/analyses/range", s.AnalysisHandler.AnalyzeRange)
		r.Post("/analyses/pr", s.AnalysisHandler.AnalyzePR)
		r.Get("/analyses/{id}", s.AnalysisHandler.GetAnalysis)

		// Phase 5: history and refinement endpoints
	})
}
