package templates

import (
	"fmt"
)

func (t *Template) init(dbName string) string {
	return fmt.Sprintf(`package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	DB *mongo.Database
	DBClient *mongo.Client
)

func NewContext() context.Context {
	c, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return c
}

func InitMongoDB() (err error) {
	DBClient, err = mongo.Connect(NewContext(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err == nil {
		DB = DBClient.Database("%s")
	}
	return
}`, dbName)
}
