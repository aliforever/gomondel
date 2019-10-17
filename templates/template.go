package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"unicode"

	"github.com/go-errors/errors"

	"github.com/gobuffalo/flect"
)

var ModelsPath = "%s/models"

type Template struct {
}

func (t *Template) Init(dbName string) (path string, err error) {
	path, err = t.CurrentPath()
	if err != nil {
		return
	}
	fileString := t.init(dbName)
	if _, err = os.Stat(fmt.Sprintf(ModelsPath, path)); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(fmt.Sprintf(ModelsPath, path), os.ModePerm)
			if err != nil {
				return
			}
		} else {
			return
		}
	}
	path = fmt.Sprintf(ModelsPath, path) + "/init.go"
	err = ioutil.WriteFile(path, []byte(fileString), os.ModePerm)
	if err == nil {
		err = t.GoFmtCurrentPath()
	}
	return
}

func (t *Template) CreateModel(modelName string, parentName *string) (path string, err error) {
	path, err = t.CurrentPath()
	if err != nil {
		return
	}
	fileName := flect.Singularize(modelName)
	fileName = flect.Capitalize(fileName)
	modelSign := ""
	r := []rune(fileName)
	for i := 0; i < len(r); i++ {
		ch := r[i]
		if unicode.IsUpper(ch) {
			modelSign += string(unicode.ToLower(ch))
		}
	}
	parentMethodStr := ""
	parentField := ""
	if parentName != nil {
		parent := fmt.Sprintf(ModelsPath, path) + "/" + *parentName + ".go"
		if _, err = os.Stat(parent); err != nil {
			err = errors.New(fmt.Sprintf("%s %s", parent, err.Error()))
			return
		} else {
			parentModelName := flect.Singularize(*parentName)
			parentModelName = flect.Capitalize(parentModelName)
			parentModelSign := ""
			r := []rune(parentModelName)
			for i := 0; i < len(r); i++ {
				ch := r[i]
				if unicode.IsUpper(ch) {
					parentModelSign += string(unicode.ToLower(ch))
				}
			}
			var tpl, tplField bytes.Buffer
			type ParentData struct {
				ModelSign       string
				ModelName       string
				ParentModelName string
				ParentModelSign string
			}
			parentMethod := t.parentMethod()
			data := ParentData{ModelName: fileName, ModelSign: modelSign, ParentModelName: parentModelName, ParentModelSign: parentModelSign}
			var tmpl *template.Template
			tmpl, err = template.New("model").Parse(parentMethod)
			if err != nil {
				return
			}
			err = tmpl.Execute(&tpl, data)
			if err != nil {
				return
			}
			parentMethodStr = tpl.String()
			type ParentFieldData struct {
				ParentModelName      string
				ParentModelNameSmall string
			}
			field := t.parentField()
			fieldData := ParentFieldData{ParentModelName: parentModelName, ParentModelNameSmall: strings.ToLower(flect.Underscore(parentModelName))}
			tmpl, err = template.New("field").Parse(field)
			if err != nil {
				return
			}
			err = tmpl.Execute(&tplField, fieldData)
			if err != nil {
				return
			}
			parentField = tplField.String()
		}
	}
	fileString := t.model()
	if _, err = os.Stat(fmt.Sprintf(ModelsPath, path)); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(fmt.Sprintf(ModelsPath, path), os.ModePerm)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	path = fmt.Sprintf(ModelsPath, path) + "/" + strings.ToLower(fileName) + ".go"
	type Data struct {
		FileName             string
		ModelName            string
		ModelSign            string
		TableName            string
		ParentMethod         string
		ParentModelName      string
		ParentModelNameSmall string
		ParentField          string
	}

	tableName := strings.ToLower(flect.Underscore(flect.Pluralize(fileName)))
	data := Data{FileName: fileName, ModelName: fileName, ModelSign: modelSign, TableName: tableName, ParentMethod: parentMethodStr, ParentField: parentField}
	tmpl, err := template.New("model").Parse(fileString)
	if err != nil {
		return
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, tpl.Bytes(), os.ModePerm)
	if err == nil {
		err = t.GoFmtCurrentPath()
	}
	return
}

func (t *Template) GoFmtCurrentPath() (err error) {
	var path string
	path, err = t.CurrentPath()
	if err != nil {
		return
	}
	cmd := exec.Command("go", "fmt", path+"/...")
	err = cmd.Run()
	return
}

func (t Template) CurrentPath() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}
	return
}
