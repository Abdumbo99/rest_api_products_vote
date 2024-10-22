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
)

// Gin equivalent of the timeHandler
func fetchProductsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, utils.FetchProducts())
	}
}

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

func main() {
	//TODO move password to secrets
	mongoPass := "gb9MPHOre4hGm5ph"
	connectionString := fmt.Sprintf("mongodb+srv://abdul:%s@cluster0.ewrnc.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", mongoPass)
	client := createMongoClient(connectionString)
	collection := client.Database("foodji").Collection("Products")
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

	/*
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
	*/

	router := gin.Default()

	//Swagger endpoint
	router.GET("/products", fetchProductsHandler())

	router.Run(":8080")
}
