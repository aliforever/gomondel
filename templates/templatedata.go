package templates

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gobuffalo/flect"
)

type TemplateData struct {
	ModelSign            string
	ModelName            string
	ParentKeyType        string
	ParentModelName      string
	ParentModelNameSmall string
	ParentModelSign      string
	FileName             string
	TableName            string
	ParentMethod         string
	ParentField          string
}

func (td TemplateData) FillModelParentMethod(content string, modelSign, modelName, parentModelName, parentModelSign string) (result string, err error) {
	td.ModelSign = modelSign
	td.ModelName = modelName
	td.ParentModelName = parentModelName
	td.ParentModelSign = parentModelSign
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("model_parent_method").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) FillModelParentField(content string, parentModelName, parentKeyType string) (result string, err error) {
	td.ParentKeyType = parentKeyType
	if td.ParentKeyType == "" {
		td.ParentKeyType = "primitive.ObjectID"
	}
	td.ParentModelName = parentModelName
	td.ParentModelNameSmall = strings.ToLower(flect.Underscore(parentModelName))
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)
	tmpl, err = template.New("model_parent_field").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) FillModel(content string, fileName, modelName, modelSign, tableName, parentMethod, parentField string) (result string, err error) {
	td.FileName = fileName
	td.ModelName = fileName
	td.ModelSign = modelSign
	td.TableName = tableName
	td.ParentMethod = parentMethod
	td.ParentField = parentField
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)
	tmpl, err = template.New("model_parent_field").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}
