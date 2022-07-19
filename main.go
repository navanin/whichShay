package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Shays struct {
	ID   uint
	shay string
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// Connect to database
	db, err := sql.Open("sqlite3", "./sqlite.db")
	checkErr(err)

	// defer close
	//defer db.Close()

	bot, _ := tgbotapi.NewBotAPI(TOKEN)
	bot.Debug = false
	log.Printf("Authorized on account %s, with ID %d", bot.Self.UserName, bot.Self.ID)
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 20
	updates := bot.GetUpdatesChan(update)

	for update := range updates {
		if update.Message == nil { // Игнорирование всего, кроме текстовых сообщений
			continue
		}
		if !update.Message.IsCommand() { // Игнорирование всего, кроме команд
			continue
		}

		// Создание сообщения. Так как наполнять ее пока нечем,
		// оставляем поле для текста пустым.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = tgbotapi.ModeMarkdown

		// Извлечение команды из полученного сообщения
		switch update.Message.Command() {

		case "help":
			msg.Text = "Для того, чтобы узнать, кто сегодня Маратович введите \n*/get*\n\n" +
				"Для добавления новой вариации имени введите \n*/add Вильдан Маратович*."

		case "get":
			if !getDay() {
				msg.Text = fmt.Sprintf("_Я уже говорил_. \n\n*%s*, - сегодня величайшего зовут так.", shayName)
				break
			}
			for {
				id = randID(db)
				if id != latest {
					break
				}
			}

			shayName = getShay(db, id)

			if id == 1 {
				msg.Text = fmt.Sprintf("*Вот это да!*\n" +
					"Сегодня Вильдан Маратович и есть Вильдан Маратович. Удивительно.")
			} else {
				msg.Text = fmt.Sprintf("*%s*, - сегодня величайшего зовут так. @FeldwebelWillman, слышал?", shayName)
			}

		case "add":
			newName := removeCommand(update.Message.Text)
			if countWords(newName) != 2 {
				msg.Text = "Некорректное имя."
				continue
			} else {
				newShay := Shays{shay: newName}
				if addShay(db, newShay) {
					msg.Text = fmt.Sprintf("Добавлено новое имя: *%s*!", newName)
				} else {
					msg.Text = "Произошла ошибка."
				}
			}

		default:
			msg.Text = "Не понял."
		}

		_, err := bot.Send(msg)
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getDay() bool {
	cmd := exec.Command("date", "+%d")
	stdout, err := cmd.Output()
	osDate := string(stdout)
	checkErr(err)

	if osDate != day {
		day = osDate
		return true
	} else {
		return false
	}

}

func removeCommand(msgText string) string {
	name := strings.ReplaceAll(msgText, "/add ", "")
	return name
}

func countWords(name string) int {
	return len(strings.Fields(name))
}

func addShay(db *sql.DB, newShay Shays) bool {

	//  check for the existence of a name in the db

	sqlStatement := "INSERT INTO shays VALUES(?, ?);"
	stmt, _ := db.Prepare(sqlStatement)
	_, err := stmt.Exec(nil, newShay.shay)
	stmt.Close()
	if err != nil {
		return false
	} else {
		return true
	}
}

func randID(db *sql.DB) int {
	minID := 1
	var maxID int
	var randedID int
	sqlStatement := "SELECT id FROM shays WHERE id=(SELECT MAX(id) FROM shays);"
	row := db.QueryRow(sqlStatement, nil)
	row.Scan(&maxID)
	randedID = minID + rand.Intn(maxID-minID)
	fmt.Printf("Rand from %d to %d. Got %d.\n", minID, maxID, randedID)
	return randedID
}

func getShay(db *sql.DB, id int) string {
	var name string
	sqlStatement := "SELECT shay FROM shays WHERE id=?;"
	row := db.QueryRow(sqlStatement, id)
	row.Scan(&name)
	return name
}
