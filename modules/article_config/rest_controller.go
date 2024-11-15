package artconfig

import (
	"emmanuel-guerreiro/stockgo/lib"
	"strconv"

	"github.com/gin-gonic/gin"
)

func findOne(c *gin.Context) {
	res, err := FindOneById(c.Param("id"), c)

	if err != nil {
		lib.AbortWithError(c, err)
		return
	}

	c.JSON(200, res)
}

func find(c *gin.Context) {
	//TODO: Implement sorting
	//TODO: Implement filtering
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	res, err := findAll(page, size, c)
	if err != nil {
		lib.AbortWithError(c, err)
		return
	}

	if res == nil {
		res = make([]ArticleConfig, 0)
	}

	jsonRes := ArticleFindResponsePaginated{
		Data:   res,
		Page:   page,
		Length: len(res),
	}

	c.JSON(200, jsonRes)
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
	prodConfig.GET("", find)
	prodConfig.PUT("/:id", put)
	prodConfig.DELETE("/:id", delete)
}
