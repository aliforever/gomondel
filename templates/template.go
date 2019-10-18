package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

func (t *Template) getSignFromName(name string) (sign string) {
	r := []rune(name)
	for i := 0; i < len(r); i++ {
		ch := r[i]
		if unicode.IsUpper(ch) {
			sign += string(unicode.ToLower(ch))
		}
	}
	return
}

func (t *Template) CreateModel(modelName string, parentName, parentKeyType *string) (path string, err error) {
	path, err = t.CurrentPath()
	if err != nil {
		return
	}
	fileName := flect.Singularize(modelName)
	fileName = flect.Capitalize(fileName)
	modelSign := t.getSignFromName(fileName)
	parentKey := ""
	if parentKeyType != nil {
		parentKey = *parentKeyType
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
			parentModelSign := t.getSignFromName(parentModelName)
			parentMethod := t.parentMethod()
			parentMethodStr, err = TemplateData{}.FillModelParentMethod(parentMethod, modelSign, modelName, parentModelName, parentModelSign)
			if err != nil {
				return
			}
			field := t.parentField()
			parentField, err = TemplateData{}.FillModelParentField(field, parentModelName, parentKey)
			if err != nil {
				return
			}
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

	tableName := strings.ToLower(flect.Underscore(flect.Pluralize(fileName)))
	fileString, err = TemplateData{}.FillModel(fileString, fileName, fileName, modelSign, tableName, parentMethodStr, parentField)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, []byte(fileString), os.ModePerm)
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
	cmd := exec.Command("go", "fmt", path+"/", "./...")
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
