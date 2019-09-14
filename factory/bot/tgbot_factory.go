package factory

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/whoismath/ssrf-pwnabot/config"
)

// BotFactory is a telegram bot factory
func BotFactory(c *config.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(c.TelegramBotToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = c.Debug

	return bot, nil
}
