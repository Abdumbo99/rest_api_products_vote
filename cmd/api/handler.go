package main

import (
	"api_assignment/api/models/product"
	"api_assignment/api/models/vote"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// application is the handler for the requests. Additionally, it holds the necessary field for the whole application
// this includes the voteService and the products.
// Products are saved here to be used for request validation and to return them when /products is called
// purpose of saving them here insread of db is becuase products do not change frequently and to reduce calls to db
type application struct {
	Products map[string]*product.Product

	// interface for easier testing
	voteService interface {
		AllVotes() ([]*vote.Vote, error)
		PostVote(newVote *vote.Vote) (*bool, error)
		GetVotesBySessionID(sessionID string) ([]*vote.Vote, error)
		GetVotesByProductID(productID string) ([]*vote.Vote, error)
		GetAvergageVotesForAllProducts(products map[string]*product.Product) (map[string]*vote.VoteResult, error)
	}
}

// NewApp creates an istancve of the application and assigns the client passed to it as its client
func NewApp(client *mongo.Client) *application {
	prs := product.FetchProducts()

	return &application{
		Products:    prs,
		voteService: vote.VoteModel{DB: client},
	}

}

// AllProductsHandler return all products in the sysetm
func (app *application) AllProductsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		if len(app.Products) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no products so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, app.Products)
	}
}

// AllVotessHandler return all votes in the sysetm. Originally this was not required but for validating results
// without peeking at the db, this is good handler
func (app *application) AllVotessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		allVotes, err := app.voteService.AllVotes()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if len(allVotes) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no votes so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, allVotes)
	}
}

// PostVoteHandler posts a vote or updates it in the system
func (app *application) PostVoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		newVote := &vote.Vote{}
		if err := c.ShouldBindJSON(newVote); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "request is invalid. Please check your request"})
			return
		}

		// could not find the product
		if _, ok := app.Products[newVote.ProductID]; !ok {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No such product"})
			return
		}

		if newVote.Rate <= 0 || newVote.Rate > 10 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "rate must be between 0 and 10!"})
			return
		}

		session := sessions.Default(c)
		// Check if the session ID exists
		sessionID := session.Get("session_id").(string)
		newVote.SessionID = sessionID

		voteExists, err := app.voteService.PostVote(newVote)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			fmt.Println(err.Error())
			return
		}
		// if vote already exists update it
		if *voteExists {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Vote already exists, your rate of the product was updated"})
			return
		}
		// if vote does not exist save it
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Your vote has been received successfully!"})
	}
}

// GetVotesBySessionIDHanlder takes a session id and return all its votes so far
func (app *application) GetVotesBySessionIDHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionID := c.Param("id")

		votes, err := app.voteService.GetVotesBySessionID(sessionID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if len(votes) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no votes for this session so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, votes)

	}
}

// GetVotesByProductIDHanlder takes a product id and return all its votes so far
// this handler was not requested but it helps to validate data
func (app *application) GetVotesByProductIDHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {

		productID := c.Param("id")

		if _, ok := app.Products[productID]; !ok {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No such product"})
			return
		}

		votes, err := app.voteService.GetVotesByProductID(productID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if len(votes) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no votes for this product so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, votes)

	}
}

// GetAvergageVotesForAllProductsHanlder calculates the avg votes for all products across all sessions
func (app *application) GetAvergageVotesForAllProductsHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {

		avgs, err := app.voteService.GetAvergageVotesForAllProducts(app.Products)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if len(avgs) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no votes so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, avgs)

	}
}
