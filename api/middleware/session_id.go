package middleware

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CheckSession is a middleware function to check and create a session if it doesn't exist
func CheckSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// Check if the session ID exists
		sessionID := session.Get("session_id")
		if sessionID == nil {
			// create a unique, new session
			newSessionID := uuid.New().String()

			// Set the session ID in the session data
			session.Set("session_id", newSessionID)

			// Save the session
			session.Save()

			fmt.Println("New session created with ID:", newSessionID)
		} else {
			fmt.Println("Existing session found with ID:", sessionID)
		}

		// Continue
		c.Next()
	}
}
