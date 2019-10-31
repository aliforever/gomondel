package templates

import (
	"bytes"
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

func (t Template) Init(projectPath, dbName string) (path string, err error) {
	if projectPath != "" {
		path = projectPath
		_, err = os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				err = errors.New(fmt.Sprintf("Path %s does not exists", projectPath))
			}
			return
		}
	} else {
		path, err = t.CurrentPath()
		if err != nil {
			return
		}
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

func (t Template) getSignFromName(name string) (sign string) {
	r := []rune(name)
	for i := 0; i < len(r); i++ {
		ch := r[i]
		if unicode.IsUpper(ch) {
			sign += string(unicode.ToLower(ch))
		}
	}
	return
}

func (t Template) CreateModel(projectPath, modelName string, modelIdType, parentName, parentIdType *string, fields []ModelField) (path string, err error) {
	if projectPath != "" {
		path = projectPath
		_, err = os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				err = errors.New(fmt.Sprintf("Path %s does not exists", projectPath))
			}
			return
		}
	} else {
		path, err = t.CurrentPath()
		if err != nil {
			return
		}
	}
	fileName := flect.Singularize(modelName)
	fileName = flect.Capitalize(fileName)
	modelSign := t.getSignFromName(fileName)
	parentKey := ""
	if parentIdType != nil {
		parentKey = *parentIdType
	}
	modelKey := ""
	if modelIdType != nil {
		modelKey = *modelIdType
	}
	parentMethodStr := ""
	parentChildMethods := ""
	parentField := ""
	parent := ""
	parentModelName := ""
	if parentName != nil {
		parent = fmt.Sprintf(ModelsPath, path) + "/" + *parentName + ".go"
		if _, err = os.Stat(parent); err != nil {
			err = errors.New(fmt.Sprintf("%s %s", parent, err.Error()))
			return
		} else {
			parentModelName = flect.Singularize(*parentName)
			parentModelName = flect.Capitalize(parentModelName)
			parentModelSign := t.getSignFromName(parentModelName)
			parentMethod := t.modelParentMethod()
			parentMethodStr, err = TemplateData{}.FillModelParentMethod(parentMethod, modelSign, modelName, parentModelName, parentModelSign)
			if err != nil {
				return
			}
			field := t.modelParentField()
			parentField, err = TemplateData{}.FillModelParentField(field, parentModelName, parentKey)
			if err != nil {
				return
			}
			parentChildMethods = t.parentChildMethods()
			parentChildMethods, err = TemplateData{}.FillParentChildMethods(parentChildMethods, parent, modelName, modelSign, modelKey, parentModelSign, parentModelName)
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
	fileString, err = TemplateData{}.FillModel(fileString, fileName, fileName, modelSign, "", modelKey, tableName, parentMethodStr, parentField, fields)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, []byte(fileString), os.ModePerm)
	if err != nil {
		return
	}
	if parent != "" && parentModelName != "" && parentChildMethods != "" {
		var parentFile []byte
		parentFile, err = ioutil.ReadFile(parent)
		if err != nil {
			return
		}
		parentFile = append(parentFile, []byte("\n\n"+parentChildMethods)...)
		err = ioutil.WriteFile(parent, parentFile, os.ModePerm)
		if err != nil {
			return
		}
	}
	err = t.GoFmtCurrentPath()
	return
}

func (t Template) GoFmtCurrentPath() (err error) {

	var path string
	path, err = t.CurrentPath()
	if err != nil {
		return
	}
	cmd := exec.Command("go", "fmt", path+"/", "./...")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	return
}

func (t Template) GoFmtPath(path string) (err error) {
	// cmd := exec.Command("go", "fmt", path+"/", "./...")
	fmt.Println("fmt path " + path)
	cmd := exec.Command("gofmt", "-w", path)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	return
}

func (t Template) CurrentPath() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}
	return
}
