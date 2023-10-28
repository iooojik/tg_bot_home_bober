package main

import (
	"github.com/spf13/viper"
	"home_chief/internal/bot"
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

	b := bot.NewBot(token, nil)
	err = b.Run()
	if err != nil {
		return err
	}

	return nil
}
