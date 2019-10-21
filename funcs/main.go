package funcs

import "github.com/aliforever/gomondel/templates"

func InitDatabase(dbPath, dbName string) (path string, err error) {
	path, err = templates.Template{}.Init(dbPath, dbName)
	return
}

func CreateModel(projectPath, modelName string, modelIdType, parentName, parentIdType *string) (path string, err error) {
	path, err = templates.Template{}.CreateModel(projectPath, modelName, modelIdType, parentName, parentIdType)
	return
}
