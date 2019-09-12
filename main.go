package main

import (
	"flag"
	"fmt"

	"github.com/aliforever/gomondel/templates"
)

func main() {
	var init string
	var model string
	flag.StringVar(&init, "init", "", "--init=database_name")
	flag.StringVar(&model, "model", "", "--model=ModelName")
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
		path, err := CreateModel(model)
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

func CreateModel(modelName string) (path string, err error) {
	t := templates.Template{}
	path, err = t.CreateModel(modelName)
	return
}
