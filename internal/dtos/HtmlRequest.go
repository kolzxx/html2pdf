package dtos

type HtmlRequest struct {
	PrintBackground     bool    `default:"false"`
	PreferCSSPageSize   bool    `default:"false"`
	DisplayHeaderFooter bool    `default:"true"`
	Landscape           bool    `default:"false"`
	MarginTop           float64 `default:"1.0"`
	MarginBottom        float64 `default:"1.0"`
	MarginRight         float64 `default:"0"`
	MarginLeft          float64 `default:"1.0"`
	PaperWidth          float64 `default:"8.27"`
	PaperHeight         float64 `default:"11.69"`
	WithScale           float64 `default:"0.57"`
	Content             string  `binding:"required" `
	ContentCss          string
	HeaderTemplate      string
	FooterTemplate      string
	WaitElementId       string `binding:"required" `
}
