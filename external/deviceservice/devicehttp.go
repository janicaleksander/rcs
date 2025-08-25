package deviceservice

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

type DeviceHTTP struct {
	ctx        *actor.Context
	serverPID  *actor.PID
	listenAddr string
	router     *chi.Mux
	//devicePID  *actor.PID
	//parent PID in the future
}

func NewHTTPDevice(addr string, ctx *actor.Context, pid *actor.PID) *DeviceHTTP {
	return &DeviceHTTP{
		serverPID:  pid,
		ctx:        ctx,
		listenAddr: addr,
		router:     nil,
		//devicePID:  pid,
	}
}

// TODO idk if panic or return err
func (d *DeviceHTTP) RunHTTPServer() {
	//setup chi router
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
	//assign to deviceHTTP
	d.router = router

	// load routes
	d.loadRoutes()
	log.Printf("Server is runnig on: %v \n", d.listenAddr)
	log.Fatalln(http.ListenAndServe(d.listenAddr, router))
}
