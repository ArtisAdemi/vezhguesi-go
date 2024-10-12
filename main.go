package main

import (
	"fmt"
	"log"

	db "vezhguesi/core/db"
	usersvc "vezhguesi/core/users"
	_ "vezhguesi/docs" // Import the generated docs package

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	// http-swagger middleware
)

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB in bytes
	})
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	apisRouter := app.Group("/api")

	apisRouter.Get("/swagger/*", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"babuqi": "dedidedi123",
		},
	}), swagger.HandlerDefault)


	userAPISvc := usersvc.NewUserHTTPTransport(
		usersvc.NewUserAPI(db),
	)

	usersvc.RegisterRoutes(apisRouter, userAPISvc)

	router := fiber.New()
	usersvc.RegisterRoutes(router, userAPISvc)

	db.AutoMigrate(&usersvc.User{})


	// Start the server
	log.Fatal(app.Listen(fmt.Sprintf(`:%d`, 8080)))
}
