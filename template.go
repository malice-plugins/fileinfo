package main

const tpl = `
{{- if .Matches}}
#### Magic
| Field       | Value            |
|-------------|------------------|
| Mime        | {{.Mime}}        |
| Description | {{.Description}} |
{{ end -}}
{{- if .SSDeep}}
#### SSDeep
 - {{.SSDeep}}
{{ end -}}
{{- if .TRiD}}
#### TRiD
{{ range .TRiD }}
 - {{.}}
{{end}}
{{ end -}}
{{- if .Exiftool}}
#### Exiftool
| Field       | Value            |
|-------------|------------------|
{{ range $key, $value := .Exiftool }}
| {{ $key }}  | {{ $value }}     |
{{end}}
{{ end -}}
`
