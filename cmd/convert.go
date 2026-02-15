package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fblissjr/crx2md/internal/crx"
	"github.com/fblissjr/crx2md/internal/extension"
	"github.com/fblissjr/crx2md/internal/markdown"
	"github.com/spf13/cobra"
)

var (
	outputPath    string
	includeBinary bool
)

var convertCmd = &cobra.Command{
	Use:   "convert <source>",
	Short: "Convert a Chrome extension to markdown",
	Long: `Convert a Chrome extension to an LLM-optimized markdown document.

Source can be:
  - A Chrome Web Store URL
  - A local .crx or .zip file path
  - A 32-character extension ID`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func init() {
	convertCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path (default: stdout)")
	convertCmd.Flags().BoolVar(&includeBinary, "include-binary", false, "list binary files with sizes instead of skipping")
	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	source := args[0]

	// Resolve the CRX file path
	crxPath, cleanup, err := resolveCRXPath(source)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	// Parse CRX to get ZIP payload
	result, err := crx.ParseFile(crxPath)
	if err != nil {
		return fmt.Errorf("parse CRX: %w", err)
	}

	// Extract extension files from ZIP
	ext, err := extension.Extract(result.Reader, result.Size)
	if err != nil {
		return fmt.Errorf("extract extension: %w", err)
	}

	// Render to markdown
	md := markdown.Render(ext, includeBinary)

	// Write output
	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(md), 0644); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Output written to %s\n", outputPath)
	} else {
		fmt.Print(md)
	}

	return nil
}

// resolveCRXPath determines the input type and returns a local file path to a CRX/ZIP file.
// Returns a cleanup function if a temp file was created (for downloads).
func resolveCRXPath(source string) (path string, cleanup func(), err error) {
	switch {
	case strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://"):
		// Chrome Web Store URL
		extID, err := crx.ParseExtensionID(source)
		if err != nil {
			return "", nil, err
		}
		return downloadExtension(extID)

	case strings.HasSuffix(source, ".crx") || strings.HasSuffix(source, ".zip"):
		// Local file
		if _, err := os.Stat(source); err != nil {
			return "", nil, fmt.Errorf("file not found: %s", source)
		}
		return source, nil, nil

	case len(source) == 32 && isAlphaLower(source):
		// Raw extension ID
		return downloadExtension(source)

	default:
		// Check if it's a directory (future: unpacked extension support)
		info, err := os.Stat(source)
		if err == nil && info.IsDir() {
			return "", nil, fmt.Errorf("unpacked extension directories not yet supported")
		}
		return "", nil, fmt.Errorf("unrecognized input: %s (expected URL, .crx/.zip file, or 32-char extension ID)", source)
	}
}

func downloadExtension(extID string) (string, func(), error) {
	fmt.Fprintf(os.Stderr, "Downloading extension %s...\n", extID)
	tmpPath, err := crx.Download(extID)
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.Remove(tmpPath) }
	return tmpPath, cleanup, nil
}

func isAlphaLower(s string) bool {
	for _, c := range s {
		if c < 'a' || c > 'z' {
			return false
		}
	}
	return true
}
