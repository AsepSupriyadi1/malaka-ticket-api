package routes

import (
	"case_study_api/container"
	"case_study_api/controllers"
	"case_study_api/middleware"

	"github.com/gin-gonic/gin"
)

func ReportRoutes(rg *gin.RouterGroup, container *container.Container) {
	reportController := controllers.NewReportController(container.ReportService)

	report := rg.Group("/reports")
	report.Use(middleware.RoleAuth("admin"))
	report.GET("/summary", reportController.SummaryReport)
	report.GET("/event/:id", reportController.EventReport)
	report.GET("/system", reportController.SystemReport)
	report.GET("/system/pdf", reportController.SystemReportPDF)
}
