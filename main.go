package main

import (
	"context"
	"fmt"
	"strings"

	// "time"

	"github.com/gin-gonic/gin"

	"goproj/db"
	"goproj/handlers"
	"goproj/repositories"
	"goproj/services"
)

func ValidateAuth(userRepository repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Authorization")
		if authToken == "" {
			c.AbortWithStatus(403)
			return
		}

		authToken = strings.ReplaceAll(authToken, "Bearer ", "")

		found, session, err := userRepository.GetSessionById(authToken)
		if err != nil {
			fmt.Printf("err :: %+v\n", err)
			c.AbortWithStatus(403)
			return
		}
		if !found {
			fmt.Println("Not found")
			c.AbortWithStatus(403)
			return
		}
		c.Set("session", session)

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "example.com")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	fmt.Println("Starting...")

	dbName := "dataset"
	client, err := db.CreateDatabaseConnection(dbName)
	if err != nil {
		fmt.Println("Failed to connect to DB")
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database(dbName)

	// Repositories
	userRepository := repositories.NewInstanceOfUserRepository(db)

	// Services
	userService := services.NewInstanceOfUserService(userRepository)

	// Handlers
	userHandler := handlers.NewInstanceOfUserHandler(userService)

	router := gin.Default()
	router.Use(CORSMiddleware())

	userAPI := router.Group("/user")
	{
		userAPI.POST("/signin", userHandler.SignIn)
		userAPI.POST("/signup", userHandler.SignUp)
		userAPI.GET("/", ValidateAuth(userRepository), userHandler.GetAllUser)
		userAPI.GET("/:id", ValidateAuth(userRepository), userRepository.GetUserByID)
		userAPI.DELETE("/:id", ValidateAuth(userRepository), userHandler.Delete)
		userAPI.UPDATE("/:id", ValidateAuth(userRepository), userHandler.Update)
	}

	router.Run(":8080")
}
