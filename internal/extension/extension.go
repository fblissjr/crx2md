package extension

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

// Extension represents a parsed Chrome extension with all its source files.
type Extension struct {
	Manifest Manifest
	Files    []File
}

// File represents a single file within the extension.
type File struct {
	Path     string // Relative path within extension
	Content  []byte
	Language string // Language identifier for fenced code blocks
	Size     int64
	IsCode   bool // true for text files we should render
}

// Manifest holds parsed manifest.json fields plus the raw JSON.
type Manifest struct {
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	ManifestVersion int      `json:"manifest_version"`
	Description     string   `json:"description"`
	Permissions     []string `json:"permissions"`
	Raw             json.RawMessage
}

// codeExtensions maps file extensions to language identifiers for fenced code blocks.
var codeExtensions = map[string]string{
	".js":   "javascript",
	".ts":   "typescript",
	".jsx":  "javascript",
	".tsx":  "typescript",
	".css":  "css",
	".html": "html",
	".htm":  "html",
	".json": "json",
	".xml":  "xml",
	".svg":  "svg",
	".md":   "markdown",
	".txt":  "text",
	".yaml": "yaml",
	".yml":  "yaml",
}

// binaryExtensions lists file extensions that should be skipped.
var binaryExtensions = map[string]bool{
	".png":   true,
	".jpg":   true,
	".jpeg":  true,
	".gif":   true,
	".webp":  true,
	".ico":   true,
	".woff":  true,
	".woff2": true,
	".ttf":   true,
	".otf":   true,
	".eot":   true,
	".mp3":   true,
	".wav":   true,
	".wasm":  true,
	".map":   true,
}

// ClassifyFile determines whether a file is code or binary and returns
// the appropriate language identifier.
func ClassifyFile(path string) (language string, isCode bool) {
	ext := strings.ToLower(filepath.Ext(path))

	if lang, ok := codeExtensions[ext]; ok {
		return lang, true
	}

	if binaryExtensions[ext] {
		return "", false
	}

	// Unknown extension - try to include it as plain text
	return "text", true
}
