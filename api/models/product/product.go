package product

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Product is simple struct that represents a product with its field, like ids, and name
type Product struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

func AddProductsToDB(DB *mongo.Client) (map[string]*Product, error) {
	// Open the JSON file
	file, err := os.Open("products.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close() // Ensure the file is closed after reading

	var products []*Product

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&products)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	for _, pr := range products {
		fmt.Printf("%+v\n", pr)
	}

	coll := DB.Database("trial").Collection("products")

	productsInterface := make([]interface{}, len(products))
	for i, pr := range products {
		productsInterface[i] = pr
	}
	_, err = coll.InsertMany(context.TODO(), productsInterface)

	if err != nil {
		return nil, err
	}
	return nil, nil

}

// FetchProducts calls the endpoint and return the products from there
func FetchProducts(DB *mongo.Client) (map[string]*Product, error) {

	coll := DB.Database("trial").Collection("products")

	// fetch everything
	filter := bson.D{{}}

	var foundProducts []*Product
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cur.All(context.TODO(), &foundProducts)

	// holder of products
	products := make(map[string]*Product)

	// save products
	for _, product := range foundProducts {
		products[product.ID] = product
	}

	fmt.Println("Fetched the following products:")
	for id, pr := range products {
		fmt.Printf("%s : %+v\n", id, pr)
	}

	return products, nil
}
