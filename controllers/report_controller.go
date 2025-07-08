package controllers

import (
	"case_study_api/services"
	"case_study_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService services.ReportService
}

func NewReportController(reportService services.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

func (rc *ReportController) SummaryReport(c *gin.Context) {
	summary, err := rc.reportService.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("failed to generate summary report"))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("success", summary))
}

func (rc *ReportController) EventReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid event id"))
		return
	}

	report, err := rc.reportService.GetEventReport(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("failed to generate event report"))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("success", report))
}

func (rc *ReportController) SystemReport(c *gin.Context) {
	report, err := rc.reportService.GetSystemReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("failed to generate system report"))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("success", report))
}

func (rc *ReportController) SystemReportPDF(c *gin.Context) {
	pdfData, err := rc.reportService.GenerateSystemReportPDF()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("failed to generate PDF report"))
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=malaka_ticket_system_report.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfData)
}
