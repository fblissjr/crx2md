# Changelog

## 0.1.0

### Added
- Initial release
- `convert` command: convert Chrome extensions to LLM-optimized markdown
- Support for CRX2 and CRX3 file formats
- Support for plain .zip extension files
- Download extensions from Chrome Web Store URLs or extension IDs
- Auto-detection of input type (URL, file path, extension ID)
- File tree and metadata in output for LLM structural context
- Language-hinted fenced code blocks for all source files
- `--output` / `-o` flag for file output (default: stdout)
- `--include-binary` flag to list binary files with sizes
