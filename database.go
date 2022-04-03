package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var dbCahce *redis.Client

const (
	MaxLifetime  int  = 600 //seconds
	MaxOpenConns int  = 5
	MaxIdleConns int  = 2
	ParseTime    bool = true
)

// Check Err status
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func connectMySQL() {
	// DB config
	config := mysql.Config{
		User:                 "b175043b2b81b6",
		Passwd:               "a79070ee",
		Addr:                 "us-cdbr-east-05.cleardb.net:3306",
		Net:                  "tcp",
		DBName:               "heroku_2a40d52308e3ec4",
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
	}

	// Connect DB
	dbConnect, err := sql.Open("mysql", config.FormatDSN())
	checkErr(err)

	// Config DB pool
	dbConnect.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	dbConnect.SetMaxOpenConns(MaxOpenConns)
	dbConnect.SetMaxIdleConns(MaxIdleConns)

	// Check connection
	err = dbConnect.Ping()
	checkErr(err)

	db = dbConnect
}

func RedisNewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	checkErr(err)

	fmt.Println("Pong response:", pong)

	// Realize redis.Client
	dbCahce = client
}
