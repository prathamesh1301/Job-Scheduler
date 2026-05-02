package main

import (
	authjwt "auth/internals/jwt"
	"auth/internals/redis"
	"auth/internals/store"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type application struct {
	config   config
	store    *store.Store
	jwt      authjwt.TokenService
	dbConfig dbConfig
	redis    *redis.Redis
}

type config struct {
	address string
}

type dbConfig struct {
	address      string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) getMuxHandler() http.Handler {
	mux := chi.NewRouter()

	// -------- Public routes --------
	mux.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		})
		r.Post("/login", app.login)
		r.Post("/signup", app.signup)
		r.Post("/refreshToken", app.refreshTokenHandler)
	})

	// -------- Protected routes --------
	mux.Group(func(r chi.Router) {
		r.Use(app.ValidateTokenMiddleware)
		r.Get("/health", app.checkhealth)
		r.Post("/addJob",app.addJobHandler)
		r.Get("/failedJobs",app.getFailedJobs)
	})
	return mux
}

func (app *application) startServer() {
	server := &http.Server{
		Addr:    app.config.address,
		Handler: app.getMuxHandler(),
	}

	fmt.Println("Starting server on ", app.config.address)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
