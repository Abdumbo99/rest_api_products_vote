package models

// Ticket holds the details of a ticket
// TODO

type Product struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
