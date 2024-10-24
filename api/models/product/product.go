package product

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Product is simple struct that represents a product with its field, like ids, and name
type Product struct {
	GloballID string `json:"global_id" bson:"global_id"`
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
}

// FetchProducts calls the endpoint and return the products from there
func FetchProducts() map[string]*Product {

	fmt.Println("Fetching products...")
	const endpoint = "https://amperoid.tenants.foodji.io/machines/4bf115ee-303a-4089-a3ea-f6e7aae0ab94"

	resp, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// holder of products
	products := make(map[string]*Product)

	rawDataMap := make(map[string]any)
	json.Unmarshal(body, &rawDataMap)

	// ugly parsing to the products
	productsRawData := rawDataMap["data"].(map[string]any)["machineProducts"].([]any)

	// save products
	for i, el := range productsRawData {
		element := el.(map[string]any)
		globalID := element["id"].(string)
		name := element["name"].(string)
		localID := strconv.Itoa(i)
		products[localID] = &Product{GloballID: globalID, ID: localID, Name: name}
	}

	fmt.Println("Fetched the following products:")
	for _, pr := range products {
		fmt.Printf("%+v\n", pr)
	}

	return products
}
