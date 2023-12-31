package main

import (
	"github.com/spf13/viper"
	"home_chief/internal/bot"
	"home_chief/internal/service"
	"log"
)

func main() {
	err := runBot()
	if err != nil {
		log.Fatal(err)
	}
}

func runBot() error {
	viper.SetConfigFile(".env.local")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	token := viper.GetString("BOT_TOKEN")
	srv := new(service.BotService)
	b := bot.NewBot(token, srv)
	err = b.Run()
	if err != nil {
		return err
	}

	return nil
}
