package vote

import (
	"api_assignment/api/models/product"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostResponse struct {
	AlreadyExist bool
	Vote
}

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

func (vModel VoteModel) PostVote(newVote *Vote) (*PostResponse, error) {

	coll := vModel.DB.Database("trial").Collection("votes")

	filter := bson.D{{Key: "product_id", Value: newVote.ProductID}, {Key: "session_id", Value: newVote.SessionID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "product_id", Value: newVote.ProductID},
		{Key: "session_id", Value: newVote.SessionID}, {Key: "rate", Value: newVote.Rate}}}}
	opts := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return &PostResponse{AlreadyExist: false, Vote: *newVote}, nil
	} else {
		return &PostResponse{AlreadyExist: true, Vote: *newVote}, nil
	}

}

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

func (vModel VoteModel) GetAvergageVotesForAllProducts(products map[string]*product.Product) (map[string]*VoteResult, error) {

	coll := vModel.DB.Database("trial").Collection("votes")

	filter := bson.D{{}}

	var foundVotes []*Vote
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cur.All(context.TODO(), &foundVotes)

	avgVotes := make(map[string]*VoteResult)
	for _, vote := range foundVotes {
		if _, ok := avgVotes[vote.ProductID]; !ok {
			avgVotes[vote.ProductID] = &VoteResult{}
		}
		avgVotes[vote.ProductID].sum += vote.Rate
		avgVotes[vote.ProductID].VotesCount++
	}

	for prodID := range avgVotes {
		avgVotes[prodID].Avg = float64(avgVotes[prodID].sum) / float64(avgVotes[prodID].VotesCount)
	}

	for prodID := range products {
		if _, ok := avgVotes[prodID]; !ok {
			avgVotes[prodID] = &VoteResult{sum: 0, Avg: 0.0, VotesCount: 0}

		}
	}

	return avgVotes, nil
}
