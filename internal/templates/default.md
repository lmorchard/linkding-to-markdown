# {{ .Title }}

_Generated: {{ .Generated }}_

---

{{ if .Options.GroupByDate -}}
{{ range $date, $bookmarks := .GroupedBookmarks -}}
## {{ $date }}

{{ range $bookmarks -}}
### [{{ if .Title }}{{ .Title }}{{ else }}{{ .URL }}{{ end }}]({{ .URL }})

{{ if .Description }}{{ .Description }}{{ end }}

{{ if and $.Options.IncludeNotes (hasContent .Notes) -}}
**Notes:**

{{ .Notes }}

{{ end -}}

{{ if and $.Options.IncludeTags (gt (len .TagNames) 0) -}}
**Tags:** {{ join .TagNames ", " }}

{{ end -}}

{{ if .WebsiteTitle -}}
_Website: {{ .WebsiteTitle }}{{ if .WebsiteDescription }} - {{ .WebsiteDescription }}{{ end }}_

{{ end -}}

_Added: {{ formatDate .DateAdded "2006-01-02 15:04:05" }}{{ if .IsArchived }} | Archived{{ end }}{{ if .Unread }} | Unread{{ end }}{{ if .Shared }} | Shared{{ end }}_

---

{{ end -}}
{{ end -}}
{{ else -}}
{{ range .Bookmarks -}}
## [{{ if .Title }}{{ .Title }}{{ else }}{{ .URL }}{{ end }}]({{ .URL }})

{{ if .Description }}{{ .Description }}{{ end }}

{{ if and .Options.IncludeNotes (hasContent .Notes) -}}
**Notes:**

{{ .Notes }}

{{ end -}}

{{ if and .Options.IncludeTags (gt (len .TagNames) 0) -}}
**Tags:** {{ join .TagNames ", " }}

{{ end -}}

{{ if .WebsiteTitle -}}
_Website: {{ .WebsiteTitle }}{{ if .WebsiteDescription }} - {{ .WebsiteDescription }}{{ end }}_

{{ end -}}

_Added: {{ formatDate .DateAdded "2006-01-02 15:04:05" }}{{ if .IsArchived }} | Archived{{ end }}{{ if .Unread }} | Unread{{ end }}{{ if .Shared }} | Shared{{ end }}_

---

{{ end -}}
{{ end -}}
