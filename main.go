package main

import (
	bot "discord_bot/bot"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	bot.BotToken = os.Getenv("BOT_TOKEN")
	bot.Run()
}
