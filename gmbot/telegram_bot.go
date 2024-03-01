package gmbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramClient struct {
	BotToken string           `json:"botToken"`
	ChatId   int64            `json:"chatId"`
	BotApi   *tgbotapi.BotAPI `json:"botApi"`
}

func InitTelegramClient(botToken string, chatId int64) (*TelegramClient, error) {
	botApi, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("tgbotapi.NewBotAPI %v", err)
	}

	return &TelegramClient{BotToken: botToken, ChatId: chatId, BotApi: botApi}, nil
}

func (tg *TelegramClient) SendMarkdownMessage(message string) error {
	msg := tgbotapi.NewMessage(tg.ChatId, message)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := tg.BotApi.Send(msg)
	if err != nil {
		return fmt.Errorf("tg.BotApi.Send %v", err)
	}
	return nil
}
