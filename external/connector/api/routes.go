package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (h *Handler) SetupRouter() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	//router.Use(middleware.AllowContentType("application/json")) TODO
	router.Use(middleware.CleanPath)
	//todo
	//router.Use(csrf)
	router.Use(httprate.Limit(
		10,             // requests
		10*time.Second, // per duration
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
		httprate.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, fmt.Sprintf(`{"error": %q}`, err), http.StatusPreconditionRequired)
		}),
	))
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"}, // change it to only http!
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	h.router = router
	h.loadRoutes()
}
func (h *Handler) loadRoutes() {
	h.router.Post("/login", h.Login)

	//have to be logged in
	h.router.Group(func(r chi.Router) {
		r.Use(GetAuthMiddlewareFunc())
		r.Get("/home", h.Home)
		r.Post("/location", h.updateLocation)
	})
}

func (h *Handler) RunHTTP() {
	// load routes
	log.Printf("HTTP server is runnig on: %v \n", h.listenAddr)
	log.Fatalln(http.ListenAndServe(h.listenAddr, h.router))
}

//TODO maybe add to user table kind: website/mobile or etc

//and block e.g. lvl3->mobile user
//and window app only for 4,5 lvl
