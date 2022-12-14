package controllers

import (
	"context"
	"log"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"task-manage/database"

	helper "task-manage/helpers"
	"task-manage/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetListUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			helper.SendResponse(c, helper.Response{Status: http.StatusBadRequest, Error: []string{err.Error()}})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// recordPerPage := 10
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			helper.SendResponse(c, helper.Response{Status: http.StatusInternalServerError, Message: []string{"error occured while listing user items"}})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}

		helper.SendResponse(c, helper.Response{Status: http.StatusOK, Data: allusers[0]})
	}
}

// GetUser is the api used to tget a single user
func GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			helper.SendResponse(c, helper.Response{Status: http.StatusBadRequest, Error: []string{err.Error()}})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			helper.SendResponse(c, helper.Response{Status: http.StatusInternalServerError, Error: []string{err.Error()}})
			return
		}

		helper.SendResponse(c, helper.Response{Status: http.StatusOK, Data: user})

	}
}

// GetUserByToken take accesstoken and decode to take info of user
func GetUserByToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("user_id")

		if !exists {
			helper.SendResponse(c, helper.Response{Status: http.StatusBadRequest, Message: []string{"Failed to call API get info user"}})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()

		if err != nil {
			helper.SendResponse(c, helper.Response{Status: http.StatusInternalServerError, Error: []string{err.Error()}})
			return
		}

		helper.SendResponse(c, helper.Response{Status: http.StatusOK, Data: user})

	}
}
