package main

import (
	"api_assignment/api/handler"
	"api_assignment/api/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func createMongoClient(connectionString string) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	return client
}

// @Summary Root endpoint
// @Description Displays a simple hello message at the root.
// @Tags default
// @Produce plain
// @Success 200 {string} string "Hello message"
// @Router / [get]
func hello() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "hello"})
	}
}

func main() {
	//read db auth info
	err := godotenv.Load()
	mongoPass := os.Getenv("MONGO_PASS")
	mongoUser := os.Getenv("MONGO_USER")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoParams := os.Getenv("MONGO_PARAMS")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s",
		mongoUser, mongoPass, mongoHost, mongoParams)
	client := createMongoClient(connectionString)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := handler.NewApp(client)

	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	// setup the cookie and use it
	store := cookie.NewStore([]byte("sessioon-key"))
	router.Use(sessions.Sessions("session_cookie", store))

	router.Use(middleware.CheckSession())
	router.Use(middleware.Log())
	router.Use(middleware.CORSMiddleware())

	// endpoints
	router.GET("/products", app.AllProductsHandler())
	router.GET("/votes", app.AllVotessHandler())
	router.POST("/votes", app.PostVoteHandler())
	router.GET("/votes/product/:id", app.GetVotesByProductIDHandler())
	router.GET("/votes/session/:id", app.GetVotesBySessionIDHandler())
	router.GET("/products/avgs", app.GetAverageVotesForAllProductsHandler())

	router.GET("/", hello())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
