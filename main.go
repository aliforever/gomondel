package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aliforever/gomondel/templates"
)

func main() {
	var init string
	var model string
	var modelParent string
	flag.StringVar(&init, "init", "", "--init=database_name")
	flag.StringVar(&model, "model", "", "--model=model_name[,int]")
	flag.StringVar(&modelParent, "parent", "", "--parent=parent_model_name[,int]")
	flag.Parse()
	if init != "" && model != "" {
		fmt.Println("You can't run init and model at the same time")
		return
	}
	if init != "" {
		path, err := InitDatabase(init)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("Database %s init file added in %s, don't forget to call InitMongoDB() in your main function", init, path))
		return
	}
	if model != "" {
		var parent *string
		var keyType *string
		var parentIdType *string
		split := strings.Split(model, ",")
		if len(split) == 2 {
			keyType = &split[1]
		}
		model = split[0]
		if modelParent != "" {
			split := strings.Split(modelParent, ",")
			if len(split) == 2 {
				parentIdType = &split[1]
			}
			parent = &split[0]
		}

		path, err := CreateModel(model, keyType, parent, parentIdType)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("Model %s file added in %s", model, path))
		return
	}
}

func InitDatabase(dbName string) (path string, err error) {
	t := templates.Template{}
	path, err = t.Init(dbName)
	return
}

func CreateModel(modelName string, modelIdType, parentName, parentIdType *string) (path string, err error) {
	t := templates.Template{}
	path, err = t.CreateModel(modelName, modelIdType, parentName, parentIdType)
	return
}
