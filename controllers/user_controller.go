package controllers

import (
	"net/http"
	"time"
	"urlier/configs"
	"urlier/models"
	"urlier/responses"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignUp(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	// check this is a new user
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)
	if err == nil {
		return c.JSON(http.StatusConflict, responses.Response{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "User with email already exists"}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	// hash password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	// create a new user with the hashedpassword
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}
	// save user to db
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}

func SignIn(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var loginInfo models.LoginInfo
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.Bind(&loginInfo); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err := userCollection.FindOne(ctx, bson.M{"email": loginInfo.Email}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	// check correct password
	if !CheckPasswordHash(loginInfo.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, responses.Response{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "Invalid password"}})
	}

	cleanUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: "",
	}
	return c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": cleanUser}})
}
