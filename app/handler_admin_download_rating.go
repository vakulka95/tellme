package app

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/app/representation"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func (s *apiserver) webAdminExpertRatingExcel(c *gin.Context) {
	qlp := &representation.QueryListParams{}

	if err := c.BindQuery(qlp); err != nil {
		log.Printf("(ERR) Bind query error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	list, err := s.repository.GetExpertRatingTable(representation.QueryExpertRatingAPItoPersistence(qlp))
	if err != nil {
		log.Printf("(ERR) Failed to fetch expert rating list: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	file := s.generateExpertRatingExcel(list)
	b, err := file.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	downloadName := time.Now().UTC().Format("expert_rating_20060102150405.xlsx")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Header("File-Name", "userInputData.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")

	c.Data(http.StatusOK, "application/octet-stream", b.Bytes())
}

func (s *apiserver) generateExpertRatingExcel(list []*model.Expert) *excelize.File {
	file := excelize.NewFile()
	sheetName := "Рейтинг психологів"
	index := file.NewSheet(sheetName)
	file.SetActiveSheet(index)

	file.SetCellValue(sheetName, "A1", "Ім'я")
	file.SetCellValue(sheetName, "B1", "Телефон")
	file.SetCellValue(sheetName, "C1", "Email")
	file.SetCellValue(sheetName, "D1", "Кількість відгуків")
	file.SetCellValue(sheetName, "E1", "Середня оцінка")
	file.SetCellValue(sheetName, "F1", "Заявок завершено")
	file.SetCellValue(sheetName, "G1", "Заявок в роботі")
	file.SetCellValue(sheetName, "H1", "Сесій проведено")
	file.SetCellValue(sheetName, "I1", "Дата реєстрації")

	for i, expert := range list {
		index := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", index), expert.Username)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", index), expert.Phone)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", index), expert.Email)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", index), expert.ReviewCount)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", index), fmt.Sprintf("%.2f", expert.AverageRating))
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", index), expert.CompletedCount)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", index), expert.ProcessingCount)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", index), expert.SessionCount)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", index), expert.CreatedAt.Format("2006-01-02 15:04"))
	}

	file.SetColWidth(sheetName, "A", "A", 30)
	file.SetColWidth(sheetName, "B", "B", 13)
	file.SetColWidth(sheetName, "C", "C", 30)
	file.SetColWidth(sheetName, "D", "H", 5)
	file.SetColWidth(sheetName, "I", "I", 13)
	file.DeleteSheet("Sheet1")
	return file
}

func (s *apiserver) webAdminExpertRatingPDF(c *gin.Context) {
	qlp := &representation.QueryListParams{}

	if err := c.BindQuery(qlp); err != nil {
		log.Printf("(ERR) Bind query error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	list, err := s.repository.GetExpertRatingTable(representation.QueryExpertRatingAPItoPersistence(qlp))
	if err != nil {
		log.Printf("(ERR) Failed to fetch expert rating list: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	file := s.generateExpertRatingPDF(list)
	if err := file.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to generate file: %v", err.Error()))
		return
	}

	b := &bytes.Buffer{}
	err = file.Output(b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to write file: %v", err.Error()))
		return
	}

	downloadName := time.Now().UTC().Format("expert_rating_20060102150405.pdf")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Header("File-Name", "userInputData.pdf")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")

	c.Data(http.StatusOK, "application/octet-stream", b.Bytes())
}

func (s *apiserver) generateExpertRatingPDF(list []*model.Expert) *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "A4", "/usr/share/tellme/static/pdf_fonts")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 9)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	pdf.SetFillColor(240, 240, 240)

	pdf.CellFormat(40, 8, tr("Ім'я"), "1", 0, "", true, 0, "")
	pdf.CellFormat(18, 8, tr("Телефон"), "1", 0, "", true, 0, "")
	pdf.CellFormat(54, 8, tr("Email"), "1", 0, "", true, 0, "")
	pdf.CellFormat(30, 8, tr("Кількість відгуків"), "1", 0, "", true, 0, "")
	pdf.CellFormat(26, 8, tr("Середня оцінка"), "1", 0, "", true, 0, "")
	pdf.CellFormat(30, 8, tr("Заявок завершено"), "1", 0, "", true, 0, "")
	pdf.CellFormat(27, 8, tr("Заявок в роботі"), "1", 0, "", true, 0, "")
	pdf.CellFormat(27, 8, tr("Сесій проведено"), "1", 0, "", true, 0, "")
	pdf.CellFormat(25, 8, tr("Дата реєстрації"), "1", 0, "", true, 0, "")

	pdf.Ln(-1)                      //newline height equals the last line height
	pdf.SetFillColor(255, 255, 255) //white color

	pdf.SetFont("Helvetica", "", 8)
	for _, expert := range list {
		pdf.CellFormat(40, 8, tr(expert.Username), "1", 0, "L", false, 0, "")
		pdf.CellFormat(18, 8, tr(expert.Phone), "1", 0, "L", false, 0, "")
		pdf.CellFormat(54, 8, expert.Email, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", expert.ReviewCount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(26, 8, fmt.Sprintf("%.2f", expert.AverageRating), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", expert.CompletedCount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, 8, fmt.Sprintf("%d", expert.ProcessingCount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, 8, fmt.Sprintf("%d", expert.SessionCount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, expert.CreatedAt.Format("2006-01-02 15:04"), "1", 0, "L", false, 0, "")

		pdf.Ln(-1)
	}
	pdf.Ln(-1)

	return pdf
}
