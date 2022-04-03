package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

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

const (
	TimeFormat string        = "2006-01-02 15:04:05"
	CacheTime  time.Duration = 10
)

// Generate a random string of [a-zA-Z0-9]
func allowedChar() byte {
	index := 48 + rand.Intn(62)
	if index >= 58 {
		index += 7
	}
	if index >= 91 {
		index += 6
	}
	return byte(index)
}

// Generate random Id
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = allowedChar()
	}
	return string(bytes)
}

// Check Err status
func checkRestErr(c *gin.Context, err error) {
	if err != nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}
}

// Check Err status
func checkRedirectErr(c *gin.Context, err error) {
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
}

func rest(c *gin.Context) {
	// Get parameter
	req := new(ReqData)
	err := c.BindJSON(&req)
	checkRestErr(c, err)

	shortId := ""
	for count := 1; count != 0; {
		// Create new shortId
		shortId = randomString(6)

		// Check usable of shortId
		err := db.QueryRow("SELECT COUNT(*) FROM `links` WHERE `Id` = ?", shortId).Scan(&count)
		checkRestErr(c, err)
	}

	// Prepare insert statement
	stmt, err := db.Prepare("INSERT links SET `Id`=?,`ExpireDate`=?,`OriginalUrl`=?")
	checkRestErr(c, err)

	// Parameter's time parse
	expireTime, err := time.Parse(time.RFC3339, req.ExpireAt)
	checkRestErr(c, err)

	// Execute insert data
	result, err := stmt.Exec(shortId, expireTime.Format(TimeFormat), req.Url)
	checkRestErr(c, err)

	_, err = result.LastInsertId()
	checkRestErr(c, err)

	// Redis cache write
	dbCahce.HSet(shortId, "ExpireDate", expireTime.Format(TimeFormat))
	dbCahce.HSet(shortId, "OriginalUrl", req.Url)

	// JSON response data
	res := new(ResData)
	res.Id = shortId
	res.ShortUrl = "http://localhost/" + res.Id
	c.JSON(http.StatusOK, res)
}

func redirect(c *gin.Context) {
	id := c.Param("id")
	var ExpireAt time.Time
	var OriginalUrl string

	// Get redis cahce status
	expireDateInCache, err := dbCahce.HExists(id, "ExpireDate").Result()
	checkRedirectErr(c, err)
	originalUrlInCache, err := dbCahce.HExists(id, "OriginalUrl").Result()
	checkRedirectErr(c, err)
	invalidInCache, err := dbCahce.HExists(id, "Invalid").Result()
	checkRedirectErr(c, err)
	// Validate cache status
	if expireDateInCache && originalUrlInCache {
		// Stored shortUrl
		ExpireDate, err := dbCahce.HGet(id, "ExpireDate").Result()
		checkRedirectErr(c, err)
		ExpireAt, err = time.Parse(TimeFormat, ExpireDate)
		checkRedirectErr(c, err)

		nowTime := time.Now()
		if ExpireAt.Before(nowTime) {
			// Validate shortUrl expiration
			fmt.Println("Invalid Time")
			c.Status(http.StatusNotFound)
			return
		} else {
			// Redirect to the originalurl by cache
			OriginalUrl, err = dbCahce.HGet(id, "OriginalUrl").Result()
			checkRedirectErr(c, err)
			c.Redirect(http.StatusMovedPermanently, OriginalUrl)
			return
		}
	} else if invalidInCache {
		// Refuse invalid request by cache in 10s
		fmt.Println("Invalid in Redis")
		dbCahce.Expire(id, CacheTime*time.Second).Result()
		c.Status(http.StatusNotFound)
		return
	}
	fmt.Printf("Id %s does not in cache\n", id)

	// Select data by RDB
	rows, err := db.Query("SELECT `ExpireDate`,`OriginalUrl` FROM `links` WHERE `Id` = ?", id)
	checkRedirectErr(c, err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&ExpireAt, &OriginalUrl)
		checkRedirectErr(c, err)
	}

	nowTime := time.Now()
	if ExpireAt.Before(nowTime) {
		// Validate shortUrl expiration
		fmt.Println("Invalid Time or url not exist")

		// Remember the invalid id in redis in 10s
		dbCahce.HSet(id, "Invalid", 1)
		dbCahce.Expire(id, CacheTime*time.Second).Result()

		c.Status(http.StatusNotFound)
		return
	}
	// Store information in cache to increase performance
	dbCahce.HSet(id, "ExpireDate", ExpireAt.Format(TimeFormat))
	dbCahce.HSet(id, "OriginalUrl", OriginalUrl)
	c.Redirect(http.StatusMovedPermanently, OriginalUrl)
}

func index(c *gin.Context) {
	c.Status(http.StatusNotFound)
}
