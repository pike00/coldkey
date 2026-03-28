package backup

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pike00/coldkey/internal/keyfile"
	"github.com/pike00/coldkey/internal/secure"
)

// Generate produces an HTML paper backup from a parsed key file.
func Generate(ki *keyfile.KeyInfo, version string) ([]byte, error) {
	// Generate QR codes from private key only (public key can be regenerated)
	secretKeyData := []byte(ki.SecretKey)
	chunks, err := GenerateQRCodes(secretKeyData)
	if err != nil {
		return nil, fmt.Errorf("generating QR codes: %w", err)
	}

	usagePct := 0
	if len(chunks) == 1 {
		usagePct = len(secretKeyData) * 100 / MaxQRCapacity
	}

	data := TemplateData{
		Date:          time.Now().Format("2006-01-02"),
		FileSize:      ki.FileSize,
		RawKeyContent: string(ki.RawContent),
		SHA256:        ki.SHA256,
		QRChunks:      chunks,
		TotalQRParts:  len(chunks),
		MaxQRCapacity: MaxQRCapacity,
		QRUsagePct:    usagePct,
		Version:       version,
	}

	var buf bytes.Buffer
	if err := backupTemplate.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering template: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteHTML generates and writes the backup HTML to the given path.
func WriteHTML(ki *keyfile.KeyInfo, outputPath, version string) error {
	html, err := Generate(ki, version)
	if err != nil {
		return err
	}
	return secure.WriteFile(outputPath, html)
}
