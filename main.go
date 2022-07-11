package main

import (
	"system/database"
	"system/routes"
	"system/memotest"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	allPieces, err := database.GetFullSymbolSet()
	if(err != nil) {
		panic("Getting symbols: "+err.Error())
	}
	if(len(allPieces)==0) {
		panic("No pieces loaded.")
	}
	memotest.CtrlStart(allPieces)
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)
	app.Listen(":8000")
}
