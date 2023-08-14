package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Player struct {
	Id int64
	PlayerTGID int64
	Name string
	Level int
	Money int
}

type BusinessType struct {
	Id int
	Name string
	Price int
	Income int
}
type Business struct {
	OwnerTGID int64
	Type int
	Amount int 
}

type DataBase struct {

}
var dbInfo string

const CountOfBusinessType = 2
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
		PlayerTGID BIGINT,
		Name TEXT,
		Level INT,
		Money INT
	);`); err != nil {
        return err
    }
	time.Sleep(5 * time.Second)
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS business_type (
		Id INT UNIQUE,
		Name TEXT,
		Price INT,
		Income INT
	);`); err != nil {
        return err
    }
	time.Sleep(5 * time.Second)
	initBusinessTypes()
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS business (
		OwnerTGID BIGINT,
		Type INT REFERENCES business_type (Id),
		Amount int

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
		err := rows.Scan(&player.Id, &player.PlayerTGID, &player.Name, &player.Level, &player.Money)
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

func initBusinessTypes() error{
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
        return err
    }
    defer db.Close()

	data := `INSERT INTO business_type (Id, Name, Price, Income) VALUES (1,'Суши-Эдо',100, 10)`
	if _, err = db.Exec(data); err != nil {
		return err
	}
	data = `INSERT INTO business_type (Id, Name, Price, Income) VALUES (2,'Лост-Дот',1000, 100)`
	if _, err = db.Exec(data); err != nil {
		return err
	}
	data = `DELETE FROM business_type WHERE Id > $1`
	if _, err = db.Exec(data, CountOfBusinessType); err != nil {
		return err
	}
	return nil
}


func(d *DataBase) AddBusinessToPlayer(p Player, biz_type int) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
        return err
    }
    defer db.Close()
	rows, err := db.Query("SELECT * FROM business WHERE OwnerTGID = $1 AND Type = $2",p.PlayerTGID,biz_type);
	if rows.Next() {
		data := `UPDATE business SET Amount = Amount + 1 WHERE OwnerTGID = $1 AND Type = $2`
		if _, err = db.Exec(data,p.PlayerTGID, biz_type); err != nil {
			return err
		}
	return nil
	}else {
		data := `INSERT INTO business (OwnerTGID, Type, Amount) VALUES ($1, $2, $3)`
		if _, err = db.Exec(data, p.PlayerTGID, biz_type, 1); err != nil {
			return err
		}
	}
	return nil
}

func(d *DataBase) GetPlayerBuisnesses(p Player) ([]Business, error) {
	//Сделать функционал
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
        return nil, err
    }
    defer db.Close()
	var playerBizes []Business
	for i := 1; i <=2;i++ {
		rows, err := db.Query("SELECT * FROM business WHERE OwnerTGID = $1 AND Type = $2", p.PlayerTGID,i)
		if err != nil {
			return nil, err
		}
		if rows.Next() {
			var biz Business
			err := rows.Scan(&biz.OwnerTGID,&biz.Type,&biz.Amount)
			if err != nil {
				return nil, err
			}
			playerBizes = append(playerBizes, biz)
		}
	}
	return playerBizes, nil
}

func (d *DataBase) GetBusinessTypeById(biz_type int) (BusinessType, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
        return (BusinessType{}), err
    }
    defer db.Close()
	rows, err := db.Query("SELECT * FROM business_type WHERE Id = $1",biz_type);
	if rows.Next() {
		var biz BusinessType
		err := rows.Scan(&biz.Id, &biz.Name, &biz.Price, &biz.Income)
		if err != nil {
			return (BusinessType{}), err
		}
		return biz, err
	}else {
		return (BusinessType{}), nil
	}
}