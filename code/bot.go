package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot_token string
	db *DataBase
}

func NewBot(token string, db *DataBase) *TelegramBot {
	return &TelegramBot{
		bot_token: token,
		db: db,
	}
}




func (tg *TelegramBot) Bot() {
	bot, err := tgbotapi.NewBotAPI(tg.bot_token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Используется на: %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			userInput := strings.Fields(update.Message.Text)

			if len(userInput) > 0 {
			//Логика создания	 нового пользователя
			if userInput[0] == "/start"{
				p, err := tg.db.GetPlayerByTGId(update.Message.From.ID)
				if err != nil {
					log.Panic(err.Error())
				}
		
				if p.PlayerTGID != 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Вы уже зарегестрированы! Ваше имя - %s", p.Name))
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg) 
					
				}else {
					if len(userInput) == 2 {
						if len(userInput[1]) <= 255 {
							p:= Player{
								Name: userInput[1],
								PlayerTGID: update.Message.From.ID,
								Level: 1,
							}
							tg.db.CreateUser(p)
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("✅Создан новый пользователь! %s",p.Name) )
							msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}else{
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Имя слишком длинное!")
							msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}
					}else if len(userInput) == 1 {
					p:= Player{
						Name: update.Message.From.FirstName,
						PlayerTGID: update.Message.From.ID,
						Level: 1,
					}
					tg.db.CreateUser(p)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("✅Создан новый пользователь: %s",p.Name) )
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
				}
				}
				//Логика Баланса
				if userInput[0] == "/balance" {
					p, err := tg.db.GetPlayerByTGId(update.Message.From.ID)
					if err != nil {
						log.Panic(err.Error())
					}
					if p != (Player{}){
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Баланс: %s$", strconv.Itoa(p.Money)))
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас еще нет аккаунта, чтоюы создать его '/start [имя]' ")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}
				}
				//Тестовая логика работы
				if userInput[0] == "/work" {
					p, err := tg.db.GetPlayerByTGId(update.Message.From.ID)
					if err != nil {
						log.Panic(err.Error())
					}
					if p != (Player{}){
						tg.db.ChangePlayerMoney(p.PlayerTGID,p.Money+10)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Вы заработали: 10$\nТеперь ваш баланс: %s$", strconv.Itoa(p.Money + 10)))
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас еще нет аккаунта, чтоюы создать его '/start [имя]' ")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}
				}
				//Логика казино
				if userInput[0] == "/casino" {
					if len(userInput) == 2 {
						bet, err := strconv.Atoi(userInput[1])
						if err != nil {
							msg:=tgbotapi.NewMessage(update.Message.Chat.ID, "Ставка введена неверно!")
							msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}else {
							p, err := tg.db.GetPlayerByTGId(update.Message.From.ID)
							if err != nil {
								log.Panic(err.Error())
							}
							if p != (Player{}){
								if p.Money >= bet {
									s1 := rand.NewSource(time.Now().UnixNano())
								r1 := rand.New(s1)
								chance := r1.Int63n(100)
								if chance >= 60 {
									tg.db.ChangePlayerMoney(p.PlayerTGID,p.Money+bet)
									msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("✅Вы выиграли!\nТеперь ваш баланс: %s$ (+%s$))", strconv.Itoa(p.Money + bet), userInput[1]))
									msg.ReplyToMessageID = update.Message.MessageID
									bot.Send(msg)
								}else {
									tg.db.ChangePlayerMoney(p.PlayerTGID,p.Money-bet)
									msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("❌Вы проиграли!\nТеперь ваш баланс: %s$ (-%s$)", strconv.Itoa(p.Money - bet), userInput[1]))
									msg.ReplyToMessageID = update.Message.MessageID
									bot.Send(msg)
								}
								}else{
									msg:= tgbotapi.NewMessage(update.Message.Chat.ID, "У вас недостаточно средств!")
									msg.ReplyToMessageID = update.Message.MessageID
									bot.Send(msg)
								}
							}else {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас еще нет аккаунта, чтоюы создать его '/start [имя]' ")
								msg.ReplyToMessageID = update.Message.MessageID
								bot.Send(msg)
							}
						}
					}else{
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите '/casino [ставка]'")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}
				}

				//Логика показа профиля
				if userInput[0] == "/profile" {
					p, err := tg.db.GetPlayerByTGId(update.Message.From.ID)
					if err != nil {
						log.Panic(err.Error())
					}
					if p != (Player{}) {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("📄Имя: %s\n💵Деньги: %s$\n✨Уровень: %s",p.Name, strconv.Itoa(p.Money), strconv.Itoa(p.Level)))
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас еще нет аккаунта, чтоюы создать его '/start [имя]' ")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}
				}

				//Показ топа по деньгам
				if userInput[0] == "/top" {
					players, err:=tg.db.GetTopPlayerByMoney()
					if err != nil {
						log.Panic(err.Error())
					}
					data := ""
					for i, p := range players{
						data += fmt.Sprintf("%s. %s - %s$\n", strconv.Itoa(i+1), p.Name, strconv.Itoa(p.Money))
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, data)
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
				//Логика покупки бизнеса (СДЕЛАТЬ КРУЧЕ ПОТОМ, ДА)
				if userInput[0] == "/buy" {
					p, err := tg.db.GetPlayerByTGId(update.Message.From.ID)
					if err != nil {
						log.Panic(err.Error())
					}
					if p != (Player{}) {
						if len(userInput) == 2 {
							biz_id, err := strconv.Atoi(userInput[1])
							if err != nil || biz_id <= 0 || biz_id > CountOfBusinessType {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌Неправильно введен тип бизнеса!")
								msg.ReplyToMessageID = update.Message.MessageID
								bot.Send(msg)
							}else{
								d, err := tg.db.GetBusinessTypeById(biz_id)
								if err != nil {
									log.Fatal(err.Error())
								}
								if p.Money >= d.Price {
									err := tg.db.ChangePlayerMoney(p.PlayerTGID, p.Money - d.Price)
									if err != nil {
										msg:= tgbotapi.NewMessage(update.Message.Chat.ID, "❌Упс! Что-то пошло не так!")
										msg.ReplyToMessageID = update.Message.MessageID
										bot.Send(msg)

									}else{
										err := tg.db.AddBusinessToPlayer(p,biz_id)
										if err != nil {
											msg:= tgbotapi.NewMessage(update.Message.Chat.ID, "❌Упс! Что-то пошло не так!")
											msg.ReplyToMessageID = update.Message.MessageID
											bot.Send(msg)
										}else{
											msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("✅Вы купили бизнес: %s!", d.Name))
											msg.ReplyToMessageID = update.Message.MessageID
											bot.Send(msg)
										}
									}
								}else{
									msg:= tgbotapi.NewMessage(update.Message.Chat.ID, "❌У вас недостаточно денег, чтобы купить этот бизнес!")
									msg.ReplyToMessageID = update.Message.MessageID
									bot.Send(msg)
								}
							}
						}else{
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌Некоректное использование комманды, введите '/buy [тип бизнеса]' ")
							msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}
					}else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌У вас еще нет аккаунта, чтоюы создать его '/start [имя]' ")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}
				}
				//Логика просмотра бизнесов игрока
				if userInput[0] == "/mybizes" {
					p, err := tg.db.GetPlayerByTGId(update.Message.From.ID)
					if err != nil {
						log.Panic(err.Error())
					}
					if p != (Player{}) {
						b, err := tg.db.GetPlayerBuisnesses(p)
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌Упс! Что-то пошло не так!")
							msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}else {
							s := "Ваши бизнесы: \n"
							for i, bi := range b {
								t, err := tg.db.GetBusinessTypeById(bi.Type)
								if err != nil {
									msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌Упс! Что-то пошло не так!")
									msg.ReplyToMessageID = update.Message.MessageID
									bot.Send(msg)
								}
								s += fmt.Sprintf("%s. %s - %s шт.\n", strconv.Itoa(i +1),t.Name,strconv.Itoa(bi.Amount))
							}
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, s)
							msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}
					}else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌У вас еще нет аккаунта, чтоюы создать его '/start [имя]' ")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}
				}
				//Логика показа всех типов бизнеса
				if userInput[0] == "/bizes" {
					s := "Бизнесы: \n"
					for i := 1; i <= CountOfBusinessType; i ++ {
						t, err := tg.db.GetBusinessTypeById(i)
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌Упс! Что-то пошло не так!")
							msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}
						s += fmt.Sprintf("%s\nID:%s\nСтоимость:%s$\nДоход:%s$\n------\n",t.Name,strconv.Itoa(t.Id),strconv.Itoa(t.Price),strconv.Itoa(t.Income))
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, s)
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
				if userInput[0] == "/help" {
					s := "/start - Создать профиль\n/work - Заработать несчастные 10 долларов\n/buy 'ID' - Купить бизнес\n/bizes - Список доступных бизнесов\n/mybizes - Ваши бизнесы\n/casino 'ставка' - Сыграть в казино\n----\nЕсли будут новые команды, они появятся тут\nhttps://github.com/AndreyFulov/lostbot-go/blob/main/CHANGELOG.md"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, s)
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}


				
		}
	}
}
}