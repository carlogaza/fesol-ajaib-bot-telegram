package main

import (
	"log"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/tucnak/telebot"
)

var API_KEY string

type Configuration struct {
	API_KEY		string   `json:"API_KEY"`
}

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

func jumlah(kata string) int {
	var sum int = 0

	byteArray := []byte(kata)
	for i := 0; i < len(byteArray); i++ {
		sum += int(byteArray[i])
	}

	return sum
}

func jawab(tanya string) string {
	var ans string

	words := strings.Split(tanya, " ")

	if strings.Contains(strings.ToLower(tanya), "atau") && len(words) > 2 {
		for i := 0; i < len(words); i++ {
			if strings.ToLower(words[i]) == "atau" {
				pilih := (jumlah(words[i-1]) + jumlah(words[i+1])) % 2
				if pilih == 0 {
					if jumlah(words[i-1]) > jumlah(words[i+1]) {
						ans = words[i-1]
					} else {
						ans = words[i+1]
					}
				} else {
					if jumlah(words[i-1]) < jumlah(words[i+1]) {
						ans = words[i-1]
					} else {
						ans = words[i+1]
					}
				}
				i = len(words) + 9999999
			}
		}
	} else if (strings.Contains(strings.ToLower(tanya), "fesol") && strings.Contains(strings.ToLower(tanya), "jembut")) ||
		(strings.Contains(strings.ToLower(tanya), "fesol") && strings.Contains(strings.ToLower(tanya), "jmbt")) {
		ans = "ABSOLUTELY!"
	} else {
		switch jumlah(tanya) % 8 {
		case 0:
			ans = "Tidak!"
		case 1:
			ans = "Jangan!"
		case 2:
			ans = "Iya!"
		case 3:
			ans = "Mungkin."
		case 4:
			ans = "Bisa jadi."
		case 5:
			ans = "Aku tidak yakin."
		case 6:
			ans = "Aku ragu."
		case 7:
			ans = "Mungkin suatu hari nanti."
		}
	}

	return ans;
}

func start() {
	bot, err := telebot.NewBot(API_KEY)
	if err != nil {
		log.Fatalln(err)
	}

	pesan := make(chan telebot.Message, 100)
	bot.Listen(pesan, 1*time.Second)

	for pesan := range pesan {
		println(pesan.Text)

		if pesan.Text == "/hi" {
			if pesan.Sender.Username == "" {
				bot.SendMessage(pesan.Chat,
					"Hello, " + pesan.Sender.FirstName + " " +
					pesan.Sender.LastName + "!", nil)
			} else {
				bot.SendMessage(pesan.Chat,
					"Hello, @" + pesan.Sender.Username + "!", nil)
			}
		} else if strings.Contains(pesan.Text, "/help") || pesan.Text == "/tanya@FesolAjaibBot" {
			bot.SendMessage(pesan.Chat,
				"Tanyakan apapun pada Fesol Ajaib!\n Ketik '@FesolAjaibBot <pertanyaan kamu>'!\n\n Puja Fesol Ajaib! Ululululululu...", nil)
		} else if strings.Contains(pesan.Text, "/tanya") {
			if strings.ToLower(pesan.Text) == "/tanya" {
				if pesan.Sender.Username == "" {
					bot.SendMessage(pesan.Chat,
						pesan.Sender.FirstName + " " + pesan.Sender.LastName + " : " +
						"Ketik pertanyaan kamu di belakang '@FesolAjaibBot'", nil)
				} else {
					bot.SendMessage(pesan.Chat,
						"@" + pesan.Sender.Username + " : " + "Ketik pertanyaan kamu di belakang '@FesolAjaibBot'", nil)
				}
			} else {
				if pesan.Sender.Username == "" {
					bot.SendMessage(pesan.Chat,
						pesan.Sender.FirstName + " " + pesan.Sender.LastName + " : " +
						jawab(pesan.Text), nil)
				} else {
					bot.SendMessage(pesan.Chat,
						"@" + pesan.Sender.Username + " : " + jawab(pesan.Text), nil)
				}
			}
		} else if strings.HasPrefix(pesan.Text, "@FesolAjaibBot") || strings.HasSuffix(pesan.Text, "@FesolAjaibBot") {
			if strings.ToLower(pesan.Text) == "@fesolajaibbot" {
				if pesan.Sender.Username == "" {
					bot.SendMessage(pesan.Chat,
						pesan.Sender.FirstName + " " + pesan.Sender.LastName + " : " +
						"Ketik pertanyaan kamu di belakang '@FesolAjaibBot'", nil)
				} else {
					bot.SendMessage(pesan.Chat,
						"@" + pesan.Sender.Username + " : " + "Ketik pertanyaan kamu di belakang '@FesolAjaibBot'", nil)
				}
			} else {
				if pesan.Sender.Username == "" {
					bot.SendMessage(pesan.Chat,
						pesan.Sender.FirstName + " " + pesan.Sender.LastName + " : " +
						jawab(pesan.Text), nil)
				} else {
					bot.SendMessage(pesan.Chat,
						"@" + pesan.Sender.Username + " : " + jawab(pesan.Text), nil)
				}
			}
		}
	}
}

func main() {
	config := LoadConfig("./config.json")
	API_KEY = config.API_KEY

	start()

//	println(API_KEY)
//	println(jawab("Tes pertanyaan"))
}