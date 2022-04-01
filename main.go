package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ResData struct {
	Id       string `json:"id"`
	ShortUrl string `json:"shortUrl"`
}
type ReqData struct {
	Url      string `json:"url"`
	ExpireAt string `json:"expireAt"`
}

// Check Err status
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Generate a random string of [a-zA-Z0-9]
func allowdChar() byte {
	index := 48 + rand.Intn(62)
	if index >= 58 {
		index += 7
	}
	if index >= 91 {
		index += 6
	}
	return byte(index)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = allowdChar()
	}
	return string(bytes)
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
	res.Id = randomString(6)
	res.ShortUrl = "http://localhost/" + res.Id
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
	db, err := sql.Open("mysql", "root:12345678@/go?charset=utf8")
	checkErr(err)

	//查詢資料
	rows, err := db.Query("SELECT * FROM go.links")
	checkErr(err)

	for rows.Next() {
		var Id string
		var ExpireAt string
		var OriginalUrl string
		err = rows.Scan(&Id, &ExpireAt, &OriginalUrl)
		checkErr(err)
		fmt.Println("-----------")
		fmt.Println("Id:,", Id)
		fmt.Println("ExpireAt:", ExpireAt)
		fmt.Println("OriginalUrl", OriginalUrl)
	}
	fmt.Println("-----End Of DB query------")

	server := gin.Default()

	rand.Seed(time.Now().UnixNano())

	server.POST("/api/v1/urls", rest)
	server.GET("/:id", redirect)
	server.GET("/", index)
	server.Run(":80")
}
