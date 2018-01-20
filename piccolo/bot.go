package piccolo

import (
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"

	"github.com/shawnsilva/piccolo/log"
	"github.com/shawnsilva/piccolo/utils"
	"github.com/shawnsilva/piccolo/version"
	"github.com/shawnsilva/piccolo/youtube"
)

type (
	// Bot holds various config and state of the current bot
	Bot struct {
		conf    *utils.Config
		version *version.Info
		dg      *discordgo.Session

		guildLookup        map[string]*guildControls
		textChannelLookup  map[string]*guildControls
		voiceChannelLookup map[string]*guildControls

		lock *sync.Mutex

		yt *youtube.Manager
	}

	guildControls struct {
		guildID        string
		textChannelIDs []string
		voiceChannelID string
		player         *player
	}
)

// NewBot will create an instance of a bot
func NewBot(c *utils.Config, v *version.Info) *Bot {
	b := &Bot{
		conf:    c,
		version: v,
		lock:    &sync.Mutex{},
	}
	return b
}

// Start will start the bot
func (b *Bot) Start() {
	b.lock.Lock()
	defer b.lock.Unlock()
	var err error
	b.dg, err = discordgo.New("Bot " + b.conf.BotToken)
	if err != nil {
		return
	}

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

	for _, guild := range b.conf.Guilds {
		vch, err := b.dg.Channel(guild.AutoJoinVoiceChannel)
		if err != nil {
			log.WithFields(log.Fields{
				"voicechannel": guild.AutoJoinVoiceChannel,
				"error":        err,
			}).Error("Failed to find channel information")
			continue
		}
		gID := vch.GuildID
		if _, ok := b.guildLookup[gID]; ok {
			log.WithFields(log.Fields{
				"guild":        gID,
				"voicechannel": guild.AutoJoinVoiceChannel,
			}).Error("This guild already has a configured voice channel")
			continue
		}
		guildInfo, err := b.dg.Guild(gID)
		if err != nil {
			log.WithFields(log.Fields{
				"guild":        gID,
				"voicechannel": guild.AutoJoinVoiceChannel,
				"error":        err,
			}).Error("Failed to find guild information")
			continue
		}
		var textChIDs []string
		for _, tChannelID := range guild.BindToTextChannels {
			tChInfo, err := b.dg.Channel(tChannelID)
			if err != nil {
				log.WithFields(log.Fields{
					"textchannel": tChannelID,
					"error":       err,
				}).Error("Failed to find text channel information")
				break
			}
			if tChInfo.GuildID != gID {
				log.WithFields(log.Fields{
					"textchannel":    tChannelID,
					"voice guild id": gID,
					"found guild id": tChInfo.GuildID,
				}).Error("Text channel isn't in same guild as voice channel")
				break
			}
			textChIDs = append(textChIDs, tChannelID)
		}

		if len(textChIDs) != len(guild.BindToTextChannels) {
			log.WithFields(log.Fields{
				"WantedTextChannels":  guild.BindToTextChannels,
				"valid text channels": textChIDs,
			}).Error("Couldn't verify text channels in same guild as voice")
			break
		}

		gControl := &guildControls{
			guildID:        gID,
			voiceChannelID: guild.AutoJoinVoiceChannel,
			textChannelIDs: textChIDs,
			player:         newPlayer(b.conf, gID, guild.AutoJoinVoiceChannel),
		}
		b.guildLookup = make(map[string]*guildControls)
		b.textChannelLookup = make(map[string]*guildControls)
		b.voiceChannelLookup = make(map[string]*guildControls)
		b.guildLookup[gID] = gControl
		if len(textChIDs) >= 1 {
			for _, tChID := range textChIDs {
				b.textChannelLookup[tChID] = gControl
			}
		} else {
			for _, channel := range guildInfo.Channels {
				if channel.Type == discordgo.ChannelTypeGuildText {
					textChIDs = append(textChIDs, channel.ID)
					b.textChannelLookup[channel.ID] = gControl
				}
			}
			gControl.textChannelIDs = textChIDs
		}
		b.voiceChannelLookup[guild.AutoJoinVoiceChannel] = gControl
	}
}

// Stop will stop the bot
func (b *Bot) Stop() {
	_ = b.dg.UpdateStatus(0, "")
	for _, voiceCh := range b.voiceChannelLookup {
		err := voiceCh.player.Shutdown()
		if err != nil {
			log.Error(err)
		}
	}
	b.dg.Close()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	_ = s.UpdateStatus(0, "development")
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, vChannel := range b.voiceChannelLookup {
		err := vChannel.player.JoinVoiceChannel(s)
		if err != nil {
			log.Error(err)
		}
	}
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, ok := b.textChannelLookup[m.ChannelID]; ok {
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
	guild, err := b.dg.State.Guild(v.GuildID)
	if err != nil {
		log.Error("Failed to determine voice state")
		return
	}
	if len(guild.VoiceStates) <= 1 {
		// pause music, nobody listening
		log.Debug("Pausing Music, nobody in voice channel")
		return
	}
	for _, vs := range guild.VoiceStates {
		if _, ok := b.voiceChannelLookup[vs.ChannelID]; ok {
			if vs.UserID != b.dg.State.User.ID {
				// at least one user, not the bot is in channel
				log.Debug("Playing Music")
				return
			}
		}
	}
	// pause music, nobody listening
	log.Debug("Pausing Music, nobody in voice channel")
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
