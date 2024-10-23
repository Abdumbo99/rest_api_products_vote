package vote

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type VoteModel struct {
	DB *mongo.Client
}

// Vote holds the data of any vote in the system
type Vote struct {
	Rate      int    `json:"rate" bson:"rate"`
	SessionID string `json:"session_id" bson:"session_id"`
	ProductID string `json:"product_id" bson:"product_id"`
}

// VoteResult is a simple container used to hold the avg of the votes of a specific product
// along with some addiational data
type VoteResult struct {
	sum        int
	Avg        float64 `json:"avg"`
	VotesCount int     `json:"votes_count"`
}
