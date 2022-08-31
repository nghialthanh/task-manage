package controllers

import (
	"context"
	"fmt"
	"log"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	helper "task-manage/helpers"
	"task-manage/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshTokenEntity struct {
	refresh_token string
}

//CreateUser is the api used to tget a single user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := models.HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, _ := helper.GenerateAccessTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type)
		refreshToken, _ := helper.GenerateRefreshTokens( *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

//Login is the api used to tget a single user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User
		var foundUser models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "posted", "error": "login or passowrd is incorrect"})
			return
		}

		passwordIsValid, msg := models.VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		token, _ := helper.GenerateAccessTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type)
		refreshToken, _ := helper.GenerateRefreshTokens( foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)

	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody models.User
		var user models.User

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, err := helper.ValidateToken(*requestBody.Refresh_token)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		
		_err := userCollection.FindOne(ctx, bson.M{"user_id": claims.Uid}).Decode(&user)

		defer cancel()

		if _err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": _err.Error()})
			return
		}

		token, _ := helper.GenerateAccessTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type)

		helper.UpdateAllTokens(token, *requestBody.Refresh_token, user.User_id)
		

		_err = userCollection.FindOne(ctx, bson.M{"user_id": user.User_id}).Decode(&user)

		if _err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": _err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
