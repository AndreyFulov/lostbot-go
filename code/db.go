package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Player struct {
	Id int
	PlayerTGID int64
	Name string
	Level int
	Money int
}

type DataBase struct {

}
var dbInfo string


func (d *DataBase) InitInfo(host, port,user,password,dbname,sslmode string) {
	dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

}

//Создаем таблицу users в БД при подключении к ней
func(d *DataBase) CreateTable() error {

    //Подключаемся к БД
    db, err := sql.Open("postgres", dbInfo)
    if err != nil {
        return err
    }
    defer db.Close()

    //Создаем таблицу users
    if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS players (
		Id SERIAL PRIMARY KEY,
		PlayerTGID INT,
		Name TEXT,
		Level INT,
		Money INT
	);`); err != nil {
        return err
    }

    return nil
}

func(d *DataBase) CreateUser(p Player) error{
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
        return err
    }
    defer db.Close()

	data := `INSERT INTO players (Name, PlayerTGID, Level, Money) VALUES ($1, $2, 1, 100)`
	if _, err = db.Exec(data,p.Name, p.PlayerTGID); err != nil {
		return err
	}
	return nil
}

func(d DataBase) GetAllPlayers() ([]Player, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
        return nil,err
    }
    defer db.Close()
	rows, err := db.Query("SELECT * FROM players");
	if err != nil{
		return nil, err
	}
	var players []Player
	for rows.Next() {
		var player Player
		err := rows.Scan(&player.Id, &player.PlayerTGID, &player.Level, &player.Money)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

func(d *DataBase) GetPlayerByTGId(id int64) (Player, error){
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return Player{}, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT * FROM players WHERE PlayerTGID = $1;`, id)
	if err != nil {
		return  Player{}, err
	}
	var player Player
	for rows.Next(){
	err := rows.Scan(&player.Id, &player.PlayerTGID, &player.Name, &player.Level, &player.Money)
	if err != nil {
		return Player{},nil
	}
}  
return player, nil
}

func(d *DataBase) ChangePlayerMoney(id int64, newMoney int) error{
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	data := `UPDATE players SET Money = $1 WHERE PlayerTGID = $2`
	if _, err = db.Exec(data,newMoney, id); err != nil {
		return err
	}
	return nil
}

func(d *DataBase) GetTopPlayerByMoney()([]Player, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
        return nil,err
    }
    defer db.Close()
	rows, err := db.Query("SELECT * FROM players ORDER BY Money DESC LIMIT 10");
	if err != nil{
		return nil, err
	}
	var players []Player
	for rows.Next() {
		var player Player
		err := rows.Scan(&player.Id, &player.PlayerTGID,&player.Name, &player.Level, &player.Money)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}