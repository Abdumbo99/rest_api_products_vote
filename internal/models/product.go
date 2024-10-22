package models

// Ticket holds the details of a ticket
// TODO

type Product struct {
	GloballID string `json:"global_id" bson:"global_id"`
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
}
