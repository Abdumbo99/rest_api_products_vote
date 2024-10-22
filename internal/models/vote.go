package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VoteModel struct {
	DB *mongo.Client
}

// Vote holds the details of purchases of tickets
// TODO
type Vote struct {
	Rate      int    `json:"rate" bson:"rate"`
	SessionID string `json:"session_id" bson:"session_id"`
	Product
}

type PostResponse struct {
	AlreadyExist bool
	Vote
}

type VoteResult struct {
	sum        int
	Avg        float64 `json:"avg"`
	VotesCount int     `json:"votes_count"`
}

func (vModel *VoteModel) All() ([]*Vote, error) {

	coll := vModel.DB.Database("foodji").Collection("Products")

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

func (vModel *VoteModel) PostVote(newVote *Vote) (*PostResponse, error) {

	coll := vModel.DB.Database("foodji").Collection("Products")

	voteBson, err := bson.Marshal(newVote)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "product_id", Value: newVote.Product.ID}, {Key: "session_id", Value: newVote.SessionID}}
	update := bson.D{{Key: "$set", Value: voteBson}}
	opts := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount > 0 {
		return &PostResponse{AlreadyExist: true, Vote: *newVote}, nil
	} else {
		return &PostResponse{AlreadyExist: false, Vote: *newVote}, nil
	}

}

func (vModel *VoteModel) GetVotesBySessionID(sessionID string) ([]*Vote, error) {

	coll := vModel.DB.Database("foodji").Collection("Products")

	filter := bson.D{{Key: "session_id", Value: sessionID}}

	var foundVotes []*Vote
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cur.All(context.TODO(), &foundVotes)

	return foundVotes, nil

}

func (vModel *VoteModel) GetVotesByProductID(productID string) ([]*Vote, error) {

	coll := vModel.DB.Database("foodji").Collection("Products")

	filter := bson.D{{Key: "product_id", Value: productID}}

	var foundVotes []*Vote
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cur.All(context.TODO(), &foundVotes)

	return foundVotes, nil

}

func (vModel *VoteModel) GetAvergageVotesByProduct(sessionID string) (map[string]*VoteResult, error) {

	coll := vModel.DB.Database("foodji").Collection("Products")

	filter := bson.D{{}}

	var foundVotes []*Vote
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cur.All(context.TODO(), &foundVotes)

	avgVotes := make(map[string]*VoteResult)
	for _, vote := range foundVotes {
		if _, ok := avgVotes[vote.Product.ID]; !ok {
			avgVotes[vote.Product.ID] = &VoteResult{}
		}
		avgVotes[vote.Product.ID].sum += vote.Rate
		avgVotes[vote.Product.ID].VotesCount++
	}

	for prodID := range avgVotes {
		avgVotes[prodID].Avg = float64(avgVotes[prodID].sum) / float64(avgVotes[prodID].VotesCount)
	}
	return avgVotes, nil
}
