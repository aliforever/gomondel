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

func (t *Template) CreateModel(modelName string) (path string, err error) {
	path, err = t.CurrentPath()
	if err != nil {
		return
	}
	fileString := t.createModel(modelName)
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
	fileName := flect.Singularize(modelName)
	fileName = flect.Capitalize(modelName)
	path = fmt.Sprintf(ModelsPath, path) + "/" + strings.ToLower(fileName) + ".go"
	type Data struct {
		FileName  string
		ModelName string
		ModelSign string
		TableName string
	}

	modelSign := ""
	r := []rune(fileName)
	for i := 0; i < len(r); i++ {
		ch := r[i]
		if unicode.IsUpper(ch) {
			modelSign += string(unicode.ToLower(ch))
		}
	}
	tableName := strings.ToLower(flect.Underscore(flect.Pluralize(fileName)))
	data := Data{FileName: fileName, ModelName: fileName, ModelSign: modelSign, TableName: tableName}
	tmpl, err := template.New("model").Parse(fileString)
	if err != nil {
		return
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
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
