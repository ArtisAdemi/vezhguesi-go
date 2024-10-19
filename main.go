package main

import (
	"fmt"
	"os"

	entitysvc "vezhguesi/app/entities"
	reportsvc "vezhguesi/app/reports"
	authsvc "vezhguesi/core/authentication/auth"
	db "vezhguesi/core/db"
	"vezhguesi/core/middleware"
	usersvc "vezhguesi/core/users"
	_ "vezhguesi/docs" // Import the generated docs package

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"gopkg.in/gomail.v2" // Import gomail
	// http-swagger middleware
)

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB in bytes
	})

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Change this to specific domains in production
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	defaultLogger := log.DefaultLogger()

	apisRouter := app.Group("/api")

	apisRouter.Get("/swagger/*", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"babuqi": "dedidedi123",
		},
	}), swagger.HandlerDefault)

	// Pass gomail dialer to user service
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_FROM"), os.Getenv("MAIL_PASSWORD"))
	

	authMiddleware := middleware.Authentication(os.Getenv("JWT_SECRET_KEY"))
	// API Services
	userAPISvc := usersvc.NewUserHTTPTransport(
		usersvc.NewUserAPI(db, os.Getenv("JWT_SECRET_KEY"), dialer, os.Getenv("UI_APP_URL"), defaultLogger),
	)
	authApiSvc := authsvc.NewAuthHTTPTransport(
		authsvc.NewAuthApi(db, os.Getenv("JWT_SECRET_KEY"), dialer, os.Getenv("UI_APP_URL"), defaultLogger),
	)
	entityApiSvc := entitysvc.NewEntitiesHTTPTransport(
		entitysvc.NewEntitiesAPI(db, defaultLogger),
	)
	reportApiSvc := reportsvc.NewReportsHTTPTransport(
		reportsvc.NewReportsAPI(db, dialer, os.Getenv("UI_APP_URL"), defaultLogger, entitysvc.NewEntitiesAPI(db, defaultLogger)),
	)

	// Register Routes
	usersvc.RegisterRoutes(apisRouter, userAPISvc, authMiddleware)
	authsvc.RegisterRoutes(apisRouter, authApiSvc, authMiddleware)
	reportsvc.RegisterRoutes(apisRouter, reportApiSvc, authMiddleware)
	entitysvc.RegisterRoutes(apisRouter, entityApiSvc)
	// Auto Migrate Core
	db.AutoMigrate(
		&usersvc.User{},
	)

	// Auto Migrate App
	db.AutoMigrate(
		&reportsvc.Report{},
		&entitysvc.Entity{},
	)

	// Start the server
	log.Fatal(app.Listen(fmt.Sprintf(`:%d`, 3001)))
}
