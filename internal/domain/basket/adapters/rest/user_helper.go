package rest

import "github.com/gin-gonic/gin"

const (
	COOKIE_USER_ID  = "USER_ID"
	DEFAULT_USER_ID = "1337"
)

// helper to simulate a user
func getUserID(c *gin.Context) string {
	userID, err := c.Cookie(COOKIE_USER_ID)
	if err != nil || userID == "" {
		return DEFAULT_USER_ID
	}

	return userID
}
