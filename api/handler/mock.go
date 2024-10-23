package handler

import (
	"api_assignment/api/models/product"
	"api_assignment/api/models/vote"
)

// MockVoteService is a mock implementation of the voteService interface
type MockVoteService struct {
	mockAllVotes          []*vote.Vote
	mockGetVotesBySession []*vote.Vote
	mockGetVotesByProduct []*vote.Vote
	mockPostVoteExists    *bool
	mockAvgVotes          map[string]*vote.VoteResult
	mockError             error
}

func (m *MockVoteService) AllVotes() ([]*vote.Vote, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockAllVotes, nil
}

func (m *MockVoteService) PostVote(newVote *vote.Vote) (*bool, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockPostVoteExists, nil
}

func (m *MockVoteService) GetVotesBySessionID(sessionID string) ([]*vote.Vote, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockGetVotesBySession, nil
}

func (m *MockVoteService) GetVotesByProductID(productID string) ([]*vote.Vote, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockGetVotesByProduct, nil
}

func (m *MockVoteService) GetAverageVotesForAllProducts(products map[string]*product.Product) (map[string]*vote.VoteResult, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockAvgVotes, nil
}
