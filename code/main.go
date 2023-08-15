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
	db.CreateTable()
	bot := NewBot(os.Getenv("TOKEN"), &db)
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				calcPlayerIncomeByBiz(&db)
				log.Printf("ОПА, ВСЕМ БАБКИ!!!")
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

	bot.Bot()
}

func calcPlayerIncomeByBiz(db *DataBase) {
	players, err := db.GetAllPlayers()
	if err != nil {
		log.Print(err.Error())
	}
	for _, p := range players {
		bizes, err := db.GetPlayerBuisnesses(p)
		if err != nil {
			log.Print(err.Error())
		}
		for _, b := range bizes {
			t,err := db.GetBusinessTypeById(b.Type)
			if err != nil {
				log.Print(err.Error())
			}
			db.ChangePlayerMoney(p.PlayerTGID,p.Money + (t.Income * b.Amount))
			log.Printf("Начисленные деньги: %s за тип бизнеса - %s в колличестве %s", p.PlayerTGID,t.Id,b.Amount)
		}
	}
}
