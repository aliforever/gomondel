package funcs

import (
	"fmt"

	"github.com/aliforever/gomondel/templates"
	"github.com/gobuffalo/flect"
)

func InitDatabase(dbPath, dbName string) (path string, err error) {
	path, err = templates.Template{}.Init(dbPath, dbName)
	if err != nil {
		return
	}
	err = templates.Template{}.GoFmtPath(dbPath + "/")
	return
}

func CreateModel(projectPath, modelName string, modelIdType, parentName, parentIdType *string, fields []templates.ModelField) (path string, err error) {
	path, err = templates.Template{}.CreateModel(projectPath, modelName, modelIdType, parentName, parentIdType, fields)
	if err != nil {
		return
	}
	err = templates.Template{}.GoFmtPath(projectPath + "/")
	return
}

func MakeModelFieldsFromMap(m map[string]string) (result []templates.ModelField) {
	for k, v := range m {
		tag := fmt.Sprintf("`"+`bson:"%s,omitempty"`+"`", flect.Underscore(k))
		result = append(result, templates.ModelField{Name: k, Type: v, Tag: tag})
	}
	return
}
