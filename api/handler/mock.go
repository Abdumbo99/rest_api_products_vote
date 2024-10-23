package handler

import (
	"api_assignment/api/models/product"
	"api_assignment/api/models/vote"
)

// MockVoteService is a mock implementation of the voteService interface
type MockVoteService struct {
	mockAllVotes          []*vote.VoteResult
	mockGetVotesBySession []*vote.VoteResult
	mockGetVotesByProduct []*vote.VoteResult
	mockPostVoteExists    *bool
	mockAvgVotes          map[string]*vote.ProductVote
	mockError             error
}

func (m *MockVoteService) AllVotes() ([]*vote.VoteResult, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockAllVotes, nil
}

func (m *MockVoteService) PostVote(newVote *vote.VoteResult) (*bool, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockPostVoteExists, nil
}

func (m *MockVoteService) GetVotesBySessionID(sessionID string) ([]*vote.VoteResult, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockGetVotesBySession, nil
}

func (m *MockVoteService) GetVotesByProductID(productID string) ([]*vote.VoteResult, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockGetVotesByProduct, nil
}

func (m *MockVoteService) GetAverageVotesForAllProducts(products map[string]*product.Product) (map[string]*vote.ProductVote, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockAvgVotes, nil
}
