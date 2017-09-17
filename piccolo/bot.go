package piccolo

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/shawnsilva/piccolo/log"
	"github.com/shawnsilva/piccolo/utils"
	"github.com/shawnsilva/piccolo/version"
	"github.com/shawnsilva/piccolo/youtube"
)

type (
	// Bot holds various config and state of the current bot
	Bot struct {
		conf           *utils.Config
		version        *version.Info
		dg             *discordgo.Session
		discordGuildID string

		player *player
		yt     *youtube.Manager
	}
)

func NewBot(c *utils.Config, v *version.Info) *Bot {
	b := &Bot{
		conf:    c,
		version: v,
	}
	return b
}

// Start will start the bot
func (b *Bot) Start() {
	var err error
	b.dg, err = discordgo.New("Bot " + b.conf.BotToken)
	if err != nil {
		return
	}

	// b.playlist = newPlaylist(b.conf.Bot.UsePlaylist, b.conf.Bot.PlaylistPath)
	b.player = newPlayer(b.conf)

	b.yt = &youtube.Manager{
		APIKey:     b.conf.GoogleAPIKey,
		YtDlPath:   b.conf.Bot.YtDlPath,
		YTCacheDir: path.Join(filepath.ToSlash(b.conf.Bot.CacheDir), "/", "ytdl"),
	}

	b.dg.AddHandler(b.ready)
	b.dg.AddHandler(b.messageCreate)
	b.dg.AddHandler(b.voiceStateChange)

	err = b.dg.Open()
	if err != nil {
		return
	}

	guilds, err := b.dg.UserGuilds(1, "", "")
	if err != nil {
		log.Fatal("Failed to determine connected guild ID")
	}
	if len(guilds) != 1 {
		log.Fatal("Too many guilds")
	}
	b.discordGuildID = guilds[0].ID
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
	if len(b.conf.BindToTextChannels) == 0 || utils.StringInSlice(m.ChannelID, b.conf.BindToTextChannels) {
		if strings.HasPrefix(m.Content, b.conf.CommandPrefix) {
			cmdName := strings.Fields(m.Content)[0][len(b.conf.CommandPrefix):]
			foundCommand, found := cmdHandler.get(cmdName)
			if !found {
				log.WithFields(log.Fields{
					"cmd": cmdName,
				}).Error("Failed to find command")
				return
			}
			cmdFunc := *foundCommand
			cmdFunc(b, m)
		}
	}
}

func (b *Bot) voiceStateChange(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	guild, err := b.dg.State.Guild(b.discordGuildID)
	if err != nil {
		log.Error("Failed to determine voice state")
		return
	}
	if len(guild.VoiceStates) <= 1 {
		// pause music, nobody listening
	}
	for _, vs := range guild.VoiceStates {
		if utils.StringInSlice(vs.ChannelID, b.conf.AutoJoinVoiceChannels) {
			if vs.UserID != b.dg.State.User.ID {
				// at least one user, not the bot is in channel
				return
			}
		}
	}
	// pause music, nobody listening
}

func (b *Bot) reply(message string, m *discordgo.MessageCreate) {
	msg, err := b.dg.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		log.WithFields(log.Fields{
			"msg":   msg,
			"error": err,
		}).Error("Failed to send message")
	}
}
