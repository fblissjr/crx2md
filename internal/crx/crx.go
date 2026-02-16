package crx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var crxMagic = []byte("Cr24")

// ParseResult holds the ZIP payload extracted from a CRX file.
type ParseResult struct {
	Reader io.ReaderAt
	Size   int64
}

// ParseFile opens a CRX or ZIP file and returns the ZIP payload.
func ParseFile(path string) (*ParseResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	offset, err := findZIPOffset(data)
	if err != nil {
		return nil, err
	}

	zipData := data[offset:]
	return &ParseResult{
		Reader: bytes.NewReader(zipData),
		Size:   info.Size() - int64(offset),
	}, nil
}

// findZIPOffset determines the byte offset where the ZIP payload begins.
// Supports CRX2, CRX3, and plain ZIP files.
func findZIPOffset(data []byte) (int, error) {
	// Check for plain ZIP (PK magic)
	if len(data) >= 2 && data[0] == 'P' && data[1] == 'K' {
		return 0, nil
	}

	// Check for CRX magic
	if len(data) < 12 {
		return 0, fmt.Errorf("file too small to be a CRX file")
	}
	if !bytes.Equal(data[:4], crxMagic) {
		return 0, fmt.Errorf("not a CRX or ZIP file (unrecognized magic bytes)")
	}

	version := binary.LittleEndian.Uint32(data[4:8])

	switch version {
	case 2:
		return parseCRX2Offset(data)
	case 3:
		return parseCRX3Offset(data)
	default:
		return 0, fmt.Errorf("unsupported CRX version: %d", version)
	}
}

// parseCRX2Offset calculates the ZIP offset for CRX2 format.
// Layout: magic(4) + version(4) + pubkey_len(4) + sig_len(4) + pubkey + sig + ZIP
func parseCRX2Offset(data []byte) (int, error) {
	if len(data) < 16 {
		return 0, fmt.Errorf("CRX2 file too small for header")
	}

	pubkeyLen := binary.LittleEndian.Uint32(data[8:12])
	sigLen := binary.LittleEndian.Uint32(data[12:16])
	offset := 16 + int(pubkeyLen) + int(sigLen)

	if offset >= len(data) {
		return 0, fmt.Errorf("CRX2 header exceeds file size")
	}
	return offset, nil
}

// parseCRX3Offset calculates the ZIP offset for CRX3 format.
// Layout: magic(4) + version(4) + header_len(4) + header + ZIP
func parseCRX3Offset(data []byte) (int, error) {
	if len(data) < 12 {
		return 0, fmt.Errorf("CRX3 file too small for header")
	}

	headerLen := binary.LittleEndian.Uint32(data[8:12])
	offset := 12 + int(headerLen)

	if offset >= len(data) {
		return 0, fmt.Errorf("CRX3 header exceeds file size")
	}
	return offset, nil
}
