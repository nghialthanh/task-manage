package controllers

import (
	"context"
	"fmt"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"task-manage/database"
	helper "task-manage/helpers"
	"task-manage/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var projectCollection *mongo.Collection = database.OpenCollection(database.Client, "project")
var roleProjectCollection *mongo.Collection = database.OpenCollection(database.Client, "roleProject")
var LogoDefault = "https://talent2win.com/wp-content/uploads/2021/09/LogoDefault.png"

func CreateProject() gin.HandlerFunc {
	return func(c *gin.Context) {

		var project models.Project
		var user models.User
		userId, exists := c.Get("user_id")

		if !exists {
			helper.SendResponse(c, helper.Response{Status: http.StatusBadRequest, Message: []string{"Failed to call API get info user"}})
			return
		}

		if err := c.BindJSON(&project); err != nil {
			helper.SendResponse(c, helper.Response{Status: http.StatusBadRequest, Error: []string{err.Error()}})
			return
		}

		validationErr := validate.Struct(project)
		if validationErr != nil {
			helper.SendResponse(c, helper.Response{Status: http.StatusBadRequest, Error: []string{validationErr.Error()}})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		project.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		project.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		project.ID = primitive.NewObjectID()
		project.Logo = &LogoDefault

		resultInsertionNumber, insertErr := projectCollection.InsertOne(ctx, project)
		defer cancel()

		if insertErr != nil {
			msg := fmt.Sprintf("Project item was not created")
			helper.SendResponse(c, helper.Response{Status: http.StatusInternalServerError, Error: []string{msg}})
			return
		}

		//find user by token
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()

		if err != nil {
			helper.SendResponse(c, helper.Response{Status: http.StatusInternalServerError, Error: []string{err.Error()}})
			return
		}
		// create role project
		var roleProject models.RoleProject

		roleProject.ID = project.ID
		roleProject.ID_User = user.ID
		roleProject.Role = "OWNER"
		_, insertErr = roleProjectCollection.InsertOne(ctx, roleProject)
		defer cancel()

		if insertErr != nil {
			msg := fmt.Sprintf("Project item was not created")
			helper.SendResponse(c, helper.Response{Status: http.StatusInternalServerError, Error: []string{msg}})
			return
		}

		helper.SendResponse(c, helper.Response{Status: http.StatusOK, Data: resultInsertionNumber})

	}
}
