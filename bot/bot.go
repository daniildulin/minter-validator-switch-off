package bot

import (
	"context"
	"github.com/daniildulin/minter-validator-switch-off/locale"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
)

func (b *TgBot) SendMsg(msg string) {
	_, err := b.tgApi.Send(b.managerChat, msg)
	if err != nil {
		b.log.Error(err)
	}
}

func (b *TgBot) Run() {
	b.tgApi.Start()
}

func New() *TgBot {
	//Init Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)

	b, err := tb.NewBot(tb.Settings{
		Token: os.Getenv("TG_TOKEN"),
	})

	if err != nil {
		logger.Panic(err)
	}

	return &TgBot{
		ctx:    context.Background(),
		tgApi:  b,
		locale: locale.NewBotLocale(),
		log: logger.WithFields(logrus.Fields{
			"version": "v1.0",
			"app":     "Minter Validator Protector Telegram Bot",
		}),
		managerChat: Chat{
			Id: os.Getenv("TG_CHANNEL_ID"),
		},
	}
}

type TgBot struct {
	ctx         context.Context
	tgApi       *tb.Bot
	log         *logrus.Entry
	locale      *locale.Locale
	managerChat Chat
}

type Chat struct {
	Id string `json:"id"`
}

func (c Chat) Recipient() string {
	return c.Id
}
