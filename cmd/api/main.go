package main

import (
	"api_assignment/internal/models"
	"api_assignment/pkg/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/google/uuid"
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
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}

func CloseMongoConnection(client *mongo.Client) {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

// Middleware to check and create a session if it doesn't exist
func checkSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// Check if the session ID exists
		sessionID := session.Get("session_id")
		if sessionID == nil {
			// If no session ID, create a new session
			newSessionID := uuid.New().String()

			// Set the session ID in the session data
			session.Set("session_id", newSessionID)

			// Save the session
			session.Save()

			fmt.Println("New session created with ID:", newSessionID)
		} else {
			fmt.Println("Existing session found with ID:", sessionID)
		}

		// Continue to the next middleware/handler
		c.Next()
	}
}

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
	prs := utils.FetchProducts()

	productsAsInterfaces := make([]interface{}, 0)

	for _, pr := range prs {
		productsAsInterfaces = append(productsAsInterfaces, pr)
	}

	collection := client.Database("trial").Collection("products")

	_, err := collection.InsertMany(context.TODO(), productsAsInterfaces)
	if err != nil {
		fmt.Println(err.Error())
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	name, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name of Index Created: " + name)

	products := make([]*models.Product, 0)

	filter := bson.D{}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("could not find elements!")
	}

	if err = cur.All(context.TODO(), &products); err != nil {
		panic(err)
	}

	for _, pr := range products {
		fmt.Println(pr)
	}

	app := &application{
		Products:    prs,
		voteService: models.VoteModel{DB: client},
	}

	router := gin.Default()

	store := cookie.NewStore([]byte("sessioon-secret-key"))
	router.Use(sessions.Sessions("session_cookie", store))

	router.Use(checkSession())

	//Swagger endpoint
	router.GET("/products", app.AllProductsHandler())
	router.GET("/votes", app.AllVotessHandler())
	router.POST("/votes", app.PostVoteHandler())
	router.GET("/votes/:id", app.GetVotesByProductIDHanlder())
	router.GET("/sessions/:id", app.GetVotesBySessionIDHanlder())
	router.GET("/products/avgs", app.GetAvergageVotesForAllProductsHanlder())

	router.GET("/", hello())
	router.Run(":8080")
}
