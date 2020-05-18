package main

import (
	"log"
	"os"
	"os/exec"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	bot, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "type /info or /status."
			case "gitpull":
				msg.Text = execGitPull()
			case "status":
				msg.Text = updateAllScripts()
			case "info":
				msg.Text = "I'm VPSService Admin Bot."
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}

}

func execGitPull() string {
	cmd := exec.Command("/opt/projects/vpsserver-go-bot/scripts/", "./git_pull.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err.Error()
	}
	return "ok"
}

func updateAllScripts() string {
	if err := os.Chmod("/opt/projects/vpsserver-go-bot/scripts/git_pull.sh", 0777); err != nil {
		return err.Error()
	}
	return "ok"
}
