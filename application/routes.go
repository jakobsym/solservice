package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jakobsym/solservice/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/watchlist", loadWatchlistRoutes)
	//router.Roue("/token", loadTokenRoutes)

	return router
}

func loadWatchlistRoutes(router chi.Router) {
	watchlistHandler := &handler.Watchlist{}

	router.Post("/", watchlistHandler.Create)
	router.Get("/", watchlistHandler.List)
	router.Get("/{id}", watchlistHandler.GetByID)
	router.Delete("/{id}", watchlistHandler.DeleteByID)
	router.Put("/{id}", watchlistHandler.UpdateByID)
}

// loadTokenRoutes(router chi.Router) {}