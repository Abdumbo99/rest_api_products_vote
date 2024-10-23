package vote

import (
	"api_assignment/api/models/product"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AllVotes fetched all votes from the db
func (vModel VoteModel) AllVotes() ([]*Vote, error) {
	coll := vModel.DB.Database("trial").Collection("votes")

	cur, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	allVotes := make([]*Vote, 0)
	if err := cur.All(context.TODO(), &allVotes); err != nil {
		return nil, err
	}
	return allVotes, nil

}

// PostVote handles the repo side of the posting/updating of a vote
func (vModel VoteModel) PostVote(newVote *Vote) (*bool, error) {

	coll := vModel.DB.Database("trial").Collection("votes")

	// create the search filter
	filter := bson.D{{Key: "product_id", Value: newVote.ProductID}, {Key: "session_id", Value: newVote.SessionID}}
	// update fields
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "product_id", Value: newVote.ProductID},
		{Key: "session_id", Value: newVote.SessionID}, {Key: "rate", Value: newVote.Rate}}}}
	// upsert; insert or update if exists
	opts := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return nil, err
	}

	// nothing matched; completely new
	if result.MatchedCount == 0 {
		alreadyExist := false
		return &alreadyExist, nil
	} else {
		alreadyExist := true
		return &alreadyExist, nil
	}

}

// GetVotesBySessionID handles the db side of returning all votes with the specified session id
func (vModel VoteModel) GetVotesBySessionID(sessionID string) ([]*Vote, error) {

	coll := vModel.DB.Database("trial").Collection("votes")

	filter := bson.D{{Key: "session_id", Value: sessionID}}

	var foundVotes []*Vote
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cur.All(context.TODO(), &foundVotes)

	return foundVotes, nil

}

// GetVotesByProductID fetches all votes by the corresponding product id
func (vModel VoteModel) GetVotesByProductID(productID string) ([]*Vote, error) {

	coll := vModel.DB.Database("trial").Collection("votes")

	filter := bson.D{{Key: "product_id", Value: productID}}

	var foundVotes []*Vote
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	cur.All(context.TODO(), &foundVotes)

	return foundVotes, nil

}

// GetAvergageVotesForAllProducts handles the actual logic of fetching votes and calculating avgs.
func (vModel VoteModel) GetAvergageVotesForAllProducts(products map[string]*product.Product) (map[string]*VoteResult, error) {

	coll := vModel.DB.Database("trial").Collection("votes")

	// fetch everything
	filter := bson.D{{}}

	var foundVotes []*Vote
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cur.All(context.TODO(), &foundVotes)

	// fill data from the quesry result
	avgVotes := make(map[string]*VoteResult)
	for _, vote := range foundVotes {
		if _, ok := avgVotes[vote.ProductID]; !ok {
			avgVotes[vote.ProductID] = &VoteResult{}
		}
		avgVotes[vote.ProductID].sum += vote.Rate
		avgVotes[vote.ProductID].VotesCount++
	}

	// calculate the avg
	for prodID := range avgVotes {
		avgVotes[prodID].Avg = float64(avgVotes[prodID].sum) / float64(avgVotes[prodID].VotesCount)
	}

	// fill products with no votes
	for prodID := range products {
		if _, ok := avgVotes[prodID]; !ok {
			avgVotes[prodID] = &VoteResult{sum: 0, Avg: 0.0, VotesCount: 0}

		}
	}

	return avgVotes, nil
}
