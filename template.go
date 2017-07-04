package main

const tpl = `{{ if .Magic}}#### Magic
| Field       | Value                  |
|-------------|------------------------|
| Mime        | {{.Magic.Mime}}        |
| Description | {{.Magic.Description}} |
{{ end -}}
{{- if .SSDeep}}
#### SSDeep
 - ` + "`" + `{{.SSDeep}}` + "`" + `
{{ end -}}
{{- if .TRiD}}
#### TRiD
{{ range .TRiD -}}
 - {{ . }}
{{end}}
{{- end }}
{{- if .Exiftool}}
#### Exiftool
| Field       | Value                |
|-------------|----------------------|
{{- range $key, $value := .Exiftool }}
| {{ $key }}  | {{ $value }}        |
{{- end }}
{{- end }}
`
