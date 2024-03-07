package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jakobsym/solservice/handler"
	"github.com/jakobsym/solservice/repository/token"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/watchlist", a.loadWatchlistRoutes)
	router.Route("/token", a.loadTokenRoutes)
	a.router = router
}

func (a *App) loadWatchlistRoutes(router chi.Router) {
	watchlistHandler := &handler.Watchlist{}

	router.Post("/", watchlistHandler.Create)
	router.Get("/", watchlistHandler.List)
	router.Get("/{id}", watchlistHandler.GetByID)
	router.Delete("/{id}", watchlistHandler.DeleteByID)
	router.Put("/{id}", watchlistHandler.UpdateByID)
}

func (a *App) loadTokenRoutes(router chi.Router) {
	tokenHandler := &handler.TokenHandler{
		Repo: &token.MySqlRepo{
			DB: a.db,
		},
	}

	router.Post("/", tokenHandler.Create)
	router.Get("/", tokenHandler.List)
	router.Get("/{ca}", tokenHandler.GetByCA)
	router.Delete("/{ca}", tokenHandler.DeleteByCA)
}
