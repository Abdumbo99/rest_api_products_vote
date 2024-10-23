package main

import (
	"api_assignment/api/handler"
	"api_assignment/api/middleware"
	"context"
	"fmt"
	"net/http"

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

func CloseMongoConnection(client *mongo.Client) {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
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
	//TODO move password to secrets
	mongoPass := "gb9MPHOre4hGm5ph"
	connectionString := fmt.Sprintf("mongodb+srv://abdul:%s@cluster0.ewrnc.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", mongoPass)
	client := createMongoClient(connectionString)

	app := handler.NewApp(client)

	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	// setup the cookie and use it
	store := cookie.NewStore([]byte("sessioon-secret-key"))
	router.Use(sessions.Sessions("session_cookie", store))

	router.Use(middleware.CheckSession())
	router.Use(middleware.Log())

	// endpoints
	router.GET("/products", app.AllProductsHandler())
	router.GET("/votes", app.AllVotessHandler())
	router.POST("/votes", app.PostVoteHandler())
	router.GET("/votes/product/:id", app.GetVotesByProductIDHandler())
	router.GET("/votes/session/:id", app.GetVotesBySessionIDHandler())
	router.GET("/products/avgs", app.GetAverageVotesForAllProductsHandler())

	router.GET("/", hello())
	router.Run(":8080")
}
