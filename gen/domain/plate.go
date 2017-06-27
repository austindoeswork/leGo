package domain

var Plate = map[string]string{
	"Domain": `
// go run gen/gen.go
// VERY GENERATED PLZ NO MODIFY

// Package {{ .Name.UpperCamel }}
// {{ .Description }}
package {{ .Name.Lower }}

import (
	"fmt"
	"encoding/json"
	{{ range $k, $v := .Imports -}}
	"{{ $v }}"
	{{- end }}
	
	"git.ottoq.com/otto-backend/valet/entity"
)

type {{ .Name.UpperCamel }} struct {
	{{- range $p := .Parameters }}
	{{ $p.Name.UpperCamel }} {{ $p.Type }}
	{{- end }}
}

func New(
  {{ range $i, $param := .Parameters -}}
  {{ if ne .ConstructorOverride "" }}
  {{- else -}}
  {{ $param.Name.LowerCamel }} {{ $param.Type }},
  {{ end }} 
  {{- end -}}
) (*{{ .Name.UpperCamel }}, error) {
	d := &{{ .Name.UpperCamel }} {
	  {{ range $i, $param := .Parameters -}}
	  {{ $param.Name.UpperCamel }}:
	  {{- if ne .ConstructorOverride "" -}}
	  {{ .ConstructorOverride }},
	  {{- else -}}
	  {{ $param.Name.LowerCamel }},
	  {{- end }}
	  {{ end }}
	}
	return d, nil
}

type Scannable interface {
	Scan(dest ...interface{}) error
}

func NewFromRow(row Scannable) (*{{ .Name.UpperCamel }}, error) {
	d := {{ .Name.UpperCamel }}{}
	err := row.Scan(
	  {{ range $i, $param := .Parameters -}}
	  &d.{{ $param.Name.UpperCamel }},
	  {{ end }}
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func Schema() string {
	return ` + "`" + `{{ .SQLSchema }} ` + "`" + `
}

func TableName() string {
	return "{{ .Name.UpperCamel }}"
}

func Random() *{{ .Name.UpperCamel }} {
	d := &{{ .Name.UpperCamel }} {
	  {{ range $i, $param := .Parameters -}}
	  {{ $param.Name.UpperCamel }}:
	  {{- if ne .ConstructorOverride "" -}}
	  {{ .ConstructorOverride }},
	  {{- else -}}
	  entity.RAND{{ $param.Type }}(),
	  {{- end }}
	  {{ end }}
	}
	return d
}

func (o *{{ .Name.UpperCamel }}) InsertString() string {
	istr := fmt.Sprintf(` + "`" + `{{ .SQLInsert }}` + "`" + `,
	  {{ range $i, $param := .Parameters -}}
	  {{ if contains $param.Type.String "time" -}}
	  o.{{ $param.Name.UpperCamel }}.Format("2006-01-02 15:04:05"),
	  {{- else -}}
	  o.{{ $param.Name.UpperCamel }},
	  {{- end }}
	  {{ end }}
	  {{ range $i, $param := .Parameters -}}
	  {{ if not $param.PrimaryKey -}}
	  {{ if contains $param.Type.String "time" -}}
	  o.{{ $param.Name.UpperCamel }}.Format("2006-01-02 15:04:05"),
	  {{- else -}}
	  o.{{ $param.Name.UpperCamel }},
	  {{- end }}
	  {{- end }}
	  {{ end }}
	)
	return istr	
}

func (o *{{ .Name.UpperCamel }}) String() string {
	b, _ := json.MarshalIndent(o, "", "    ")
	return string(b)	
}

func (o *{{ .Name.UpperCamel }}) PPrint() {
	fmt.Println(o.String())
}
`,
}
