package interfaces

import (
	"context"
	"net/http"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/kolzxx/html2pdf/internal/dtos"
)

type Html2PdfServiceInterface interface {
	HtmlToPdf(request dtos.HtmlRequest) (dtos.PdfResponse, error)
	PdfGrabber(url string, res *[]byte, request dtos.HtmlRequest) chromedp.Tasks
	DoPdfActions(res *[]byte, request dtos.HtmlRequest, ctx context.Context) error
	WriteHTML() http.Handler
	DoHandler(w http.ResponseWriter, h *http.Request)
	// DoEventLoad(ctx context.Context, wg *sync.WaitGroup, ch *chan int) error
	DoGetFrameTree(ctx context.Context) error
	Wait(wg *sync.WaitGroup) chromedp.ActionFunc
	WgWait(wg *sync.WaitGroup)
}
