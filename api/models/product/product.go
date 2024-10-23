package product

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Ticket holds the details of a ticket
// TODO

type Product struct {
	GloballID string `json:"global_id" bson:"global_id"`
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
}

func FetchProducts() map[string]*Product {

	fmt.Println("Fetching products...")
	// TODO add retries
	const link = "https://amperoid.tenants.foodji.io/machines/4bf115ee-303a-4089-a3ea-f6e7aae0ab94"

	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// create the slice to hold products
	products := make(map[string]*Product)

	rawDataMap := make(map[string]any)
	json.Unmarshal(body, &rawDataMap)

	// ugly parsing to the products
	productsRawData := rawDataMap["data"].(map[string]any)["machineProducts"].([]any)

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
