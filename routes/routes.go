package routes

import (
	"system/controllers"
	memotest "system/memotest"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	app.Post("/api/memo", memotest.CreateGame)
	app.Get("/api/memo/:gameId", memotest.ShowGame)
	app.Post("/api/memo/:gameId", memotest.JoinGame)
	app.Put("/api/memo/:gameId/:pieceId", memotest.SelectPiece)
	app.Delete("/api/memo/:gameId/:pieceId", memotest.DeselectPiece)
}
