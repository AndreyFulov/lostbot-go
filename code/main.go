package main

import (
	"log"
	"os"
	"time"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")


func main() {
	var db DataBase
	db.InitInfo(host, port,user,password,dbname,sslmode)
	time.Sleep(5 * time.Second)
	err := db.CreateTable()
	if err != nil {
		log.Panic(err)
	}
	bot := NewBot(os.Getenv("TOKEN"), &db)
	bot.Bot()
}