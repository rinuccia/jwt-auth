package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rinuccia/jwt-auth/database"
	"github.com/rinuccia/jwt-auth/models"
	"github.com/rinuccia/jwt-auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var users *mongo.Collection = database.OpenCollection(database.Client, "users")

var validate = validator.New()

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(userPassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		return false
	}
	return true
}

func SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	user := models.User{}
	err := c.BindJSON(&user)
	validationErr := validate.Struct(user)
	if err != nil || validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input body"})
		return
	}

	count, err := users.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
		return
	}

	password, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	user.Password = password

	user.CreatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}

	user.UpdatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}

	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	token, refreshToken, err := utils.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	user.Token = token
	user.RefreshToken = refreshToken
	resultInsertionNumber, insertErr := users.InsertOne(ctx, user)
	defer cancel()
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, resultInsertionNumber)
}

func SignIn(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input body"})
		return
	}

	err := users.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
		return
	}

	if ok := VerifyPassword(user.Password, foundUser.Password); !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
		return
	}

	token, refreshToken, err := utils.GenerateTokens(foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	err = utils.UpdateTokens(token, refreshToken, foundUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	err = users.FindOne(ctx, bson.M{"user_id": foundUser.UserID}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": foundUser.Token, "refresh_token": foundUser.RefreshToken})
}
