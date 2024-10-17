package reports

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, reportsHttpApi ReportsHTTPTransport, authMiddleware func(c *fiber.Ctx) error) {
	reportsRoutes := router.Group("/reports")
	reportsRoutes.Post("", authMiddleware, reportsHttpApi.Create)
	reportsRoutes.Get("", authMiddleware, reportsHttpApi.GetReports)
	reportsRoutes.Get("/:id", authMiddleware, reportsHttpApi.GetReportByID)
	reportsRoutes.Put("/:id", authMiddleware, reportsHttpApi.UpdateReport)
}
