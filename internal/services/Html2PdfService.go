package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/kolzxx/html2pdf/internal/dtos"
	"github.com/kolzxx/html2pdf/internal/interfaces"
	"github.com/kolzxx/html2pdf/internal/logger"
)

type html2PdfService struct {
	logger          logger.Logger
	Content         string
	ContentCss      string
	WaitElementId   string
	chromedpService *ChromedpService
	lock            *sync.Mutex
}

func NewHtml2PdfService(l logger.Logger, chromedpService *ChromedpService) interfaces.Html2PdfServiceInterface {
	obj := &html2PdfService{
		logger:          l,
		chromedpService: chromedpService,
		lock:            &sync.Mutex{},
	}
	return obj
}

func (r *html2PdfService) HtmlToPdf(request dtos.HtmlRequest) (dtos.PdfResponse, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	resp := new(dtos.PdfResponse)
	r.WaitElementId = request.WaitElementId
	r.Content = request.Content
	r.Content = strings.ReplaceAll(r.Content, "\r\n", "\n")
	r.Content = strings.ReplaceAll(r.Content, "\r", "")
	r.ContentCss = request.ContentCss
	r.ContentCss = strings.ReplaceAll(r.ContentCss, "\r\n", "\n")
	r.ContentCss = strings.ReplaceAll(r.ContentCss, "\r", "")
	request.ContentCss = strings.ReplaceAll(request.ContentCss, "\r\n", "\n")
	request.ContentCss = strings.ReplaceAll(request.ContentCss, "\r", "")
	request.FooterTemplate = strings.ReplaceAll(request.FooterTemplate, "\r\n", "\n")
	request.HeaderTemplate = strings.ReplaceAll(request.HeaderTemplate, "\r\n", "\n")
	request.FooterTemplate = strings.ReplaceAll(request.FooterTemplate, "\r", "")
	request.HeaderTemplate = strings.ReplaceAll(request.HeaderTemplate, "\r", "")

	ts := httptest.NewServer(r.WriteHTML())

	defer ts.Close()

	taskCtx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		r.logger.Error("browser process isn't started", err)
	}

	cxtt, cancelt := context.WithTimeout(taskCtx, time.Second*20)
	defer cancelt()

	var pdfBuffer []byte

	if err := chromedp.Run(cxtt,
		r.PdfGrabber(ts.URL, &pdfBuffer, request)); err != nil {
		if len(pdfBuffer) > 0 {
			resp.Content = pdfBuffer
			return *resp, nil
		}
		if err2 := chromedp.Run(cxtt,
			r.PdfGrabber(ts.URL, &pdfBuffer, request),
		); err2 != nil {
			r.logger.Error("context timeout reached, attempting to perform actions", err)
		}
	}

	resp.Content = pdfBuffer

	return *resp, nil
}

const (
	makeVisibleScript = `setTimeout(function() { document.querySelector('#hash_assinatura').style.display = '';	}, 3000);`
)

func (r *html2PdfService) PdfGrabber(url string, res *[]byte, request dtos.HtmlRequest) chromedp.Tasks {
	var nodes []*cdp.Node
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, lcancel := context.WithCancel(ctx)
			defer lcancel()
			var wg sync.WaitGroup
			wg.Add(1)
			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					// It's a good habit to remove the event listener if we don't need it anymore.
					lcancel()
					wg.Done()
				}
			})
			_, exp, err := runtime.Evaluate(makeVisibleScript).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			r.logger.Info("Getting Frame Tree finish")
			if len(r.ContentCss) > 0 {
				if err := page.SetDocumentContent(frameTree.Frame.ID, fmt.Sprintf(r.Content, r.ContentCss)).Do(ctx); err != nil {
					return err
				}
			} else {
				if err := page.SetDocumentContent(frameTree.Frame.ID, r.Content).Do(ctx); err != nil {
					return err
				}
			}
			wg.Wait()

			return nil
		}),
		chromedp.Nodes("#"+r.WaitElementId, &nodes, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.ActionFunc(r.pdfActions(res, request)),
	}
}

func (r *html2PdfService) Wait(wg *sync.WaitGroup) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		r.WgWait(wg)
		return nil
	}
}

func (r *html2PdfService) WgWait(wg *sync.WaitGroup) {
	wg.Wait()
}

func (r *html2PdfService) pdfActions(res *[]byte, request dtos.HtmlRequest) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return r.DoPdfActions(res, request, ctx)
	}
}

func (r *html2PdfService) DoPdfActions(res *[]byte, request dtos.HtmlRequest, ctx context.Context) error {
	buf, err := doPrint(ctx, request)
	if err != nil {
		return err
	}
	*res = buf
	return nil
}

func (r *html2PdfService) getFrameTree() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return r.DoGetFrameTree(ctx)
	}
}

func (r *html2PdfService) DoGetFrameTree(ctx context.Context) error {
	r.logger.Info("Getting Frame Tree")
	frameTree, err := page.GetFrameTree().Do(ctx)
	if err != nil {
		return err
	}
	r.logger.Info("Getting Frame Tree finish")
	if len(r.ContentCss) > 0 {
		return page.SetDocumentContent(frameTree.Frame.ID, fmt.Sprintf(r.Content, r.ContentCss)).Do(ctx)
	} else {
		return page.SetDocumentContent(frameTree.Frame.ID, r.Content).Do(ctx)
	}

}

func doPrint(ctx context.Context, request dtos.HtmlRequest) ([]byte, error) {
	buf, _, err := page.PrintToPDF().
		WithDisplayHeaderFooter(request.DisplayHeaderFooter).
		WithPrintBackground(true).
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

func (r *html2PdfService) WriteHTML() http.Handler {
	return http.HandlerFunc(r.DoHandler)
}

func (r *html2PdfService) DoHandler(w http.ResponseWriter, h *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	io.WriteString(w, strings.TrimSpace(r.ContentCss))
}
