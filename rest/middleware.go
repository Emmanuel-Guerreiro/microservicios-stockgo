package rest

import (
	"emmanuel-guerreiro/stockgo/lib"
	"emmanuel-guerreiro/stockgo/security"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProtectedMiddleware(c *gin.Context) {
	_, err := validateToken(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	// ctx := GinCtx(c)
	// c.Set("logger", log.Get(ctx...).WithField(log.LOG_FIELD_USER_ID, user.ID))
}

// get token from Authorization header
func HeaderToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(tokenString, "bearer ") != 0 {
		return "", lib.UnauthorizedError
	}
	return tokenString[7:], nil
}

func validateToken(c *gin.Context) (*security.User, error) {
	tokenString, err := HeaderToken(c)
	if err != nil {
		return nil, lib.UnauthorizedError
	}

	// ctx := GinCtx(c)
	ctx := c
	user, err := security.Validate(tokenString, ctx)
	if err != nil {
		return nil, lib.UnauthorizedError
	}

	c.Set("tokenString", tokenString)
	c.Set("user", *user)

	return user, nil
}
