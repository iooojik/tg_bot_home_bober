package bot

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"strconv"
	"strings"
)

const (
	cmdStart = "/start"
)

type Service interface {
	CheckUser(userId int) (string, error)
	ChangeDate(date, userId int) (string, error)
}

type Bot struct {
	token string
	debug bool
	// chat id - command
	chatContext map[int64]string
	srv         Service
}

func NewBot(token string, srv Service) *Bot {
	b := &Bot{
		token:       token,
		srv:         srv,
		chatContext: map[int64]string{},
		debug:       viper.GetBool("DEBUG"),
	}
	return b
}

func (b *Bot) Run() error {
	bot, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return err
	}
	bot.Debug = b.debug
	slog.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	for update := range updates {
		message := update.Message
		if update.Message == nil {
			continue
		}
		log.Printf("[%v][%v] %s", message.From.ID, message.From.UserName, message.Text)
		msg, err := b.handleCommand(message)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		_, err = bot.Send(msg)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	var (
		msg string
		err error
	)
	if !message.IsCommand() {
		chatCtx, exists := b.chatContext[message.Chat.ID]
		if !exists {
			message.Text = cmdStart
		}
		msg, err = b.handleChatCtx(chatCtx, message)
		if err != nil {
			return nil, err
		}
	}
	switch strings.ToLower(strings.TrimSpace(message.Text)) {
	case cmdStart:
		_, exists := b.chatContext[message.Chat.ID]
		if exists {
			delete(b.chatContext, message.Chat.ID)
		}
		msg, err = b.srv.CheckUser(message.From.ID)
		if msg == "" {
			//	пользователь впервые воспользовался ботом
			msg = fmt.Sprint("Напишите число месяца, когда будет приходить напоминание")
			b.chatContext[message.Chat.ID] = cmdStart
		}
	}
	if err != nil {
		return nil, err
	}
	msgConf := tgbotapi.NewMessage(message.Chat.ID, msg)
	return &msgConf, nil
}

func (b *Bot) handleChatCtx(chatCtx string, message *tgbotapi.Message) (string, error) {
	switch chatCtx {
	case cmdStart:
		text := strings.TrimSpace(message.Text)
		val, err := strconv.Atoi(text)
		if err != nil {
			b.chatContext[message.Chat.ID] = cmdStart
			return "Введите число. Например, 10", nil
		}
		_, err = b.srv.ChangeDate(val, message.From.ID)
		if err != nil {
			return "", err
		}
		delete(b.chatContext, message.Chat.ID)
		return "Дата успешно добавлена. Чтобы обновить дату, напишите /start", nil
	}
	return "impossible!!!", nil
}
