# gomondel
Golang MongoDB Model Generator with CRUD

[![gomondel](http://i3.ytimg.com/vi/I1ncw9C696E/maxresdefault.jpg)](https://www.youtube.com/watch?v=I1ncw9C696E&)

## Description

Create ready-to-go models when working with the [Official Go Driver for MongoDB](https://github.com/mongodb/mongo-go-driver) including CRUD Operations (Create, Read, Update, Delete).



## Installation 
```go get -u github.com/aliforever/gomondel```



## Usage
###### Note: Run `go install` after `go get` to install gomondel in bin folder located at %GOPATH%/bin (make sure bin folder is added to path variable).

### Initialize New Database for Project
To add database connection to your project: 

1) cd into your project directory, forexample: `cd %GOPATH%/src/myproject/`, and then run:

    `gomondel --init=your_database_name_here`

    (This will create a models folder with a db.go file)

2) call `models.InitMongoDB()` in your main file to initialize the database.

### Create New Model 
To create a new model based on your MongoDB collections, run:

`gomondel --model=ModelName`
    
(This will create modelname.go file in models folder)


## Note:
gomondel uses [Inflect Package](github.com/gobuffalo/flect) to pluralize model names, so forexample if your collection name is `users`, you should use `User` as model name.
This is to comply to golang naming conventions, your model struct type will be named User.

## Contribution
You're free to create issues and pull requests to help complete gomondel!

