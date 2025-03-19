package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/dtos"
	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/kolzxx/html2pdf/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestHtml2PdfServiceMock(t *testing.T) {
	t.Parallel()

	t.Run("NewHtml2PdfServiceMock", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		hs := services.NewHtml2PdfServiceMock(logger)

		assert.NotNil(t, hs)
	})

	t.Run("TestHtmlToPdfMock", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		hs := services.NewHtml2PdfServiceMock(logger)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent

		pdfResponse, err := hs.HtmlToPdf(obj)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)
	})

	t.Run("TestHtmlToPdfErrorMock", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		hs := services.NewHtml2PdfServiceMock(logger)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = "&¨&$%¨&&*(())"

		pdfResponse, err := hs.HtmlToPdf(obj)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)
	})

	t.Run("TestWriteHTMLMock", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		hs := services.NewHtml2PdfServiceMock(logger)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent

		pdfResponse, err := hs.HtmlToPdf(obj)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)

		handler := hs.WriteHTML()

		assert.NotNil(t, hs)
		assert.NotNil(t, handler)
	})

	t.Run("TestPdfActionsMock", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		hs := services.NewHtml2PdfServiceMock(logger)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		var pdfBuffer []byte

		action := hs.PdfActions(&pdfBuffer, obj)

		assert.NotNil(t, action)
	})

	t.Run("TestPdfGrabberMock", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		hs := services.NewHtml2PdfServiceMock(logger)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		var pdfBuffer []byte
		pdfResponse, err := hs.HtmlToPdf(obj)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)

		ts := httptest.NewServer(hs.WriteHTML())

		task := hs.PdfGrabber(ts.URL, &pdfBuffer, obj)

		assert.NotNil(t, task)
	})

	t.Run("TestDoPrintMock", func(t *testing.T) {
		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent

		buf, err := services.DoPrintMock(context.Background(), obj)

		assert.Nil(t, buf)
		assert.NotNil(t, err)

	})

	t.Run("TestDoHandlerMock", func(t *testing.T) {
		logger := logger.NewFakeLogger()
		// w := httptest.NewRecorder()
		//sc, _ := gin.CreateTestContext(w)

		hs := services.NewHtml2PdfServiceMock(logger)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		pdfResponse, err := hs.HtmlToPdf(obj)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)

		buf := hs.DoHandler

		assert.NotNil(t, buf)
		assert.Nil(t, err)

	})

	t.Run("TestDoPdfActions", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		hs := services.NewHtml2PdfServiceMock(logger)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss
		_, err := hs.HtmlToPdf(obj)
		assert.Nil(t, err)

		var pdfBuffer []byte

		err = hs.DoPdfActions(&pdfBuffer, obj, context.Background())

		assert.EqualError(t, err, "invalid context")
		assert.NotNil(t, &pdfBuffer)
	})

	t.Run("TestDoHandler", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		hs := services.NewHtml2PdfServiceMock(logger)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss
		pdfResponse, err := hs.HtmlToPdf(obj)
		assert.Nil(t, err)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)

		buf := hs.DoHandler

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		w := httptest.NewRecorder()
		jsonValue, _ := json.Marshal(obj)
		req, _ := http.NewRequest(http.MethodPost, "/v1/html2pdf", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		//Assertion
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "404 page not found", w.Body.String())

		hs.DoHandler(w, req)
		assert.NotNil(t, buf)
		assert.Nil(t, err)

	})

}
