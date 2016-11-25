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

	if strings.Contains(strings.ToLower(query), " atau ") && len(words) > 3 {
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

	// This is just for fun bro. :D
	} else if (strings.Contains(strings.ToLower(query), "fesol") &&
			strings.Contains(strings.ToLower(query), "jembut")) ||
			(strings.Contains(strings.ToLower(query), "fesol") &&
			strings.Contains(strings.ToLower(query), "jmbt")) {
		ans = "ABSOLUTELY!"

	// For asking time
	} else if strings.Contains(strings.ToLower(query), "jam") &&
			(strings.Contains(strings.ToLower(query), "berapa") ||
			strings.Contains(strings.ToLower(query), "brp")) &&
			(strings.Contains(strings.ToLower(query), "sekarang") ||
			strings.Contains(strings.ToLower(query), "skrg")) {
		ans = "Sekarang jam " + time.Now().Format(time.Kitchen)

	// For asking day
	} else if strings.Contains(strings.ToLower(query), "hari") &&
			(strings.Contains(strings.ToLower(query), "sekarang") ||
			strings.Contains(strings.ToLower(query), "skrg")) {
		ans = "Sekarang hari " + time.Now().Weekday().String()

	// For asking month
	} else if strings.Contains(strings.ToLower(query), "bulan") &&
			(strings.Contains(strings.ToLower(query), "sekarang") ||
			strings.Contains(strings.ToLower(query), "skrg")) {
		ans = "Sekarang bulan " + time.Now().Month().String()

	// For asking date
	} else if (strings.Contains(strings.ToLower(query), "tanggal") ||
			strings.Contains(strings.ToLower(query), "tgl")) &&
			(strings.Contains(strings.ToLower(query), "berapa") ||
			strings.Contains(strings.ToLower(query), "brp")) &&
			(strings.Contains(strings.ToLower(query), "sekarang") ||
			strings.Contains(strings.ToLower(query), "skrg")) {
		ans = "Sekarang tanggal " + strconv.Itoa(time.Now().Day()) +
			" " + time.Now().Month().String() + " " + strconv.Itoa(time.Now().Year())

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
			(strings.Contains(strings.ToLower(query), "berapa lama") &&
			strings.Contains(strings.ToLower(query), "akan")) {
		switch calculate(query) % 7 {
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

	// For 'why' question
	} else if strings.Contains(strings.ToLower(query), "kenapa") ||
			strings.Contains(strings.ToLower(query), "mengapa") {
		switch calculate(query) % 5 {
		case 0:
			ans = "Ngeyel sih!"
		case 1:
			ans = "Noob sih!"
		case 2:
			ans = "Mungkin sudah takdir."
		case 3:
			ans = "Sudah sepantasnya seperti itu."
		case 4:
			ans = "Memang harus seperti itu."
		}

	// For 'where' question
	} else if strings.Contains(strings.ToLower(query), "dimana") ||
			strings.Contains(strings.ToLower(query), "di mana") ||
			strings.Contains(strings.ToLower(query), "dmn ") {
		switch calculate(query) % 5 {
		case 0:
			ans = "Sepertinya di rumah."
		case 1:
			ans = "Mungkin di kost."
		case 2:
			ans = "Di kampus mungkin."
		case 3:
			ans = "Lagi di jalan."
		case 4:
			ans = "Di tempat kerja."
		}

	} else if strings.Contains(strings.ToLower(query), "siapa") &&
			(strings.Contains(strings.ToLower(query), "namamu") ||
			strings.Contains(strings.ToLower(query), "kamu")) {
		ans = "Aku Fesol kak, si kulit kerang ajaib, yang dapat menjawab " +
			"segala pertanyaan kamu. Kalau kakak mau tau lebih lanjut tentang " +
			"aku, kakak bisa lihat di sini /tentang."

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
			ans = "Sepertinya iya."
		}
	}
	return ans;
}

// Just say hello to user
func hi(bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat, "Hello, " + mentionUser(message) + "!", nil)
}

// /bantuan function
func help(bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		"Tanyakan apapun kepada Fesol Ajaib!\n" +
		"Ketik \"/tanya <pertanyaan kamu>\"!\n\n" +
		"Bot juga sama seperti manusia, butuh tidur. Oleh " +
		"karena itu, jika bot ini sedang tidur atau tidak " +
		"merespon pertanyaan kamu, coba mention bot dengan " +
		"cara mengetik /tanya@FesolAjaibBot dan ketik ulang " +
		"pertanyaan kamu.\n\n" +
		"Puja Fesol Ajaib! Ululululululu...", nil)
}

// /tentang function
func about(bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		"Fesol adalah nama dari sebuah kulit kerang ajaib yang " +
		"mahatahu dan dapat menjawab segala pertanyaan kamu. Kamu " +
		"harus percaya dengan Fesol Ajaib!\n\nSelamat bergabung " +
		"di klub pemuja Fesol Ajaib!\n" +
		"Puja Fesol Ajaib! Ululululululu...", nil)
}

// /pengembang function
func developer(bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		"Bot Fesol Ajaib ini dibuat berdasarkan inspirasi dari " +
		"kulit kerang ajaib pada kartun SpongeBob Squarepants. " +
		"Kode sumber bot ini tersedia secara terbuka di github " +
		"https://github.com/carlogaza/fesol-ajaib-bot-telegram " +
		"untuk anda yang tertarik mengembangkan bot ini.", nil)
}

// /saran function
func recommendation(bot *telebot.Bot, message telebot.Message, logApp string) {
	if strings.ToLower(message.Text) == "/saran" || strings.ToLower(message.Text) == "/saran@fesolajaibbot" {
		bot.SendMessage(message.Chat,
			"Jika kamu mempunyai saran untuk bot ini, silahkan kirimkan " +
			"saran kamu dengan cara mengetik /saran <spasi> saran kamu. " +
			"Saran kamu sangat bermanfaat untuk pengembangan bot ini sel" +
			"anjutnya.\nTerima kasih. :)", nil)
	} else {
		// Add user recommendation to recommendation file.
		writeLog(logApp + "\n", "saran.txt")
		bot.SendMessage(message.Chat, mentionUser(message) + " : Saran kamu telah kami tampung. Terima kasih. :D", nil)
	}
}

// /tanya function.
// This function to get user question and answer it
func ask(bot *telebot.Bot, message telebot.Message, logApp string) string {
	if strings.ToLower(message.Text) == "/tanya" || strings.ToLower(message.Text) == "/tanya@fesolajaibbot" {
		bot.SendMessage(message.Chat, mentionUser(message) +
		" : Ketik pertanyaan kamu di belakang '/tanya'!", nil)
	} else {
		bot.SendMessage(message.Chat, mentionUser(message) + " : " + answer(message.Text), nil)
		// Add answer to log file
		logApp += "  <->  " + answer(message.Text)
	}
	return logApp
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
		logApp := time.Now().Format(time.RFC1123) + "  ->  " + mentionUser(message) + "  ->  " + message.Text

		if message.Text == "/hi" {
			hi(bot, message)
		} else if strings.Contains(message.Text, "/help") {
			help(bot, message)
		} else if strings.Contains(message.Text, "/pengembang") {
			developer(bot, message)
		} else if strings.Contains(message.Text, "/tentang") {
			about(bot, message)
		} else if strings.Contains(message.Text, "/saran") {
			recommendation(bot, message, logApp)
		} else if strings.Contains(message.Text, "/tanya") {
			logApp = ask(bot, message, logApp)
		} else if strings.HasPrefix(message.Text, "@FesolAjaibBot") ||
					strings.HasSuffix(message.Text, "@FesolAjaibBot") {
			if strings.ToLower(message.Text) == "@fesolajaibbot" {
				bot.SendMessage(message.Chat, mentionUser(message) +
				" : Ketik pertanyaan kamu di belakang '/tanya'!", nil)
			} else {
				bot.SendMessage(message.Chat, mentionUser(message) + " : " + answer(message.Text), nil)
				// Add answer to log file
				logApp += "  <->  " + answer(message.Text)
			}
		}

		// Write to log file
		writeLog(logApp + "\n", "log.txt")
		println(logApp)
	}
}

func main() {
//	config := LoadConfig("./config.json")
//	start(config.API_KEY)
	start("258624514:AAGNsbjp3z2udCV6Aq7t2FXs0FwQ43IFLQM")

//	println(answer("/tanya hari apa sekarang?"))
//	println(answer("/tanya bulan apa sekarang?"))
//	println(answer("/tanya tanggal berapa sekarang?"))
}

// Write to file
// First must create log file named 'log.txt' or 'saran.txt' or whatever
func writeLog(text string, fileName string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed open log file. ", err)
	}

	defer f.Close()

	if _,err = f.WriteString(text); err != nil {
		log.Fatalln("Failed writing to log file. ", err)
	}
}