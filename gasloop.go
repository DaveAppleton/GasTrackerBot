package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

func priceMessage(data GasData, advert string) (message string) {

	return fmt.Sprintf("*Gas prices*\n\n```"+`
Tx Speed Gas needed  expected wait
--------+-----------+-------------
Fastest  %5v GWei   %v
Fast     %5v GWei   %v
Safelow  %5v GWei   %v
Average  %5v GWei   %v
`+"```"+`
%s
`, data.Fastest, data.FastestWait, data.Fast, data.FastWait, data.SafeLow, data.SafeLowWait, data.Average, data.AverageWait, advert)

}

func gasLoop() {
	log.Println("gasloop started")
	for {

		gasHTML, _, data, err := multiGas()
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}
		err = data.save()
		if err != nil {
			log.Println("error saving data", err)
		}

		// users, err := getAllBelowOrEqual(median)
		// if err != nil {
		// 	log.Println("error getting Users", err)
		// } else {
		// 	for _, user := range users {
		// 		msg := tgbotapi.NewMessage(user, fmt.Sprint("Safe Low gas price is now ", data.SafeLow, " GWei"))
		// 		msg.ParseMode = "MarkdownV2"
		// 		msg.DisableWebPagePreview = true
		// 		_, err := bot.Send(msg)
		// 		if err != nil {
		// 			log.Println("could not send message to ", user, err)
		// 		}
		// 	}
		// }
		users, err := getAllWhoWantInfo()
		if err != nil {
			log.Println("error getting Users", err)
		} else {
			pass := 0
			fail := 0
			for _, user := range users {
				msg := tgbotapi.NewMessage(user, gasHTML)
				msg.ParseMode = "MarkdownV2"
				msg.DisableWebPagePreview = true
				_, err := bot.Send(msg)
				if err != nil {
					fail++
					log.Println("could not send gas info message to ", user, err)
					//log.Println(gasHTML)
				} else {
					pass++
				}
			}
			log.Printf("%d users want info, %d pass, %d fail\n", len(users), pass, fail)
		}
		version := viper.GetFloat64("Version")
		list, err := getAllUsersBelowVersion(version)
		if err != nil {
			log.Println("All usersBelowVersion", err)
		} else {

			WhatsNew, err := newInVersion(version)
			if err == nil {
				for _, user := range list {
					msg := tgbotapi.NewMessage(user, WhatsNew)
					msg.ParseMode = "MarkdownV2"
					msg.DisableWebPagePreview = true
					_, err := bot.Send(msg)
					if err != nil {
						log.Println("could not send Whatsnew to ", user, err)
					} else {
						err = updateUserVersion(user, version)
						if err != nil {
							log.Println("update user version", user, version, err)
						}
					}
				}
			}
		}
		time.Sleep(time.Hour)
	}
}
