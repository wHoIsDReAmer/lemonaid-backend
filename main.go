package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"lemonaid-backend/db"
	"lemonaid-backend/dotenv"
	v1 "lemonaid-backend/routes/v1"
	"os"
)

func main() {
	dotenv.Load(".env")

	var port = os.Getenv("PORT")

	db.Init() // init the db

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		BodyLimit:             30 * 1024 * 1024,
	})

	if os.Getenv("DEV") == "true" {
		app.Use(cors.New(cors.Config{
			AllowCredentials: true,
			ExposeHeaders:    "Content-Disposition",
		}))
	}

	v1.Controller(app)

	app.Static("/", "./public")

	app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	fmt.Println("Server listening at " + port)
	if err := app.Listen(":" + port); err != nil {
		fmt.Println("Server got error: ", err)
	}
}
