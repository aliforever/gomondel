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
	ModelNameChild       string
	ModelNameChildPlural string
	CustomModelId        bool
	ModelId              string
	ModelIdDefaultValue  string
	ModelIdSmall         string
	ModelIdType          string
	ParentModelIdType    string
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

func (td TemplateData) FillModelParentField(content string, parentModelName, parentIdType string) (result string, err error) {
	td.ParentModelIdType = parentIdType
	if td.ParentModelIdType == "" {
		td.ParentModelIdType = "primitive.ObjectID"
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

func (td TemplateData) FillModel(content string, fileName, modelName, modelSign, modelId, modelIdType, tableName, parentMethod, parentField string) (result string, err error) {
	td.ModelId = "Id"
	td.ModelIdDefaultValue = "primitive.NewObjectID()"
	td.ModelIdType = "primitive.ObjectID"
	if modelIdType != "" && modelIdType != "primitive.ObjectID" {
		td.CustomModelId = true
		td.ModelIdSmall = strings.ToLower(flect.Underscore(td.ModelId))
		td.ModelIdDefaultValue = td.ModelIdSmall
		td.ModelIdType = modelIdType
	}
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

func (td TemplateData) FillParentChildMethods(content string, fileName, modelName, modelSign, modelIdType, parentModelSign, parentModelName string) (result string, err error) {
	td.FileName = fileName
	td.ModelName = modelName
	td.ModelNameChild = flect.Underscore(td.ModelName)
	split := strings.SplitN(td.ModelNameChild, "_", 2)
	if len(split) == 2 {
		td.ModelNameChild = flect.Pascalize(split[1])
	}
	td.ModelNameChildPlural = flect.Pascalize(flect.Pluralize(td.ModelNameChild))
	td.ModelSign = modelSign
	td.ParentModelSign = parentModelSign
	td.ParentModelName = parentModelName
	td.ModelIdType = modelIdType
	if td.ModelIdType == "" {
		td.ModelIdType = "primitive.ObjectID"
	}
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)
	tmpl, err = template.New("parent_child_methods").Parse(content)
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
