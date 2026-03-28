package backup

import (
	"strings"
	"testing"
)

func TestGenerateQRCodesSingleChunk(t *testing.T) {
	data := []byte("short data for QR code")

	chunks, err := GenerateQRCodes(data)
	if err != nil {
		t.Fatalf("GenerateQRCodes: %v", err)
	}
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if chunks[0].Part != 1 || chunks[0].Total != 1 {
		t.Errorf("chunk = %d/%d, want 1/1", chunks[0].Part, chunks[0].Total)
	}
	if chunks[0].DataLen != len(data) {
		t.Errorf("DataLen = %d, want %d", chunks[0].DataLen, len(data))
	}
	svg := string(chunks[0].SVG)
	if !strings.Contains(svg, "<svg") {
		t.Error("SVG output should contain <svg element")
	}
}

func TestGenerateQRCodesMultiChunk(t *testing.T) {
	// Create data larger than maxPayload to force multi-QR
	data := make([]byte, maxPayload+500)
	for i := range data {
		data[i] = byte('A' + (i % 26))
	}

	chunks, err := GenerateQRCodes(data)
	if err != nil {
		t.Fatalf("GenerateQRCodes: %v", err)
	}
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	for i, c := range chunks {
		if c.Part != i+1 {
			t.Errorf("chunk %d: Part = %d, want %d", i, c.Part, i+1)
		}
		if c.Total != len(chunks) {
			t.Errorf("chunk %d: Total = %d, want %d", i, c.Total, len(chunks))
		}
	}
}
