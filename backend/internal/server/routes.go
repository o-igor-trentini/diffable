package server

import "github.com/go-chi/chi/v5"

func (s *Server) RegisterRoutes() {
	s.Router.Get("/healthz", s.HealthHandler.Healthz)
	s.Router.Get("/readyz", s.HealthHandler.Readyz)

	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Post("/analyses/commit", s.AnalysisHandler.AnalyzeCommit)
		r.Post("/analyses/range", s.AnalysisHandler.AnalyzeRange)
		r.Post("/analyses/pr", s.AnalysisHandler.AnalyzePR)
		r.Get("/analyses", s.HistoryHandler.ListAnalyses)
		r.Get("/analyses/{id}", s.AnalysisHandler.GetAnalysis)
		r.Post("/analyses/{id}/refine", s.AnalysisHandler.RefineDescription)
		r.Get("/analyses/{id}/refinements", s.HistoryHandler.GetRefinements)
	})
}
