package backup

import (
	"embed"
	"html/template"
)

//go:embed templates/backup.html.tmpl
var templateFS embed.FS

// TemplateData is passed to the HTML template for rendering.
type TemplateData struct {
	Date          string
	FileSize      int64
	RawKeyContent string
	SHA256        string
	QRChunks      []QRChunk
	TotalQRParts  int
	MaxQRCapacity int
	QRUsagePct    int
	Version       string
}

var backupTemplate = template.Must(
	template.ParseFS(templateFS, "templates/backup.html.tmpl"),
)
