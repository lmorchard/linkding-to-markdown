# linkding-to-markdown

A command-line tool to fetch bookmarks from [Linkding](https://github.com/sissbruecker/linkding) and export them to Markdown format.

## Features

- Fetch bookmarks from a Linkding instance over a specified timespan
- Filter bookmarks by date range (relative or absolute)
- Search bookmarks with query filters
- Customize output format (include/exclude notes, tags, grouping)
- Support for both configuration files and command-line flags
- Output to file or stdout

## Installation

### From Source

```bash
git clone https://github.com/lmorchard/linkding-to-markdown.git
cd linkding-to-markdown
make build
```

The binary will be created as `linkding-to-markdown` in the current directory.

### Using Go Install

```bash
go install github.com/lmorchard/linkding-to-markdown@latest
```

## Quick Start

Initialize the configuration and template files:

```bash
linkding-to-markdown init
```

This creates:
- `linkding-to-markdown.yaml` - Configuration file
- `linkding-to-markdown.md` - Customizable template

Edit the config file with your Linkding URL and API token, then fetch bookmarks:

```bash
linkding-to-markdown fetch --output bookmarks.md
```

## Configuration

You can configure the tool using either:

1. A YAML configuration file (default: `linkding-to-markdown.yaml`)
2. Environment variables (prefixed with `LINKDING_TO_MARKDOWN_`)
3. Command-line flags (highest precedence)

### Initialize Configuration

The easiest way to get started is to use the `init` command:

```bash
linkding-to-markdown init
```

This generates:
- `linkding-to-markdown.yaml` - Configuration file with all options documented
- `linkding-to-markdown.md` - Customizable Markdown template

You can specify a custom template filename:

```bash
linkding-to-markdown init --template-file my-template.md
```

Use `--force` to overwrite existing files:

```bash
linkding-to-markdown init --force
```

### Manual Configuration

Alternatively, copy the example configuration:

```bash
cp linkding-to-markdown.yaml.example linkding-to-markdown.yaml
```

Edit the file with your Linkding instance details:

```yaml
linkding:
  url: "https://your-linkding-instance.com"
  token: "your-api-token"
  timeout: 30s

fetch:
  days: 7
  title: "My Bookmarks"
  template: "linkding-to-markdown.md"  # Optional custom template
```

### Environment Variables

```bash
export LINKDING_TO_MARKDOWN_LINKDING_URL="https://your-linkding-instance.com"
export LINKDING_TO_MARKDOWN_LINKDING_TOKEN="your-api-token"
```

## Usage

### Basic Usage

Fetch bookmarks from the last 7 days:

```bash
linkding-to-markdown fetch --url https://linkding.example.com --token your-token
```

### Fetch with Date Range

Fetch bookmarks from a specific date range:

```bash
linkding-to-markdown fetch \
  --url https://linkding.example.com \
  --token your-token \
  --since 2025-01-01 \
  --until 2025-01-31
```

### Fetch with Search Query

Filter bookmarks by search query:

```bash
linkding-to-markdown fetch \
  --url https://linkding.example.com \
  --token your-token \
  --query "golang" \
  --days 30
```

### Save to File

Output to a file instead of stdout:

```bash
linkding-to-markdown fetch \
  --url https://linkding.example.com \
  --token your-token \
  --output bookmarks.md
```

### Customize Output Format

```bash
linkding-to-markdown fetch \
  --url https://linkding.example.com \
  --token your-token \
  --title "Weekly Bookmarks" \
  --no-notes \
  --no-group-by-date
```

### Use Custom Template

Use a custom Markdown template for output:

```bash
linkding-to-markdown fetch \
  --url https://linkding.example.com \
  --token your-token \
  --template linkding-to-markdown.md \
  --output bookmarks.md
```

## Custom Templates

The tool uses Go's `text/template` system for generating Markdown. You can create custom templates to control the exact output format.

### Template Variables

Templates have access to the following data:

- `.Title` - Document title (string)
- `.Generated` - Generation timestamp (string, RFC3339 format)
- `.Bookmarks` - Array of all bookmarks
- `.GroupedBookmarks` - Map of date strings to bookmark arrays (when grouping by date)
- `.Options` - Generation options (IncludeNotes, IncludeTags, GroupByDate, DateFormat)

### Bookmark Fields

Each bookmark has these fields:

- `.ID` - Bookmark ID (int)
- `.URL` - Bookmark URL (string)
- `.Title` - User-set title (string)
- `.Description` - User-set description (string)
- `.Notes` - User notes (string)
- `.WebsiteTitle` - Auto-fetched website title (string)
- `.WebsiteDescription` - Auto-fetched website description (string)
- `.IsArchived` - Whether bookmark is archived (bool)
- `.Unread` - Whether bookmark is unread (bool)
- `.Shared` - Whether bookmark is shared (bool)
- `.TagNames` - Array of tag names ([]string)
- `.DateAdded` - When bookmark was added (time.Time)
- `.DateModified` - When bookmark was last modified (time.Time)

### Template Functions

Available template functions:

- `formatDate <time> <format>` - Format a time.Time value (uses Go time format)
- `join <slice> <separator>` - Join string slice with separator
- `hasContent <string>` - Check if string has non-whitespace content

### Example Template

```markdown
# {{ .Title }}

{{ range .Bookmarks -}}
## [{{ .Title }}]({{ .URL }})

{{ .Description }}

Tags: {{ join .TagNames ", " }}

---
{{ end }}
```

## Command-Line Options

### Global Flags

- `--config` - Configuration file path (default: `./linkding-to-markdown.yaml`)
- `--verbose, -v` - Enable verbose output
- `--debug` - Enable debug output
- `--log-json` - Output logs in JSON format

### Fetch Command Flags

**Connection:**
- `--url` - Linkding instance URL (required)
- `--token` - Linkding API token (required)
- `--timeout` - HTTP request timeout (default: 30s)

**Time Range:**
- `--days` - Number of days to fetch from now (default: 7)
- `--since` - Fetch bookmarks added since this date (YYYY-MM-DD)
- `--until` - Fetch bookmarks added until this date (YYYY-MM-DD)

**Filtering:**
- `--query` - Search query to filter bookmarks

**Output:**
- `--output, -o` - Output file (default: stdout)
- `--title` - Title for the markdown document (default: "Bookmarks")
- `--template` - Custom template file (default: built-in template)
- `--no-notes` - Exclude notes from output
- `--no-tags` - Exclude tags from output
- `--no-group-by-date` - Don't group bookmarks by date
- `--date-format` - Date format for grouping (Go time format, default: "2006-01-02")

## Output Format

The generated Markdown includes:

- Document title and generation timestamp
- Bookmarks grouped by date (optional)
- For each bookmark:
  - Title (linked to URL)
  - Description
  - Notes (optional)
  - Tags (optional)
  - Website metadata
  - Added date and status flags (archived, unread, shared)

### Example Output

```markdown
# Bookmarks

_Generated: 2025-11-12T10:30:00Z_

---

## 2025-11-12

### [Example Bookmark](https://example.com)

A description of the bookmark.

**Notes:**

My personal notes about this bookmark.

**Tags:** golang, programming, tutorial

_Website: Example Site - An example website_

_Added: 2025-11-12 10:15:30 | Unread_

---
```

## Development

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Formatting

```bash
make format
```

## License

MIT License - see LICENSE file for details

## See Also

- [Linkding](https://github.com/sissbruecker/linkding) - Self-hosted bookmark manager
- [mastodon-to-markdown](https://github.com/lmorchard/mastodon-to-markdown) - Similar tool for Mastodon
- [linkding-to-opml](https://github.com/lmorchard/linkding-to-opml) - Export Linkding bookmarks to OPML
