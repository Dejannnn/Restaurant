package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Dejannnn/Restaurant.git/database"
	"github.com/Dejannnn/Restaurant.git/helpers"
	"github.com/Dejannnn/Restaurant.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// defer cancel()
		// recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		// if err != nil || recordPerPage < 1 {
		// 	recordPerPage = 10
		// }
		// currentPage, err := strconv.Atoi(c.Query("currentPage"))
		// if err != nil || currentPage < 1 {
		// 	currentPage = 1
		// }

		// startIndex := (currentPage - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// matchStage := bson.E{"$match", bson.D{{}}}
		// projectStage := bson.D{
		// 	{"$project", bson.D{
		// 		{"_id", 0},
		// 		{"total_count", 1},
		// 		{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		// 	}}}

		// result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage})

		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// }

		// var users []bson.M
		// if err := result.All(ctx, &users); err != nil {
		// 	log.Fatal(err)
		// }

		// c.JSON(http.StatusOK, users[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userId := c.Param("userId")

		var user models.User
		filter := bson.M{"_id": userId}
		err := userCollection.FindOne(ctx, filter).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fetch user by id din't pass"})
		}

		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		countRes, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if countRes > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exist"})
			return
		}
		user.ID = primitive.NewObjectID()
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.User_id = user.ID.Hex()

		//HASH PASSWORD
		hashPassord, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Hash password faild"})
			return
		}
		fmt.Println("Hasshed password", hashPassord)

		user.Password = hashPassord
		//CREATE TOKEN
		// helpers.GenerateAllTokens(*user.)

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"email": user.Email}
		err := userCollection.FindOne(ctx, filter).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong credentials"})
			return
		}
		isValidPassword := VerifyPassword(foundUser.Password, user.Password)
		if !isValidPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong credentials"})
			return
		}

		//Generate tokens
		token, refresh_token, err := helpers.GenerateAllTokens(foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Generate new token faild"})
		}

		helpers.UpdateAllTokens(token, refresh_token, foundUser.User_id)
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
