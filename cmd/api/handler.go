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

type application struct {
	Products    map[string]*product.Product
	voteService interface {
		AllVotes() ([]*vote.Vote, error)
		PostVote(newVote *vote.Vote) (*vote.PostResponse, error)
		GetVotesBySessionID(sessionID string) ([]*vote.Vote, error)
		GetVotesByProductID(productID string) ([]*vote.Vote, error)
		GetAvergageVotesForAllProducts(products map[string]*product.Product) (map[string]*vote.VoteResult, error)
	}
}

func NewApp(client *mongo.Client) *application {
	prs := product.FetchProducts()

	return &application{
		Products:    prs,
		voteService: vote.VoteModel{DB: client},
	}

}

func (app *application) AllProductsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(app.Products) == 0 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Looks like there are no products so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, app.Products)
	}
}

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

func (app *application) PostVoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		// Check if the session ID exists
		sessionID := session.Get("session_id").(string)

		newVote := &vote.Vote{}
		if err := c.ShouldBindJSON(newVote); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "request is invalid. Please check your request"})
			return
		}

		if _, ok := app.Products[newVote.ProductID]; !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No such product"})
			return
		}

		if newVote.Rate <= 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "rate must be higher than 0!"})
			return
		}
		if newVote.ProductID == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Prodcut does not exist!"})
			return
		}

		newVote.SessionID = sessionID

		resp, err := app.voteService.PostVote(newVote)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			fmt.Println(err.Error())
			return
		}
		if resp.AlreadyExist {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Vote already exists, your rate of the product was updated"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Your vote has been received successfully!"})
	}
}

func (app *application) GetVotesBySessionIDHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		fmt.Println(sessionID)
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

func (app *application) GetVotesByProductIDHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		fmt.Println(productID)

		if _, ok := app.Products[productID]; !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No such product"})
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

func (app *application) GetAvergageVotesForAllProductsHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		avgs, err := app.voteService.GetAvergageVotesForAllProducts(app.Products)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if len(avgs) == 0 {
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Looks like there are no votes so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, avgs)

	}
}
