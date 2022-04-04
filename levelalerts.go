package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func processAlerts() {
	_, median, _, err := multiGas()
	users, err := getAllBelowOrEqual(median)
	if err != nil {
		log.Println("error getting Users", err)
	} else {
		for _, user := range users {
			msg := tgbotapi.NewMessage(user, fmt.Sprint("Safe Low gas price is now ", median, " GWei"))
			msg.ParseMode = "MarkdownV2"
			msg.DisableWebPagePreview = true
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("could not send level alert message to ", user, err)
			} else {
				updateLastPinged(user)
			}
		}
	}
}

func levelAlertLoop(exit chan bool) {
	for {
		select {
		case <-time.After(10 * time.Minute):
			go processAlerts()
		case <-exit:
			return
		}
	}
}
