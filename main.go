package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

//const endpoint = "https://ethgasstation.info/api/ethgasAPI.json"

// GasData holds data returned from ETHGasStation
type GasData struct {
	Fast        decimal.Decimal
	Fastest     decimal.Decimal
	SafeLow     decimal.Decimal
	Average     decimal.Decimal
	FastWait    decimal.Decimal
	FastestWait decimal.Decimal
	SafeLowWait decimal.Decimal
	AverageWait decimal.Decimal
	BlockTime   decimal.Decimal `json:"block_time"`
	BlockNum    uint64
	DateAdded   time.Time
}

func (myData GasData) save() error {
	return saveGasData(&myData)
}

func (myData *GasData) load() (err error) {
	myData, err = loadLatestGasData()
	return
}

func getGasData() (data GasData, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", viper.GetString("GAS_STATION"), nil)
	if err != nil {
		log.Println("GetGasData, NewRequest:", err)
		return
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Println("GetGasData, Do:", err)
		return
	}
	dat, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("GetGasData, ReadAll:", err)
		return
	}
	//pretty.Println(string(dat))
	err = json.Unmarshal(dat, &data)
	if err != nil {
		log.Println("GetGasData, Unmarshall:", err)
		log.Println(string(dat))
	}
	return
}

func onOff(val bool) string {
	if val {
		return "ON"
	}
	return "OFF"
}

func messages(header string) string {
	msg := viper.GetString(header)
	if len(msg) != 0 {
		return msg + "\n\n" + viper.GetString("valid_commands")
	}
	log.Println("unknown message", header)
	switch header {
	case "start_error":
		return "sorry - an error has occurred.\nPlease try later"
	case "welcome_msg":
		return "Welcome to the GasBot\ncommand(s):\nusage :" + messages("level_usage")
	case "level_error":
		return "usage : " + messages("level_usage")
	case "level_usage":
		return "/level <gas level for alert>"
	default:
		log.Println("unknown message", header)
		return "unknown message : " + header
	}
}

func newMessage(chatID int64, tag string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, messages(tag))
	parseMode := viper.GetString("parsemode")
	if len(parseMode) != 0 {
		msg.ParseMode = parseMode
	}
	return msg
}

func messageLoop(gStream chan gasNowData) {
	// message loop waits for replies from spamees
	// only thing handles are
	//    /ack_nnnn which sets the ack flag
	//    /nak_nnnn which indicates that person cannot handle it
	var gnd gasNowData
	var gndReceived time.Time
	chi := getTokenPrice("chi-gastoken")
	gst2 := getTokenPrice("gastoken")
	gasTokenUpdate := time.Now()

	log.Println("Messageloop started")
	var msg tgbotapi.MessageConfig

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		fmt.Println("👀")
		repeat := true
		for repeat {
			select {
			case gnd = <-gStream:
				fmt.Println("received gasNow data")
				gndReceived = time.Now()
			default:
				fmt.Println("break")
				repeat = false
			}
		}
		if update.Message == nil {
			log.Println("weird update", update)
			continue
		}
		fmt.Println(update.Message.Text)
		user := update.Message.From
		msgStrings := strings.Split(update.Message.Text, " ")
		if len(msgStrings) == 0 {
			continue
		}
		command := msgStrings[0]
		chatID := update.Message.Chat.ID
		msgID := update.Message.MessageID
		oldUser, err := saveUserToDatabase(user.FirstName, user.LastName, user.UserName, user.ID)
		if err != nil {
			log.Println(user.FirstName, user.LastName, user.UserName, user.ID, err)
			msg = newMessage(chatID, "start_error")
		} else if !oldUser {
			message, err := welcomeMessage()
			if err != nil {
				msg = newMessage(chatID, "welcome_msg")
				msg.Text += viper.GetString("ADVERT")
			} else {
				msg = tgbotapi.NewMessage(chatID, message)
				msg.ParseMode = "markdown"
			}
		}

		switch {
		case command == "/start":
			if msg.Text == "" {
				message, err := welcomeMessage()
				if err != nil {
					msg = newMessage(chatID, "welcome_msg")
					msg.Text += viper.GetString("ADVERT")
				} else {
					msg = tgbotapi.NewMessage(chatID, message)
					msg.ParseMode = "markdown"
				}
			}
		case command == "/alert":
			var message string
			if len(msgStrings) != 2 {
				log.Println(user.FirstName, user.LastName, user.UserName, user.ID, update.Message.Text)
				message = readMD("level.md")
			} else {
				data := strings.ToLower(msgStrings[1])
				minutes, err := strconv.Atoi(data)
				if err != nil {
					log.Println("level number ", minutes)
					message = levelError(data)
				} else if err := saveAlertToDatabase(user.ID, minutes); err != nil {
					log.Println(user.FirstName, user.LastName, user.UserName, user.ID, err)
					message = readMD("alert_save_error.md")
				} else {
					message, err = alertResponse(minutes, gnd.Safe)
					if err != nil {
						message = "you will be alerted every " + data + " minutes once the gas falls below your target level"
					}
				}
			}
			msg = tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "markdown"
		case command == "/level":
			var message string
			if len(msgStrings) != 2 {
				log.Println(user.FirstName, user.LastName, user.UserName, user.ID, update.Message.Text)
				message = readMD("level.md")
			} else {
				data := strings.ToLower(msgStrings[1])
				if data == "off" || data == "no" {
					message = readMD("level_off.md")
					if err = turnOffLevelAlerts(int64(user.ID)); err != nil {
						log.Println("level off ", data, user.ID, err)
						message = readMD("level_off_error.md")
					} else {
						message = readMD("level_off.md")
					}
				} else {
					amount, err := decimal.NewFromString(data)
					if err != nil {
						log.Println("level number ", data)
						message = levelError(data)
					} else if err := saveLevelToDatabase(user.ID, amount); err != nil {
						log.Println(user.FirstName, user.LastName, user.UserName, user.ID, err)
						message = readMD("level_save_error.md")
					} else {
						message, err = levelResponse(amount, gnd.Safe)
						if err != nil {
							message = "you will be alerted once the gas falls below *" + data + " GWei*"
						}
					}
				}
			}
			msg = tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "markdown"
		case command == "/info":
			if len(msgStrings) != 2 {
				log.Println(user.FirstName, user.LastName, user.UserName, user.ID, update.Message.Text)
				message := readMD("info.md")
				msg = tgbotapi.NewMessage(chatID, message)
				msg.ParseMode = "markdown"
			} else {
				data := strings.ToLower(msgStrings[1])
				if data == "on" || data == "yes" {
					if err := turnOnInfo(int64(user.ID)); err != nil {
						log.Println("turn on info for ", user.ID, err)
						msg = newMessage(chatID, "tech_error")
					} else {
						msg = newMessage(chatID, "info_on")
					}
				} else if data == "off" || data == "no" {
					if err := turnOffInfo(int64(user.ID)); err != nil {
						log.Println("turn off info for ", user.ID, err)
						msg = newMessage(chatID, "tech_error")
					} else {
						msg = newMessage(chatID, "info_off")
					}
				} else {
					log.Println("parameter err : info for ", user.ID, update.Message.Text)
					msg = newMessage(chatID, "info_parameter_error")
				}
			}
		case command == "/help":
			message := readMD("info.md")
			msg = tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "markdown"
		case command == "/status":
			levelState, infoState, levelValue, err := getUserInfo(user.ID)
			if err != nil {
				msg = newMessage(chatID, "tech_error")
			}
			message := "conditional gas level notifications are *%v*\n"
			if levelState {
				message += "alerts when safe low gas is below *%v*\n"
			}
			message += "hourly complete gas info is *%v*"
			if levelState {
				message = fmt.Sprintf(message, onOff(levelState), levelValue, onOff(infoState))
			} else {
				message = fmt.Sprintf(message, onOff(levelState), onOff(infoState))
			}
			message = message + viper.GetString("ADVERT")
			msg = tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "markdown"
		case command == "/week":
			gds, err := loadWeeklyData()
			if err != nil {
				log.Println("weekly chart ", err)
			}
			fname := fmt.Sprintf("week%d.png", msgID)
			w, _ := os.Create(fname)
			defer w.Close()
			m := buildMap(&gds)

			if m != nil {
				png.Encode(w, m) //Encode writes the Image m to w in PNG format.
				w.Close()
				msg := tgbotapi.NewPhotoUpload(chatID, fname)
				msg.ReplyToMessageID = msgID
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "error")
				msg.ReplyToMessageID = msgID
				bot.Send(msg)
			}
			os.Remove(fname)
			continue
		case command == "/day":
			gds, err := loadDailyData()
			if err != nil {
				log.Println("daily chart ", err)
			}
			fname := fmt.Sprintf("day%d.png", msgID)
			w, _ := os.Create(fname)
			defer w.Close()
			m := buildMap(&gds)
			if m != nil {
				png.Encode(w, m) //Encode writes the Image m to w in PNG format.
				w.Close()
				msg := tgbotapi.NewPhotoUpload(chatID, fname)
				msg.ReplyToMessageID = msgID
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "error")
				msg.ReplyToMessageID = msgID
				bot.Send(msg)
			}
			os.Remove(fname)
			continue
		case command == "/gasnow":
			fmt.Println("gasnow")
			message := gasNowMessage(gnd, gndReceived, viper.GetString("ADVERT"))
			msg = tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "markdown"
		case command == "/gastoken":
			if time.Since(gasTokenUpdate).Minutes() > 5 {
				chi = getTokenPrice("CHI")
				gst2 = getTokenPrice("GST2")
				gasTokenUpdate = time.Now()
			}
			message := ""
			if chi.OK {
				message += "\n" + chi.Message
			}
			if gst2.OK {
				message += "\n" + gst2.Message
			}
			message = strings.Replace(message, ".", "\\.", 5)
			message = strings.Replace(message, "-", "\\-", 5)
			if message == "" {
				message = "No GasToken information found"
			} else {
				message += "\nprices courtesy of [Coingecko](https://coingecko.com)\n\n"
			}
			message += readMD("advert.md")
			fmt.Println(message)
			msg = tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "MarkdownV2"
		case command == "/sources":
			message := readMD("sources.md")
			msg = tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "MarkdownV2"
		case command == "/whatsnew":
			message := "sorry cannot process that right now..."
			version := viper.GetFloat64("Version")
			update, err := newInVersion(version)
			if err == nil {
				message = update
			} else {
				log.Println(err)
			}
			msg = tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "MarkdownV2"
		default:
			log.Println("unknown command", user.FirstName, user.LastName, user.UserName, user.ID, command)
			if len(msg.Text) == 0 {
				msg = newMessage(chatID, "unrecognised_command")
			}
		}
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = msgID
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func rowExists(db *sql.DB, subquery string, args ...interface{}) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT exists (%s)", subquery)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func initViper() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("viper changed", e.Name)
	})
}

func main() {
	gndChan := make(chan gasNowData, 100)
	var BotToken string
	initViper()
	logName := viper.GetString("log")
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/" + logName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
	data, err := getGasData()
	if err != nil {
		log.Fatal("Gas Data not working")
	}
	fmt.Println("Fast", data.Fast.String())
	fmt.Println("Fastest", data.Fastest.String())
	fmt.Println("SafeLow", data.SafeLow.String())
	fmt.Println("Average", data.Average.String())
	fmt.Println("FastWait", data.FastWait.String())
	fmt.Println("FastestWait", data.FastestWait.String())
	fmt.Println("SafeLowWait", data.SafeLowWait.String())
	fmt.Println("AverageWait", data.AverageWait.String())
	fmt.Println("BlockTime", data.BlockTime.String())
	fmt.Println("BlockNum", data.BlockNum)
	if viper.GetBool("live") {
		BotToken = viper.GetString("ALERTBOT_TELE_TOKEN")
		log.Println("Using live BOT")
	} else {
		BotToken = viper.GetString("ALERTBOT_STAGING_TOKEN")
		log.Println("Using staging BOT")
	}
	if len(BotToken) == 0 {
		log.Fatal("Token must be specified", BotToken)
	}

	//authToken = viper.GetString("ALERTBOT_AUTH_TOKEN")
	bot, err = tgbotapi.NewBotAPI(BotToken) //"240960422:AAFYlmg7nIfVKebw3NWqo2oMUbhURtddg2g")
	bot.Debug = false
	exit := make(chan bool, 5)
	go gasNowLoop(gndChan, exit)
	go messageLoop(gndChan)
	go gasLoop()
	go levelAlertLoop(exit)
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	defer signal.Stop(stop)
	_ = <-stop
	exit <- true
	log.Println("Stop sequence initiated...")

}
