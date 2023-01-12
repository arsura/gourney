package api

import "github.com/gofiber/fiber/v2"

func (app *Application) routes(server *fiber.App) {
	server.Get("/posts/count-by-social-network", app.Handlers.Post.CountPostBySocialNetworkTypeHandler)
	server.Get("/posts/:id", app.Handlers.Post.FindPostByIdHandler)
	server.Post("/posts", app.Handlers.Post.CreatePostHandler)
	server.Patch("/posts", app.Handlers.Post.UpdatePostByIdHandler)
	server.Delete("/posts/:id", app.Handlers.Post.DeletePostByIdHandler)
}
