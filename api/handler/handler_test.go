package handler

import (
	"api_assignment/api/middleware"
	"api_assignment/api/models/product"
	"api_assignment/api/models/vote"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter sets up test router
func setupRouter(app *Application) *gin.Engine {
	router := gin.Default()

	// Use the session middleware with a mock cookie store
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session_cookie", store))
	router.Use(middleware.CheckSession())

	router.GET("/products", app.AllProductsHandler())
	router.GET("/votes", app.AllVotessHandler())
	router.POST("/votes", app.PostVoteHandler())
	router.GET("/votes/session/:id", app.GetVotesBySessionIDHandler())
	router.GET("/votes/product/:id", app.GetVotesByProductIDHandler())
	router.GET("/products/avgs", app.GetAverageVotesForAllProductsHandler())
	return router
}

func TestAllProductsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	products := map[string]*product.Product{
		"p1": {ID: "p1", Name: "Product 1"},
	}
	app := &Application{
		Products:    products,
		voteService: &MockVoteService{},
	}

	router := setupRouter(app)

	// Test case: Products exist
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var respProducts map[string]*product.Product
	err := json.Unmarshal(w.Body.Bytes(), &respProducts)
	assert.NoError(t, err)
	assert.Equal(t, products, respProducts)

	// Test case: No products
	app.Products = map[string]*product.Product{}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Looks like there are no products so far.")
}

func TestAllVotessHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockVotes := []*vote.VoteResult{
		{ProductID: "p1", SessionID: "s1", Rate: 8},
	}
	app := &Application{
		voteService: &MockVoteService{
			mockAllVotes: mockVotes,
		},
	}

	router := setupRouter(app)

	// Test case: Votes exist
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/votes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var respVotes []*vote.VoteResult
	err := json.Unmarshal(w.Body.Bytes(), &respVotes)
	assert.NoError(t, err)
	assert.Equal(t, mockVotes, respVotes)

	// Test case: No votes
	app.voteService = &MockVoteService{}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/votes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Looks like there are no votes so far.")
}

func TestPostVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := &Application{
		Products: map[string]*product.Product{
			"p1": {ID: "p1", Name: "Product 1"},
		},
		voteService: &MockVoteService{
			mockPostVoteExists: func() *bool { v := false; return &v }(),
		},
	}

	// Set up router with session middleware
	router := setupRouter(app)

	// Test case: Successful vote post
	w := httptest.NewRecorder()
	voteBody := `{"product_id": "p1", "rate": 8}`
	req, _ := http.NewRequest(http.MethodPost, "/votes", strings.NewReader(voteBody))
	req.Header.Set("Content-Type", "application/json")

	// Now send the request with the session cookie
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Your vote has been received successfully!")

	// Test case: Invalid product
	voteBody = `{"product_id": "invalid", "rate": 8}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/votes", strings.NewReader(voteBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "No such product")

	// Test case: Invalid rate
	voteBody = `{"product_id": "p1", "rate": 15}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/votes", strings.NewReader(voteBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "rate must be between 0 and 10!")
}

func TestGetVotesBySessionIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockVotes := []*vote.VoteResult{
		{ProductID: "p1", SessionID: "s1", Rate: 8},
	}
	app := &Application{
		voteService: &MockVoteService{
			mockGetVotesBySession: mockVotes,
		},
	}

	router := setupRouter(app)

	// Test case: Votes found by session ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/votes/session/s1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var respVotes []*vote.VoteResult
	err := json.Unmarshal(w.Body.Bytes(), &respVotes)
	assert.NoError(t, err)
	assert.Equal(t, mockVotes, respVotes)

	// Test case: No votes for session
	app.voteService = &MockVoteService{}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/votes/session/s1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Looks like there are no votes for this session so far.")
}

func TestGetVotesByProductIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockVotes := []*vote.VoteResult{
		{ProductID: "p1", SessionID: "s1", Rate: 8},
	}
	app := &Application{
		Products: map[string]*product.Product{
			"p1": {ID: "p1", Name: "Product 1"},
		},
		voteService: &MockVoteService{
			mockGetVotesByProduct: mockVotes,
		},
	}

	router := setupRouter(app)

	// Test case: Votes found by product ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/votes/product/p1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var respVotes []*vote.VoteResult
	err := json.Unmarshal(w.Body.Bytes(), &respVotes)
	assert.NoError(t, err)
	assert.Equal(t, mockVotes, respVotes)

	// Test case: Invalid product ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/votes/product/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "No such product")
}

func TestGetAverageVotesForAllProductsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAvgVotes := map[string]*vote.ProductVote{
		"p1": {VotesCount: 2, Avg: 7.5},
	}
	app := &Application{
		Products: map[string]*product.Product{
			"p1": {ID: "p1", Name: "Product 1"},
		},
		voteService: &MockVoteService{
			mockAvgVotes: mockAvgVotes,
		},
	}

	router := setupRouter(app)

	// Test case: Average votes calculated
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/products/avgs", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var respAvgVotes map[string]*vote.ProductVote
	err := json.Unmarshal(w.Body.Bytes(), &respAvgVotes)
	assert.NoError(t, err)
	assert.Equal(t, mockAvgVotes, respAvgVotes)

	// Test case: No votes found
	app.voteService = &MockVoteService{}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/products/avgs", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Looks like there are no votes so far.")
}
