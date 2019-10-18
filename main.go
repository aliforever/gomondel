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
	flag.StringVar(&model, "model", "", "--model=model_name")
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
		var parentKeyType *string
		if modelParent != "" {
			split := strings.Split(modelParent, ",")
			if len(split) == 2 {
				parentKeyType = &split[1]
			}
			parent = &split[0]
		}

		path, err := CreateModel(model, parent, parentKeyType)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("Model %s file added in %s", model, path))
		return
	}
	//fmt.Println(flect.Underscore("UserForm"))
	//models.InitMongoDB()
}

func InitDatabase(dbName string) (path string, err error) {
	t := templates.Template{}
	path, err = t.Init(dbName)
	return
}

func CreateModel(modelName string, parentName, parentKeyType *string) (path string, err error) {
	t := templates.Template{}
	path, err = t.CreateModel(modelName, parentName, parentKeyType)
	return
}
