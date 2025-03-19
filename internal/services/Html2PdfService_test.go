package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/dtos"
	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/kolzxx/html2pdf/internal/services"
	"github.com/stretchr/testify/assert"
)

const jsonHeader = `"<div style="color: black; font-size: 18pt; padding-top: 5px; text-align: center; width: 100%; height: 100pt;"> <p> <span style="font-family: 'Quattrocento Sans'">N&ordm; do contrato: </span> <span style="font-weight: bold;font-family: 'Quattrocento Sans'">@NumeroContrato</span> </p> </div>`
const jsonFooter = `<div style="text-align: right;width: 297mm;font-size: 8px;opacity: 0.7;"><span style="margin-right: 1cm"><span class="pageNumber"></span> / <span class="totalPages"></span></span></div>`
const jsonContent = `<html><body><h1>My First Heading</h1><p>My first paragraph.</p>/body></html>`
const jsonContentCss = `* {            font-family: system-ui, system-ui, sans-serif;            font-size: 18pt;        }        body {            padding: 10pt;            margin: 15pt 20pt;        }        header, footer, h2, h5 {            text-align: center;            color: #979797;        }        h2 {            font-size: 19pt;        }        h5 {            font-size: 15pt;        }        h2, h5 {            margin: 0;        }        .tab {            tab-size: 4;        }        .assinatura {            margin: 20pt 0pt;            text-align: center;        }        table,        h2,        h5,        header,        footer,        #data,        #emitente {            border-spacing: 0pt;            width: 100%;        }        table {            margin-top: 30pt;            table-layout: fixed;            border: 1pt solid black;        }        th {            border-top: 1pt solid black;            border-left: 1pt dotted black;            border-right: 1pt dotted black;            border-bottom: 1pt dotted black;        }        td {            border: 1pt dotted black;        }        td {            vertical-align: bottom;        }        th, td {            padding: 5pt;            text-align: left;        }            td.small {                width: 5%;            }        .assin {            text-align: center;            margin-top: 20pt;            font-size: 18pt;        }        .nc {            text-align: center;            font-weight: bold;        }       ol {            padding: 0;            margin-left: 15pt;            margin-right: 15pt;            text-align: justify;        }        .center {            text-align: center;        }        .espaco {            margin-left: 15pt;        }        footer {            margin: 30pt 0;        }            footer div {                position: relative;            }            footer #page, footer #info {                position: absolute;            }            footer #page {                right: 8%;                border: 1pt solid #979797;                padding: 5pt 10pt;                top: -10pt;            }            footer #info {                left: 0;                bottom: 0;                font-size: 11pt;            }        p {            text-align: justify;        }`

func TestHtml2PdfService(t *testing.T) {

	t.Run("NewHtml2PdfService", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()

		hs := services.NewHtml2PdfService(logger, cdp)

		assert.NotNil(t, hs)
		assert.NotNil(t, cdp)
	})

	t.Run("TestHtmlToPdf", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()
		hs := services.NewHtml2PdfService(logger, cdp)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss

		pdfResponse, err := hs.HtmlToPdf(obj)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)
	})

	t.Run("TestHtmlToPdfError", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()
		hs := services.NewHtml2PdfService(logger, cdp)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = "&¨&$%¨&&*(())"
		obj.ContentCss = jsonContentCss

		pdfResponse, err := hs.HtmlToPdf(obj)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)
	})

	t.Run("TestHtmlToPdfContextNil", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()
		hs := services.NewHtml2PdfService(logger, cdp)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = "&¨&$%¨&&*(())"
		obj.ContentCss = jsonContentCss

		cdp.Context = nil

		pdfResponse, err := hs.HtmlToPdf(obj)

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)
	})

	t.Run("TestWriteHTML", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()
		hs := services.NewHtml2PdfService(logger, cdp)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss

		pdfResponse, err := hs.HtmlToPdf(obj)
		defer cdp.Cancelf()

		assert.NotNil(t, hs)
		assert.NotNil(t, pdfResponse)
		assert.Nil(t, err)

		handler := cdp.WriteHTML()

		assert.NotNil(t, hs)
		assert.NotNil(t, handler)
	})

	t.Run("TestDoPdfActions", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()
		hs := services.NewHtml2PdfService(logger, cdp)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss

		cxtt, cancelt := context.WithTimeout(cdp.Context, time.Second*20)
		defer cancelt()
		taskCtx, cancel := chromedp.NewContext(cxtt)
		defer cancel()

		var pdfBuffer []byte

		ts := httptest.NewServer(hs.WriteHTML())
		defer ts.Close()

		chromedp.Run(taskCtx, hs.PdfGrabber(ts.URL, &pdfBuffer, obj))

		err := hs.DoPdfActions(&pdfBuffer, obj, taskCtx)
		defer cdp.Cancelf()

		assert.EqualError(t, err, "invalid context")
		assert.NotNil(t, &pdfBuffer)
	})

	t.Run("TestDoHandler", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()
		hs := services.NewHtml2PdfService(logger, cdp)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss
		pdfResponse, err := hs.HtmlToPdf(obj)

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

	t.Run("TestDoPrint", func(t *testing.T) {
		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent

		buf, err := services.DoPrint(context.Background(), obj)

		assert.Nil(t, buf)
		assert.NotNil(t, err)

	})

	t.Run("TestDoHandlerChrome", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		w := httptest.NewRecorder()
		jsonValue, _ := json.Marshal(obj)
		req, _ := http.NewRequest(http.MethodPost, "/v1/html2pdf", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		//Assertion
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "404 page not found", w.Body.String())

		cdp.DoHandler(w, req)
		assert.NotNil(t, w)

	})

	// t.Run("DoEventLoad", func(t *testing.T) {
	// 	logger := logger.NewFakeLogger()

	// 	cdp := services.NewChromedpService(context.Background(), logger)
	// 	cdp.RunChromeDp()
	// 	hs := services.NewHtml2PdfService(logger, cdp)

	// 	obj := dtos.HtmlRequest{}
	// 	obj.HeaderTemplate = jsonHeader
	// 	obj.FooterTemplate = jsonFooter
	// 	obj.Content = jsonContent
	// 	obj.ContentCss = jsonContentCss

	// 	gin.SetMode(gin.TestMode)
	// 	router := gin.Default()

	// 	w := httptest.NewRecorder()
	// 	jsonValue, _ := json.Marshal(obj)
	// 	req, _ := http.NewRequest(http.MethodPost, "/v1/html2pdf", bytes.NewBuffer(jsonValue))
	// 	router.ServeHTTP(w, req)
	// 	//Assertion
	// 	assert.Equal(t, http.StatusNotFound, w.Code)
	// 	assert.Equal(t, "404 page not found", w.Body.String())

	// 	cxtt, cancelt := context.WithTimeout(cdp.Context, time.Second*20)
	// 	defer cancelt()
	// 	// taskCtx, cancel := chromedp.NewContext(cxtt)
	// 	// defer cancel()

	// 	var wg sync.WaitGroup
	// 	ch := make(chan int, 5)
	// 	wg.Add(1)
	// 	ch <- 1

	// 	err := hs.DoEventLoad(cxtt, &wg, &ch)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, w)
	// })

	t.Run("DoGetFrameTree", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		cdp := services.NewChromedpService(context.Background(), logger)
		cdp.RunChromeDp()
		hs := services.NewHtml2PdfService(logger, cdp)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss

		if cdp.Context == nil {
			cdp.Context = context.Background()
		}
		hs.HtmlToPdf(obj)

		var wg sync.WaitGroup
		ch := make(chan int, 5)
		wg.Add(1)
		ch <- 1

		wait := hs.Wait(&wg)
		hs.Wait(&wg)
		assert.NotNil(t, wait)
	})

}
