package lib

import (
	"github.com/gin-gonic/gin"
)

func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}
