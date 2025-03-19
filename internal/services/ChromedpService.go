package services

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/kolzxx/html2pdf/internal/logger"
)

const content = `
<body>
<script>
	// Show the current cookies.
	var p = document.createElement("p")
	p.innerText = document.cookie
	p.setAttribute("id", "cookies")
	document.body.appendChild(p)
	// Override the cookies.
	document.cookie = "foo=bar"
</script>
</body>
	`

type ChromedpService struct {
	logger  logger.Logger
	Context context.Context
	Cancelf context.CancelFunc
	Content string
}

func NewChromedpService(c context.Context, l logger.Logger) *ChromedpService {
	obj := &ChromedpService{
		logger: l,
	}

	return obj
}

func (c *ChromedpService) RunChromeDp() error {
	c.Content = content

	ts := httptest.NewServer(c.WriteHTML())

	defer ts.Close()

	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
		// chromedp.WithDebugf(log.Printf),
	)
	c.Cancelf = cancel
	c.Context = taskCtx
	chromedp.Flag("headless", false)
	chromedp.Flag("disable-gpu", false)
	chromedp.Flag("enable-automation", false)
	chromedp.Flag("disable-extensions", false)

	if err := chromedp.Run(taskCtx); err != nil {
		c.logger.Error(err.Error())
	}

	return nil
}

func (c *ChromedpService) WriteHTML() http.Handler {
	return http.HandlerFunc(c.DoHandler)
}

func (c *ChromedpService) DoHandler(w http.ResponseWriter, h *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, strings.TrimSpace(c.Content))
}
