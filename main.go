package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	analysesvc "vezhguesi/app/analyses"
	"vezhguesi/app/articles"
	articlesvc "vezhguesi/app/articles"
	entitysvc "vezhguesi/app/entities"
	entity_reportsvc "vezhguesi/app/entity_reports"
	orgsvc "vezhguesi/app/orgs"
	reportsvc "vezhguesi/app/reports"
	subscriptionsvc "vezhguesi/app/subscriptions"
	session "vezhguesi/core/authentication"
	authsvc "vezhguesi/core/authentication/auth"
	rolesvc "vezhguesi/core/authorization/role"
	db "vezhguesi/core/db"
	dbseeds "vezhguesi/core/db/seeds"
	"vezhguesi/core/middleware"
	usersvc "vezhguesi/core/users"
	_ "vezhguesi/docs" // Import the generated docs package
	server "vezhguesi/sentiment-communication"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"gopkg.in/gomail.v2" // Import gomail
	"gorm.io/gorm"
	// http-swagger middleware
)

type StringArray []string

func (a StringArray) Value() (interface{}, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(b, a)
}

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
			"influxo": "123123123",
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
		reportsvc.NewReportsAPI(db, dialer, os.Getenv("UI_APP_URL"), defaultLogger, entitysvc.NewEntitiesAPI(db, defaultLogger), server.NewServerAPI(db, defaultLogger)),
	)
	orgApiSvc := orgsvc.NewOrgHTTPTransport(
		orgsvc.NewOrgAPI(db, defaultLogger),
		defaultLogger,
	)

	
	// Register Routes
	usersvc.RegisterRoutes(apisRouter, userAPISvc, authMiddleware)
	authsvc.RegisterRoutes(apisRouter, authApiSvc, authMiddleware, middleware.SessionMiddleware(db))
	reportsvc.RegisterRoutes(apisRouter, reportApiSvc, authMiddleware)
	entitysvc.RegisterRoutes(apisRouter, entityApiSvc)
	orgsvc.RegisterRoutes(apisRouter, orgApiSvc, authMiddleware)
	// Auto Migrate Core
	db.AutoMigrate(
		&usersvc.User{},
		&rolesvc.Role{},
		&rolesvc.Permission{},
		&session.Session{}, // Add this line
	)
	
	// Auto Migrate App
	db.AutoMigrate(
		&reportsvc.Report{},
		&entitysvc.Entity{},
		&orgsvc.Org{},
		&orgsvc.UserOrgRole{},
		&subscriptionsvc.Subscription{},
		&subscriptionsvc.Feature{},
		&articles.Article{},
		&articles.ArticleEntity{},
		&entity_reportsvc.EntityReport{},
		&entity_reportsvc.EntityReportArticle{},
		&entity_reportsvc.UserEntityReport{},
		&analysesvc.Analysis{},
	)

	dbseeds.SeedDefaultRolesAndPermissions(db)

	// Start article fetching in a separate goroutine
	go scheduledArticleFetch(server.NewServerAPI(db, defaultLogger))

	// go scheduledEntityCheck(db, defaultLogger)

	// Start the server
	log.Fatal(app.Listen(fmt.Sprintf(`:%d`, 3001)))
}

func scheduledArticleFetch(api server.ServerAPI) {
	ticker := time.NewTicker(1 * time.Hour) // Adjust interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := api.FetchAndStoreArticles(); err != nil {
				fmt.Printf("Error fetching articles: %v", err)
			}
		}
	}
}

func scheduledEntityCheck(db *gorm.DB, logger log.AllLogger) {
	ticker := time.NewTicker(10 * time.Second) // Adjust interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := checkArticlesForEntities(db, logger); err != nil {
				logger.Errorf("Error checking articles for entities: %v", err)
			}
		}
	}
}

func checkArticlesForEntities(db *gorm.DB, logger log.AllLogger) error {
	// Begin transaction
	tx := db.Begin()

	// Get all entities
	var entities []entitysvc.Entity
	if err := tx.Find(&entities).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch entities: %v", err)
	}

	// Get all articles
	var articles []articlesvc.Article
	if err := tx.Find(&articles).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch articles: %v", err)
	}

	// For each article, check for entity mentions
	for _, article := range articles {
		content := strings.ToLower(article.Content)
		title := strings.ToLower(article.Title)

		for _, entity := range entities {
			entityName := strings.ToLower(entity.Name)
			
			// Check if entity is mentioned in title or content
			if strings.Contains(content, entityName) || strings.Contains(title, entityName) {
				// Check if relation already exists
				var existingRelation articlesvc.ArticleEntity
				err := tx.Where("article_id = ? AND entity_name = ?", article.ID, entity.Name).
					First(&existingRelation).Error

				if err == gorm.ErrRecordNotFound {
					// Create new relation if it doesn't exist
					relation := articlesvc.ArticleEntity{
						ArticleID:      article.ID,
						EntityName:     entity.Name,
						SentimentScore: 0,
						SentimentLabel: "neutral",
					}
					
					if err := tx.Create(&relation).Error; err != nil {
						tx.Rollback()
						logger.Errorf("Failed to create article-entity relation: %v", err)
						return fmt.Errorf("failed to create article-entity relation: %v", err)
					}
					
					logger.Infof("Created new relation: Article %d - Entity %s", article.ID, entity.Name)
				}
			}
		}
	}

	return tx.Commit().Error
}