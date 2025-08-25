package deviceservice

func (d *DeviceHTTP) loadRoutes() {
	d.router.Post("/login", d.Login)
	d.router.With(GetAuthMiddlewareFunc()).Get("/home", d.Home)
}
