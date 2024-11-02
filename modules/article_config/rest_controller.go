package artconfig

import (
	"emmanuel-guerreiro/stockgo/lib"

	"github.com/gin-gonic/gin"
)

func findOne(c *gin.Context) {
	res, err := findOneById(c.Param("id"), c)

	if err != nil {
		lib.AbortWithError(c, err)
		return
	}

	c.JSON(200, res)
}

func post(c *gin.Context) {
	var body CreateArticleConfigDto
	if err := c.ShouldBindJSON(&body); err != nil {
		lib.AbortWithError(c, err)
		return
	}

	res, err := createArticleConfig(&body, c)

	if err != nil {
		lib.AbortWithError(c, err)
		return
	}

	c.JSON(201, res)
}

func put(c *gin.Context) {
	var body ReplaceArticleConfigDto
	if err := c.ShouldBindJSON(&body); err != nil {
		lib.AbortWithError(c, err)
		return
	}
	res, err := replaceArticleConfig(c.Param("id"), &body, c)
	if err != nil {
		lib.AbortWithError(c, err)
		return
	}

	c.JSON(200, res)

}
func delete(c *gin.Context) {

	res, err := deleteArticleConfig(c.Param("id"), c)
	if err != nil {
		lib.AbortWithError(c, err)
		return
	}

	c.JSON(200, res)
}

func InitController(router *gin.RouterGroup) {
	prodConfig := router.Group("/article-config")
	prodConfig.GET("/:id", findOne)
	prodConfig.POST("", post)
	prodConfig.PUT("/:id", put)
	prodConfig.DELETE("/:id", delete)
}
