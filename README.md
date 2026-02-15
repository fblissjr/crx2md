# crx2md

Extract Chrome extension source code into a single markdown document. Useful for feeding to LLMs.

## Install

```bash
go install github.com/fblissjr/crx2md@latest
```

Or build from source:

```bash
go build -o crx2md .
```

To update after `git pull`:

```bash
go build -o crx2md .
```

## Usage

```bash
# From Chrome Web Store URL
crx2md convert https://chromewebstore.google.com/detail/unhook-remove-youtube-rec/khncfooichmfjbepaaaebmommgaepoid

# From extension ID
crx2md convert khncfooichmfjbepaaaebmommgaepoid

# From local .crx or .zip file
crx2md convert extension.crx

# Write to file instead of stdout
crx2md convert <source> -o output.md

# Include binary file listings
crx2md convert <source> --include-binary
```

## Output

Produces markdown with:
- Extension metadata (name, version, permissions)
- File tree
- All source files in fenced code blocks with language hints
- Binary files skipped (count reported at the bottom)

Supports CRX2, CRX3, and plain ZIP formats.
