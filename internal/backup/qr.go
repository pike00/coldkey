package backup

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	qr "github.com/piglig/go-qr"
)

// QRChunk represents one part of a potentially multi-QR encoding.
type QRChunk struct {
	Part    int
	Total   int
	SVG     template.HTML
	DataLen int
}

// maxPayload is the conservative max bytes per QR code (v40, EC-H = 1273 minus header overhead).
const maxPayload = 1250

// GenerateQRCodes encodes data into one or more QR code SVGs.
// If data fits in a single QR code, no framing header is added.
// If splitting is needed, each chunk is prefixed with COLDKEY:<part>/<total>: for reassembly.
func GenerateQRCodes(data []byte) ([]QRChunk, error) {
	if len(data) <= maxPayload {
		svg, err := encodeQR(data)
		if err != nil {
			return nil, err
		}
		return []QRChunk{{
			Part:    1,
			Total:   1,
			SVG:     template.HTML(svg),
			DataLen: len(data),
		}}, nil
	}

	// Multi-QR: split with framing headers
	// Header like "COLDKEY:01/03:" is 15 bytes max
	chunkData := maxPayload - 20 // conservative header allowance
	total := (len(data) + chunkData - 1) / chunkData

	chunks := make([]QRChunk, 0, total)
	for i := 0; i < total; i++ {
		start := i * chunkData
		end := start + chunkData
		if end > len(data) {
			end = len(data)
		}
		payload := fmt.Sprintf("COLDKEY:%d/%d:", i+1, total) + string(data[start:end])

		svg, err := encodeQR([]byte(payload))
		if err != nil {
			return nil, fmt.Errorf("encoding QR part %d/%d: %w", i+1, total, err)
		}
		chunks = append(chunks, QRChunk{
			Part:    i + 1,
			Total:   total,
			SVG:     template.HTML(svg),
			DataLen: end - start,
		})
	}
	return chunks, nil
}

// MaxQRCapacity is the byte limit for QR version 40 at error correction level H.
const MaxQRCapacity = 1273

func encodeQR(data []byte) (string, error) {
	seg, err := qr.MakeBytes(data)
	if err != nil {
		return "", fmt.Errorf("creating QR segment: %w", err)
	}
	code, err := qr.EncodeSegments([]*qr.QrSegment{seg}, qr.High, 1, 40, -1, true)
	if err != nil {
		return "", fmt.Errorf("encoding QR: %w", err)
	}
	config := qr.NewQrCodeImgConfig(4, 0)
	var buf bytes.Buffer
	if err := code.WriteAsSVG(config, &buf, "#FFFFFF", "#000000"); err != nil {
		return "", fmt.Errorf("rendering QR SVG: %w", err)
	}
	// Strip XML declaration for embedding in HTML
	svg := buf.String()
	if idx := strings.Index(svg, "<svg"); idx > 0 {
		svg = svg[idx:]
	}
	return svg, nil
}
