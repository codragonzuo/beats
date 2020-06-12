// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package asset

import (
	"bytes"
	"go/format"
	"math"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type Priority int32

const (
	Highest          Priority = 1
	ECSFieldsPri     Priority = 5
	LibbeatFieldsPri Priority = 10
	BeatFieldsPri    Priority = 50
	ModuleFieldsPri  Priority = 100
	Lowest           Priority = math.MaxInt32
)

var Template = template.Must(template.New("normalizations").Parse(`
{{ .License }}

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package {{ .Package }}

import (
	"github.com/codragonzuo/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("{{ .Beat }}", "{{ .Name }}", {{ .Priority }}, Asset{{ .GoTypeName }}); err != nil {
		panic(err)
	}
}

// Asset{{ .GoTypeName }} returns asset data.
// This is the base64 encoded gzipped contents of {{ .Path }}.
func Asset{{ .GoTypeName }}() string {
	return "{{ .Data }}"
}

`))

type Data struct {
	License    string
	Beat       string
	Name       string
	Priority   string
	Data       string
	Package    string
	Path       string
	GoTypeName string
}

func CreateAsset(license string, beat string, name string, pkg string, data []byte, priority string, path string) ([]byte, error) {

	// Depending on OS or tools configuration, files can contain carriages (\r),
	// what leads to different results, remove them before encoding.
	encData, err := EncodeData(strings.Replace(string(data), "\r", "", -1))
	if err != nil {
		return nil, errors.Wrap(err, "error encoding the data")
	}

	goTypeName := goTypeName(name)
	var buf bytes.Buffer
	Template.Execute(&buf, Data{
		License:    license,
		Beat:       beat,
		Name:       name,
		Data:       encData,
		Priority:   priority,
		Package:    pkg,
		Path:       path,
		GoTypeName: goTypeName,
	})

	bs, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "error creating golang file from template")
	}

	return bs, nil
}

// goTypeName removes special characters ('_', '.', '@') and returns a
// camel-cased name.
func goTypeName(name string) string {
	var b strings.Builder
	for _, w := range strings.FieldsFunc(name, isSeparator) {
		b.WriteString(strings.Title(w))
	}
	return b.String()
}

// isSeparate returns true if the character is a field name separator. This is
// used to detect the separators in fields like ephemeral_id or instance.name.
func isSeparator(c rune) bool {
	switch c {
	case '.', '_', '/':
		return true
	case '@':
		// This effectively filters @ from field names.
		return true
	default:
		return false
	}
}
