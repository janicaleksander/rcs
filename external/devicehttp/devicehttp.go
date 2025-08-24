package devicehttp

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type DeviceHTTP struct {
	listenAddr string
	router     *chi.Mux
	//parent PID in the future
}

func New(addr string) *DeviceHTTP {
	return &DeviceHTTP{
		listenAddr: addr,
		router:     nil,
	}
}

// TODO idk if panic or return err
func (d *DeviceHTTP) RunHTTPServer() {
	//setup chi router
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	//assign to deviceHTTP
	d.router = router

	// load routes
	d.loadRoutes()

	//run
	log.Printf("Server is runnig on: %v \n", d.listenAddr)
	log.Fatalln(http.ListenAndServe(d.listenAddr, router))
}
