package main

import (
	"api_assignment/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) AllHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		allVotes, err := app.VoteService.All()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if len(allVotes) == 0 {
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Looks like there are no votes so far."})
			return
		}
		c.IndentedJSON(http.StatusOK, allVotes)
	}
}

func (app *application) PostVoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newVote *models.Vote
		newVote.SessionID = "" //TODO
		if err := c.ShouldBindBodyWithJSON(newVote); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "request is invalid. Please check your request"})
		}
		newVote.Product.Name = "" //TODO
		resp, err := app.VoteService.PostVote(newVote)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if resp.AlreadyExist {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Vote already exists, your rate of the product was updated"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Your vote has been received successfully!"})
	}
}

func (app *application) GetVotesByProductIDHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		votes, err := app.VoteService.GetVotesByProductID(productID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if len(votes) == 0 {
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Looks like there are no votes for this product so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, votes)

	}
}

func (app *application) GetAvergageVotesByProductHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := "" // TODO
		avgs, err := app.VoteService.GetAvergageVotesByProduct(sessionId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
			return
		}
		if len(avgs) == 0 {
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Looks like there are no votes for this product so far."})
			return
		}

		c.IndentedJSON(http.StatusOK, avgs)

	}
}
