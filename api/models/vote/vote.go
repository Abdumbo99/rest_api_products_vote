package vote

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type VoteModel struct {
	DB *mongo.Client
}

// Vote holds the details of purchases of tickets
// TODO
type Vote struct {
	Rate      int    `json:"rate" bson:"rate"`
	SessionID string `json:"session_id" bson:"session_id"`
	ProductID string `json:"product_id" bson:"product_id"`
}

type VoteResult struct {
	sum        int
	Avg        float64 `json:"avg"`
	VotesCount int     `json:"votes_count"`
}
