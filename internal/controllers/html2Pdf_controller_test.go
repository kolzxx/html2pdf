package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/controllers"
	"github.com/kolzxx/html2pdf/internal/dtos"
	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/stretchr/testify/assert"
)

const jsonHeader = `"<div style="color: black; font-size: 18pt; padding-top: 5px; text-align: center; width: 100%; height: 100pt;"> <p> <span style="font-family: 'Quattrocento Sans'">N&ordm; do contrato: </span> <span style="font-weight: bold;font-family: 'Quattrocento Sans'">@NumeroContrato</span> </p> </div>`
const jsonFooter = `<div style="text-align: right;width: 297mm;font-size: 8px;opacity: 0.7;"><span style="margin-right: 1cm"><span class="pageNumber"></span> / <span class="totalPages"></span></span></div>`
const jsonContent = `<html><body><h1>My First Heading</h1><p>My first paragraph.</p>/body></html>`
const jsonContentCss = `* {            font-family: system-ui, system-ui, sans-serif;            font-size: 18pt;        }        body {            padding: 10pt;            margin: 15pt 20pt;        }        header, footer, h2, h5 {            text-align: center;            color: #979797;        }        h2 {            font-size: 19pt;        }        h5 {            font-size: 15pt;        }        h2, h5 {            margin: 0;        }        .tab {            tab-size: 4;        }        .assinatura {            margin: 20pt 0pt;            text-align: center;        }        table,        h2,        h5,        header,        footer,        #data,        #emitente {            border-spacing: 0pt;            width: 100%;        }        table {            margin-top: 30pt;            table-layout: fixed;            border: 1pt solid black;        }        th {            border-top: 1pt solid black;            border-left: 1pt dotted black;            border-right: 1pt dotted black;            border-bottom: 1pt dotted black;        }        td {            border: 1pt dotted black;        }        td {            vertical-align: bottom;        }        th, td {            padding: 5pt;            text-align: left;        }            td.small {                width: 5%;            }        .assin {            text-align: center;            margin-top: 20pt;            font-size: 18pt;        }        .nc {            text-align: center;            font-weight: bold;        }       ol {            padding: 0;            margin-left: 15pt;            margin-right: 15pt;            text-align: justify;        }        .center {            text-align: center;        }        .espaco {            margin-left: 15pt;        }        footer {            margin: 30pt 0;        }            footer div {                position: relative;            }            footer #page, footer #info {                position: absolute;            }            footer #page {                right: 8%;                border: 1pt solid #979797;                padding: 5pt 10pt;                top: -10pt;            }            footer #info {                left: 0;                bottom: 0;                font-size: 11pt;            }        p {            text-align: justify;        }`

const respJson = `{
    "Content": "JVBERi0xLjQKJdPr6eEKMSAwIG9iago8PC9DcmVhdG9yIChDaHJvbWl1bSkKL1Byb2R1Y2VyIChTa2lhL1BERiBtMTI3KQovQ3JlYXRpb25EYXRlIChEOjIwMjQwODA4MTg0NDU2KzAwJzAwJykKL01vZERhdGUgKEQ6MjAyNDA4MDgxODQ0NTYrMDAnMDAnKT4+CmVuZG9iagozIDAgb2JqCjw8L2NhIDEKL0JNIC9Ob3JtYWw+PgplbmRvYmoKNSAwIG9iago8PC9GaWx0ZXIgL0ZsYXRlRGVjb2RlCi9MZW5ndGggMjk3Pj4gc3RyZWFtCnictZTZagJBEEXf6yvqOZCyq6urFwiBmKjPhoZ8QBYhYEDz/xAcNRrIzQI6TzVcbq1nRqK14eHAgS/l6LUmlaatVX5c0ooC5yjO1kyKsxVxXj/TwwW/0YpMNPqQ4hA9LikMwf2Mt8F6QaOZ8eKdxp1G08QaJPumXOX+QrrtQdkiW+C+pKsQ1K65v5JLi9WLZQ7cn3gjjJEQByFKy5qrx4NwwlSOHLAGdBQghCkoHgx1lVCNioQ7lAp1tRMmneY0H7CwFJpYcY6mAyFnRkM16A6OhA7nBcyFBXg4uCIIBxTgfSAcN0jQf0N+QgecA+4KAgg5g5NDx+3PX/EnsntO99xqrmf/pZVUpMaqdozvbo+T/pvZTUr2/NUc/2rOUixZ/b7yZiUf77A/zgplbmRzdHJlYW0KZW5kb2JqCjIgMCBvYmoKPDwvVHlwZSAvUGFnZQovUmVzb3VyY2VzIDw8L1Byb2NTZXQgWy9QREYgL1RleHQgL0ltYWdlQiAvSW1hZ2VDIC9JbWFnZUldCi9FeHRHU3RhdGUgPDwv4NiAwMDAwMCBuIAowMDAwMDE0MDc0IDAwMDAwIG4gCnRyYWlsZXIKPDwvU2l6ZSAxMgovUm9vdCA3IDAgUgovSW5mbyAxIDAgUj4+CnN0YXJ0eHJlZgoxNDU3MQolJUVPRgo="
}`

func TestHandleHtml2Pdf(t *testing.T) {
	t.Parallel()

	t.Run("HandleHttp2Pdf", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)

		hc := controllers.NewHtml2PdfController(logger)

		path := "/html2pdf"
		r.POST(path, hc.HandleHttp2Pdf)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.Content = jsonContent
		obj.ContentCss = jsonContentCss

		body, err := json.Marshal(obj)
		if err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequest("POST", path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("version", "1")
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)

		assert.NotEqual(t, respJson, string(responseData))
	})

	t.Run("HandleHttp2Pdf404Json", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)

		hc := controllers.NewHtml2PdfController(logger)

		path := "/html2pdf"
		r.POST(path, hc.HandleHttp2Pdf)

		obj := dtos.HtmlRequest{}
		obj.HeaderTemplate = jsonHeader
		obj.FooterTemplate = jsonFooter
		obj.ContentCss = jsonContentCss

		body, err := json.Marshal(obj)
		if err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequest("POST", path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("version", "1")
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)

		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.NotEqual(t, respJson, string(responseData))
	})

}
