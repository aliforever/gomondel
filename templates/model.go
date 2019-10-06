package templates

func (t *Template) createModel(modelName string) string {
	return `package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type {{.ModelName}} struct {
	Id primitive.ObjectID ` + "`" + `bson:"_id,omitempty"` + "`" + `
}

func ({{.ModelSign}} {{.ModelName}}) New() (n{{.ModelSign}} *{{.ModelName}}, err error) {
	n{{.ModelSign}} = &{{.ModelName}}{Id: primitive.NewObjectID()}
	_, err = DB.Collection("{{.TableName}}").InsertOne(NewContext(), &n{{.ModelSign}})
	return
}

func ({{.ModelSign}} *{{.ModelName}}) Create() (err error) {
	_, err = DB.Collection("{{.TableName}}").InsertOne(NewContext(), &{{.ModelSign}})
	return
}

func ({{.ModelSign}} {{.ModelName}}) FindById(id primitive.ObjectID) (f{{.ModelSign}} *{{.ModelName}}, err error) {
	err = DB.Collection("{{.TableName}}").FindOne(NewContext(), bson.M{"_id": id}).Decode(&f{{.ModelSign}})
	return
}

func ({{.ModelSign}} *{{.ModelName}}) FindOne() (err error) {
	err = DB.Collection("{{.TableName}}").FindOne(NewContext(), *{{.ModelSign}}).Decode(&{{.ModelSign}})
	return
}

func ({{.ModelSign}} {{.ModelName}}) Find() ({{.ModelSign}}s []{{.ModelName}}, err error) {
	ctx := NewContext()
	var cur *mongo.Cursor
	cur, err = DB.Collection("{{.TableName}}").Find(ctx, bson.D{})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result {{.ModelName}}
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		{{.ModelSign}}s = append({{.ModelSign}}s, result)
	}
	if err = cur.Err(); err != nil {
		return
	}
	return
}

func ({{.ModelSign}} *{{.ModelName}}) Save() (err error) {
	_, err = DB.Collection("{{.TableName}}").UpdateOne(NewContext(), bson.M{"_id": {{.ModelSign}}.Id}, bson.M{"$set": &{{.ModelSign}}})
	return
}

func ({{.ModelSign}} *{{.ModelName}}) Remove() (err error) {
	_, err = DB.Collection("{{.TableName}}").DeleteOne(NewContext(), bson.M{"_id": {{.ModelSign}}.Id})
	return
}
`
}
