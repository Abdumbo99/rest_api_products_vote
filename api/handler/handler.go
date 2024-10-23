package handler

import (
	"api_assignment/api/models/product"
	"api_assignment/api/models/vote"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// @title Voting API
// @version 1.0
// @description API for voting on products in a session-based system.
// @host localhost:8080
// @BasePath /

// application is the handler for the requests. Additionally, it holds the necessary field for the whole application
// this includes the voteService and the products.
// Products are saved here to be used for request validation and to return them when /products is called
// purpose of saving them here instead of db is becuase products do not change frequently and to reduce calls to db
type application struct {
	Products map[string]*product.Product

	// interface for easier testing
	voteService interface {
		AllVotes() ([]*vote.Vote, error)
		PostVote(newVote *vote.Vote) (*bool, error)
		GetVotesBySessionID(sessionID string) ([]*vote.Vote, error)
		GetVotesByProductID(productID string) ([]*vote.Vote, error)
		GetAverageVotesForAllProducts(products map[string]*product.Product) (map[string]*vote.VoteResult, error)
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

// @Summary Get all products
// @Description Retrieves all the available products in the system.
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string]*product.Product
// @Failure 204 {object} map[string]string
// @Router /products [get]
func (app *application) AllProductsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		if len(app.Products) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no products so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, app.Products)
	}
}

// @Summary Get all votes
// @Description Retrieves all the votes from the system.
// @Tags votes
// @Accept json
// @Produce json
// @Success 200 {array} vote.Vote
// @Failure 500 {object} map[string]string
// @Router /votes [get]
func (app *application) AllVotessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		allVotes, err := app.voteService.AllVotes()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		if len(allVotes) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no votes so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, allVotes)
	}
}

// @Summary Post or update a vote
// @Description Posts a new vote or updates an existing vote based on the session and product.
// @Tags votes
// @Accept json
// @Produce json
// @Param vote body vote.Vote true "Vote to post or update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /votes [post
func (app *application) PostVoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		newVote := &vote.Vote{}
		if err := c.ShouldBindJSON(newVote); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "request is invalid. Please check your request"})
			fmt.Fprintln(os.Stderr, err.Error())
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

// @Summary Get votes by session ID
// @Description Retrieves all votes for a given session ID.
// @Tags votes
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {array} vote.Vote
// @Failure 500 {object} map[string]string
// @Router /votes/session/{id} [get]
func (app *application) GetVotesBySessionIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionID := c.Param("id")

		votes, err := app.voteService.GetVotesBySessionID(sessionID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		if len(votes) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no votes for this session so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, votes)

	}
}

// @Summary Get votes by product ID
// @Description Retrieves all votes for a given product ID.
// @Tags votes
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {array} vote.Vote
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /votes/product/{id} [get]
func (app *application) GetVotesByProductIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		productID := c.Param("id")

		if _, ok := app.Products[productID]; !ok {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No such product"})
			return
		}

		votes, err := app.voteService.GetVotesByProductID(productID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		if len(votes) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no votes for this product so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, votes)

	}
}

// @Summary Get average votes for all products
// @Description Calculates and retrieves the average votes for all products across all sessions.
// @Tags votes
// @Accept json
// @Produce json
// @Success 200 {object} map[string]vote.VoteResult
// @Failure 500 {object} map[string]string
// @Router /products/avgs [get]
func (app *application) GetAverageVotesForAllProductsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		avgs, err := app.voteService.GetAverageVotesForAllProducts(app.Products)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		if len(avgs) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no votes so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, avgs)

	}
}
