package deviceservice

import "github.com/go-chi/chi/v5"

func (d *DeviceHTTP) loadRoutes() {
	d.router.Post("/login", d.Login)

	//have to be logged in
	d.router.Group(func(r chi.Router) {
		r.Use(GetAuthMiddlewareFunc())
		r.Get("/home", d.Home)
	})
}
