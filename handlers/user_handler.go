package handlers

import (
	// "fmt"
	"context"
	"encoding/json"
	"goproj/db"
	"goproj/models"
	"goproj/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type UserHandler struct {
	UserService services.UserService
}

/////////////////////////////////NewInstanceOfUserHandler/////////////////////////////////////////////////////

func NewInstanceOfUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

////////////////////////////SignIn//////////////////////////////////////////////////////////

func (u *UserHandler) SignIn(c *gin.Context) {
	var body models.SignInBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	token, err := u.UserService.SignIn(body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Signed in", "token": token})
	return
}

///////////////////////////////SignUp///////////////////////////////////////////////////////

func (u *UserHandler) SignUp(c *gin.Context) {
	var body models.SignUpBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	token, err := u.UserService.SignUp(body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Signed up", "token": token})
	return
}

// func (u *UserHandler) DeletePerson(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Set("content-type", "application/json")
// 	params := mux.Vars(request)
// 	id, _ := primitive.ObjectIDFromHex(params["id"])
// 	// var person Person
// 	dbName := "demoDB"
// 	collectionn := "people" // move to env
// 	client, _ := db.CreateDatabaseConnection(dbName)
// 	collection := client.Database(dbName).Collection(collectionn)
// 	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	// err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
// 	_, err := collection.DeleteOne(ctx, models.User{ID: id})

// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
// 		return
// 	}
// 	// json.NewEncoder(response).Encode(person)

// }

//////////////////////////////////////////////////////////////////////////////////////

// func (u *UserHandler) GetUser(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Set("content-type", "application/json")
// 	params := mux.Vars(request)
// 	id, _ := primitive.ObjectIDFromHex(params["id"])
// 	var person models.User
// 	dbName := "dataset"
// 	collectionn := "user" // move to env
// 	client, _ := db.CreateDatabaseConnection(dbName)
// 	collection := client.Database(dbName).Collection(collectionn)
// 	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	err := collection.FindOne(ctx, models.User{ID: id}).Decode(&person)
// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
// 		return
// 	}
// 	json.NewEncoder(response).Encode(person)
// }

////////////////////////////////////GetAllUser//////////////////////////////////////////////////

func (u *UserHandler) GetAllUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []models.User
	dbName := "dataset"
	collectionn := "user" // move to env
	client, _ := db.CreateDatabaseConnection(dbName)
	collection := client.Database(dbName).Collection(collectionn)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person models.User
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

/////////////////////////////////Update/////////////////////////////////////////////////////

func (u *UserHandler) Update(c *gin.Context) {
	carsID := c.Param("id")

	var body models.UpdateUser
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := body.Valid(); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	session, exists := u.GetSession(c)
	if !exists {
		c.JSON(403, gin.H{"message": "error: unauthorized"})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	err := u.UserService.Update(session, carsID, body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Updated car"})
	return
}

//////////////////////////////Delete////////////////////////////////////////////////////////

func (u *UserHandler) Delete(c *gin.Context) {
	session, exists := u.GetSession(c)
	if !exists {
		c.JSON(403, gin.H{"message": "error: unauthorized"})
		return
	}

	carsID := c.Param("id")

	err := u.UserService.Delete(session, carsID)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted car"})
	return
}

////////////////////////////////GetSession//////////////////////////////////////////////////////

func (u *UserHandler) GetSession(c *gin.Context) (models.Session, bool) {
	i, exists := c.Get("session")
	if !exists {
		return models.Session{}, false
	}
	session, ok := i.(models.Session)
	if !ok {
		return models.Session{}, false
	}
	return session, true
}
