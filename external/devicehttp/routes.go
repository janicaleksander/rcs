package devicehttp

import (
	"github.com/go-chi/chi/v5"
)

func (d *DeviceHTTP) loadRoutes() {
	d.router.Route("/", d.loadHomeRoutes)
	d.router.Route("/login", d.loadLoginRoutes)
}

func (d *DeviceHTTP) loadHomeRoutes(r chi.Router) {
	r.Get("/", d.home)
}

func (d *DeviceHTTP) loadLoginRoutes(r chi.Router) {
	r.Post("/login", d.login)
}
