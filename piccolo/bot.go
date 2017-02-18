package piccolo

import (
	"github.com/bwmarrin/discordgo"

	"github.com/shawnsilva/piccolo/log"
	"github.com/shawnsilva/piccolo/utils"
)

type (
	// Bot holds various config and state of the current bot
	Bot struct {
		Conf *utils.Config

		dg *discordgo.Session
	}
)

// Start will start the bot
func (b *Bot) Start() {
	var err error
	b.dg, err = discordgo.New("Bot " + b.Conf.BotToken)
	if err != nil {
		return
	}

	b.dg.AddHandler(b.ready)
	b.dg.AddHandler(b.messageCreate)

	err = b.dg.Open()
	if err != nil {
		return
	}
}

// Stop will stop the bot
func (b *Bot) Stop() {
	_ = b.dg.UpdateStatus(0, "")
	b.dg.Close()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	_ = s.UpdateStatus(0, "development")
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(b.Conf.BindToTextChannels) == 0 || utils.StringInSlice(m.ChannelID, b.Conf.BindToTextChannels) {
		log.Info("Message: " + m.Content)
	}
}
