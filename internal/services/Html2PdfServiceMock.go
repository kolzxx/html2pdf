package services

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/kolzxx/html2pdf/internal/dtos"
	"github.com/kolzxx/html2pdf/internal/logger"
)

type Html2PdfServiceMock struct {
	logger  logger.Logger
	Content string
}

func NewHtml2PdfServiceMock(l logger.Logger) *Html2PdfServiceMock {
	obj := &Html2PdfServiceMock{
		logger: l,
	}
	return obj
}

func (r *Html2PdfServiceMock) HtmlToPdf(request dtos.HtmlRequest) (dtos.PdfResponse, error) {
	resp := new(dtos.PdfResponse)
	r.Content = request.Content

	ts := httptest.NewServer(r.WriteHTML())

	defer ts.Close()

	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var pdfBuffer []byte

	if err := chromedp.Run(taskCtx, r.PdfGrabber(ts.URL, &pdfBuffer, request)); err != nil {
		r.logger.Error(err.Error())
	}

	resp.Content = pdfBuffer

	return *resp, nil
}

func (r *Html2PdfServiceMock) PdfGrabber(url string, res *[]byte, request dtos.HtmlRequest) chromedp.Tasks {
	return chromedp.Tasks{
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(url),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.ScrollIntoView(`body`),
		// chromedp.Text(`h1`, &res, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.ActionFunc(r.PdfActions(res, request)),
	}
}

func (r *Html2PdfServiceMock) PdfActions(res *[]byte, request dtos.HtmlRequest) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return r.DoPdfActions(res, request, ctx)
	}
}
func (r *Html2PdfServiceMock) DoPdfActions(res *[]byte, request dtos.HtmlRequest, ctx context.Context) error {
	buf, err := doPrint(ctx, request)
	if err != nil {
		return err
	}
	*res = buf
	return nil
}

func doPrintMock(ctx context.Context, request dtos.HtmlRequest) ([]byte, error) {
	buf, _, err := page.PrintToPDF().
		WithDisplayHeaderFooter(request.DisplayHeaderFooter).
		WithPrintBackground(request.PrintBackground).
		WithPreferCSSPageSize(request.PreferCSSPageSize).
		WithScale(request.WithScale).
		WithPaperWidth(request.PaperWidth).
		WithPaperHeight(request.PaperHeight).
		WithLandscape(request.Landscape).
		WithMarginTop(request.MarginTop).
		WithMarginRight(request.MarginRight).
		WithMarginBottom(request.MarginBottom).
		WithMarginLeft(request.MarginLeft).
		WithHeaderTemplate(request.HeaderTemplate).
		WithFooterTemplate(request.FooterTemplate).
		Do(ctx)
	return buf, err
}

func (r *Html2PdfServiceMock) WriteHTML() http.Handler {
	return http.HandlerFunc(r.DoHandler)
}

func (r *Html2PdfServiceMock) DoHandler(w http.ResponseWriter, h *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, strings.TrimSpace(r.Content))
}
