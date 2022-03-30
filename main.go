package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResData struct {
	Id       string `json:"id"`
	ShortUrl string `json:"shortUrl"`
}
type ReqData struct {
	Url      string `json:"url"`
	ExpireAt string `json:"expireAt"`
}

func rest(c *gin.Context) {
	req := new(ReqData)
	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("ExpireAt", req.ExpireAt)
	fmt.Println("Url", req.Url)

	res := new(ResData)
	c.JSON(http.StatusOK, res)
}

func redirect(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	c.JSON(http.StatusOK, id)
}

func index(c *gin.Context) {
	c.Status(http.StatusNotFound)
}

func main() {
	server := gin.Default()
	server.POST("/api/v1/urls", rest)
	server.Any("/:id", redirect)
	server.GET("/", index)
	server.Run(":8888")
}
