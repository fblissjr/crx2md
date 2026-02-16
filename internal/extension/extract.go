package extension

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// Extract reads a ZIP archive and builds an Extension model from its contents.
func Extract(reader io.ReaderAt, size int64) (*Extension, error) {
	zr, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	ext := &Extension{}
	var manifestRaw []byte

	for _, zf := range zr.File {
		// Skip directories
		if zf.FileInfo().IsDir() {
			continue
		}

		path := zf.Name
		content, err := readZipFile(zf)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}

		language, isCode := ClassifyFile(path)

		file := File{
			Path:     path,
			Content:  content,
			Language: language,
			Size:     int64(zf.UncompressedSize64),
			IsCode:   isCode,
		}
		ext.Files = append(ext.Files, file)

		if path == "manifest.json" {
			manifestRaw = content
		}
	}

	// Parse manifest
	if manifestRaw != nil {
		if err := json.Unmarshal(manifestRaw, &ext.Manifest); err != nil {
			return nil, fmt.Errorf("parse manifest.json: %w", err)
		}
		ext.Manifest.Raw = manifestRaw
	}

	// Sort: manifest.json first, then alphabetically
	sort.Slice(ext.Files, func(i, j int) bool {
		if ext.Files[i].Path == "manifest.json" {
			return true
		}
		if ext.Files[j].Path == "manifest.json" {
			return false
		}
		return ext.Files[i].Path < ext.Files[j].Path
	})

	return ext, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	rc, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}
