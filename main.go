package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aliforever/gomondel/utils"

	"github.com/aliforever/gomondel/funcs"
	"github.com/aliforever/gomondel/templates"
)

func main() {
	var init, model, modelParent, path, fields string
	flag.StringVar(&init, "init", "", "--init=database_name")
	flag.StringVar(&path, "path", "", "--path=/home/go/src/project_name")
	flag.StringVar(&model, "model", "", "--model=model_name[,int]")
	flag.StringVar(&fields, "fields", "", "--fields=username,string-password,string")
	flag.StringVar(&modelParent, "parent", "", "--parent=parent_model_name[,int]")
	flag.Parse()
	if init != "" && model != "" {
		fmt.Println("You can't run init and model at the same time")
		return
	}
	var err error
	if path == "" {
		path, err = utils.CurrentPath()
		if err != nil {
			fmt.Println("invalid path", path)
			return
		}
	} else {
		_, err = os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				err = errors.New(fmt.Sprintf("Path %s does not exists", path))
			}
			return
		}
	}
	if init != "" {
		dbPath, err := funcs.InitDatabase(path, init)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("Database %s init file added in %s, don't forget to call InitMongoDB() in your main function", init, dbPath))
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
		var modelFields []templates.ModelField
		if fields != "" {
			splitFields := strings.Split(fields, "-")
			m := map[string]string{}
			for _, field := range splitFields {
				splitFieldType := strings.Split(field, ",")
				if len(splitFieldType) != 2 {
					fmt.Println("invalid fields argument")
					return
				}
				fieldName := splitFieldType[0]
				fieldType := splitFieldType[1]
				m[fieldName] = fieldType
			}
			modelFields = funcs.MakeModelFieldsFromMap(m)
		}
		modelPath, err := funcs.CreateModel(path, model, keyType, parent, parentIdType, modelFields)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("Model %s file added in %s", model, modelPath))
		return
	}
}
