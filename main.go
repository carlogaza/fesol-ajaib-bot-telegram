package main

import (
	"os"
	"log"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/tucnak/telebot"
	"strconv"
)

// Configuration for read json file that contains API KEY
type Configuration struct {
	API_KEY		string   `json:"API_KEY"`
}

// Function to load json configuration file and parse to configuration struct
func LoadConfig(path string) Configuration {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Config file is missing. ", err)
	}

	var config Configuration
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln("Config parse error: ", err)
	}

	return config
}

// Calculate the sum of string
func calculate(words string) int {
	var sum int = 0

	byteArray := []byte(words)
	for i := 0; i < len(byteArray); i++ {
		sum += int(byteArray[i])
	}

	return sum
}

// Mention user with username, if empty mention with name.
func mentionUser(pesan telebot.Message) string {
	username := ""

	if pesan.Sender.Username == "" {
		username = pesan.Sender.FirstName + " " + pesan.Sender.LastName
	} else {
		username = "@" + pesan.Sender.Username
	}

	return username
}

// Generate the answer of question
func answer(query string) string {
	var ans string

	words := strings.Split(query, " ")

	if strings.Contains(strings.ToLower(query), "atau") && len(words) > 3 {
		choose := ""
		for i := 0; i < len(words); i++ {
			if strings.ToLower(words[i]) == "atau"{
				if i == len(words)-1 {
					println(i)
					println(len(words))
					choose = "Pertanyaan kurang lengkap!"
					i += 99999
				} else {
					choose = words[i-1]
					mod := (calculate(choose) + calculate(words[i+1])) % 2
					if mod == 0 {
						if calculate(choose) < calculate(words[i+1]) {
							choose = words[i+1]
						}
					} else {
						if calculate(words[i-1]) > calculate(words[i+1]) {
							choose = words[i+1]
						}
					}
				}
			}
		}
		ans = choose

	// This is just for make funny. :D
	} else if (strings.Contains(strings.ToLower(query), "fesol") &&
		strings.Contains(strings.ToLower(query), "jembut")) ||
		(strings.Contains(strings.ToLower(query), "fesol") &&
		strings.Contains(strings.ToLower(query), "jmbt")) {
		ans = "ABSOLUTELY!"

	// For 'when' question
	} else if strings.Contains(strings.ToLower(query), "kapan") {
		switch calculate(query) % 11 {
		case 0:
			ans = "Sekarang!"
		case 1:
			ans = "Besok!"
		case 2:
			ans = "Mungkin minggu depan."
		case 3:
			ans = "Mungkin bulan depan."
		case 4:
			ans = "Bisa jadi tahun depan."
		case 5:
			ans = strconv.Itoa(calculate(query) % 3 + 2) + " tahun lagi mungkin."
		case 6:
			ans = "Aku tidak tahu."
		case 7:
			ans = "Mungkin suatu hari nanti."
		case 8:
			ans = "Saat lebaran kuda."
		case 9:
			ans = "Tidak akan pernah!"
		case 10:
			ans = "Itu mustahil!"
		}

	// For question `how long ...`
	} else if strings.Contains(strings.ToLower(query), "berapa lama lagi") ||
		(strings.Contains(strings.ToLower(query), "berapa lama") && strings.Contains(strings.ToLower(query), "akan")) {
		switch calculate(query) % 5 {
		case 0:
			ans = "Selamanya!"
		case 1:
			ans = "Beberapa hari lagi."
		case 2:
			ans = "Bisa jadi beberapa minggu lagi."
		case 3:
			ans = "Mungkin beberapa bulan lagi."
		case 4:
			ans = "Sampai saatnya tiba nanti."
		case 5:
			ans = "Beberapa menit lagi."
		case 6:
			ans = "Beberapa jam lagi."
		}

	// For question `how long` only
	} else if strings.Contains(strings.ToLower(query), "berapa lama") {
		switch calculate(query) % 11 {
		case 0:
			ans = "Satu jam."
		case 1:
			ans = "Dua jam."
		case 2:
			ans = "Beberapa jam."
		case 3:
			ans = "Satu hari."
		case 4:
			ans = "Dua hari."
		case 5:
			ans = "Beberapa hari."
		case 6:
			ans = "Satu minggu."
		case 7:
			ans = "Satu bulan."
		case 8:
			ans = "Beberapa bulan."
		case 9:
			ans = "Satu tahun."
		case 10:
			ans = "Beberapa tahun."
		}

	// For general question
	} else {
		switch calculate(query) % 7 {
		case 0:
			ans = "Tidak!"
		case 1:
			ans = "Iya!"
		case 2:
			ans = "Mungkin."
		case 3:
			ans = "Bisa jadi."
		case 4:
			ans = "Aku tidak yakin."
		case 5:
			ans = "Aku ragu."
		case 6:
			ans = "Mungkin suatu hari nanti."
		}
	}
	return ans;
}

// Just say hello to user
func hi(bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat, "Hello, " + mentionUser(message) + "!", nil)
}

// /help function
func help(bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		"Tanyakan apapun pada Fesol Ajaib!\n" +
		"Ketik '/tanya <pertanyaan kamu>'!\n\n" +
		"Puja Fesol Ajaib! Ululululululu...", nil)
}

// /tanya function.
// This function to get user question and answer it
func ask(bot *telebot.Bot, message telebot.Message) {
	if strings.ToLower(message.Text) == "/tanya" {
		bot.SendMessage(message.Chat, mentionUser(message) + " : Ketik pertanyaan kamu di belakang '/tanya'!", nil)
	} else {
		bot.SendMessage(message.Chat, mentionUser(message) + " : " + answer(message.Text), nil)
	}
}

// Start the bot
func start(botKEY string) {
	// Initial new bot with token `BotKEY`
	bot, err := telebot.NewBot(botKEY)
	if err != nil {
		log.Fatalln(err)
	}

	messages := make(chan telebot.Message, 100)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		// Capture to log file
		writeLog(time.Now().Format(time.RFC1123) + "  ->  " + mentionUser(message) + "  ->  " + message.Text + "\n")
		println(time.Now().Format(time.RFC1123) + "  ->  " + mentionUser(message) + "  ->  " + message.Text)

		if message.Text == "/hi" {
			hi(bot, message)
		} else if strings.Contains(message.Text, "/help") || message.Text == "/tanya@FesolAjaibBot" {
			help(bot, message)
		} else if strings.Contains(message.Text, "/tanya") {
			ask(bot, message)
		} else if strings.HasPrefix(message.Text, "@FesolAjaibBot") || strings.HasSuffix(message.Text, "@FesolAjaibBot") {
			if strings.ToLower(message.Text) == "@fesolajaibbot" {
				bot.SendMessage(message.Chat, mentionUser(message) + " : Ketik pertanyaan kamu di belakang '/tanya'!", nil)
			} else {
				bot.SendMessage(message.Chat, mentionUser(message) + " : " + answer(message.Text), nil)
			}
		}
	}
}

func main() {
	config := LoadConfig("./config.json")
	start(config.API_KEY)

//	println(ans)
//	println(jawab("Tes pertanyaan"))
}

// Write log message to log file
// First must create log file named 'log.txt'
func writeLog(text string) {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed open log file. ", err)
	}

	defer f.Close()

	if _,err = f.WriteString(text); err != nil {
		log.Fatalln("Failed writing to log file. ", err)
	}
}