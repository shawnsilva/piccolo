package piccolo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"

	"github.com/shawnsilva/piccolo/utils"
)

type (
	player struct {
		conf           *utils.Config
		playlist       *playlist
		guildID        string
		voiceChannelID string
		vc             *discordgo.VoiceConnection
		stream         *dca.StreamingSession
	}
)

func newPlayer(confpointer *utils.Config, guildID string, voiceChID string) *player {
	p := &player{conf: confpointer, guildID: guildID, voiceChannelID: voiceChID}
	p.playlist = newPlaylist(p.conf.Bot.UsePlaylist, p.conf.Bot.PlaylistPath)
	return p
}

func (p *player) Shutdown() error {
	if p.stream != nil {
		p.stream.SetPaused(true)
	}
	if p.vc != nil {
		err := p.vc.Disconnect()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *player) JoinVoiceChannel(dg *discordgo.Session) error {
	vc, err := dg.ChannelVoiceJoin(p.guildID, p.voiceChannelID, false, true)
	if err != nil {
		return err
	}
	p.vc = vc
	return nil
}
