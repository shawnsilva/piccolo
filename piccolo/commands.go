package piccolo

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/shawnsilva/piccolo/log"
	"github.com/shawnsilva/piccolo/utils"
)

type (
	command    func(b *Bot, m *discordgo.MessageCreate)
	commandMap map[string]command

	commandHandler struct {
		commands commandMap
	}
)

var (
	cmdHandler *commandHandler
)

func init() {
	cmdHandler = &commandHandler{make(commandMap)}
	cmdHandler.addCommand("help", help)
	cmdHandler.addCommand("version", botVersion)
	cmdHandler.addCommand("play", play)
	cmdHandler.addCommand("savePlaylist", savePlaylist)
	cmdHandler.addCommand("showPlaylist", printPlaylist)
}

func (h commandHandler) addCommand(name string, c command) {
	h.commands[name] = c
}

func (h commandHandler) getAllCommands() commandMap {
	return h.commands
}

func (h commandHandler) get(name string) (*command, bool) {
	cmd, found := h.commands[name]
	return &cmd, found
}

func help(b *Bot, m *discordgo.MessageCreate) {
	var msg string
	var cmdListStr string
	var cmdList []string
	for cmdN := range cmdHandler.getAllCommands() {
		cmdList = append(cmdList, b.conf.CommandPrefix+cmdN)
	}
	sort.Strings(cmdList)
	cmdListStr = fmt.Sprintf("```%s```", utils.StrJoin(cmdList, " "))
	msg = fmt.Sprintf("<@%s>, **Commands**\n%s\n%s", m.Author.ID, cmdListStr, "https://github.com/shawnsilva/piccolo/wiki/Commands")
	b.reply(msg, m)
}

func botVersion(b *Bot, m *discordgo.MessageCreate) {
	msg := fmt.Sprintf("<@%s>, `VERSION: %s`", m.Author.ID, b.version.GetVersionString())
	b.reply(msg, m)
}

func play(b *Bot, m *discordgo.MessageCreate) {
	if _, ok := b.textChannelLookup[m.ChannelID]; !ok {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
		}).Error("Failed to find controller from channel id")
		return
	}
	splitString := strings.SplitN(m.Content, " ", 2)
	if len(splitString) <= 1 {
		log.WithFields(log.Fields{
			"command": m.Content,
		}).Warn("!play doesn't appear to contain a search string")
		b.reply(fmt.Sprintf("<@%s> - Sorry, your command didn't appear to have a song to search for: **%s**", m.Author.ID, m.Content), m)
		return
	}
	song := splitString[1]
	result, err := b.yt.SearchFirstResult(song)
	if err != nil {
		log.WithFields(log.Fields{
			"song":  song,
			"error": err,
		}).Debug("Failed to find song")
		b.reply(fmt.Sprintf("<@%s> - Sorry, couldn't find a result for: **%s**", m.Author.ID, song), m)
		return
	}
	b.textChannelLookup[m.ChannelID].player.playlist.addSong(m.Author, m.ChannelID, result.ID.VideoID, result.Snippet.Title)
	go b.textChannelLookup[m.ChannelID].player.downloadNextSong()
	b.reply(fmt.Sprintf("<@%s> - Enqueued **%s** to be played.", m.Author.ID, result.Snippet.Title), m)
}

func savePlaylist(b *Bot, m *discordgo.MessageCreate) {
	if _, ok := b.textChannelLookup[m.ChannelID]; !ok {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
		}).Error("Failed to find controller from channel id")
		return
	}
	err := b.textChannelLookup[m.ChannelID].player.playlist.savePlaylist()
	if err == nil {
		b.reply(fmt.Sprintf("<@%s> - Saved current playlist to disk.", m.Author.ID), m)
	} else {
		b.reply(fmt.Sprintf("<@%s> - %s.", m.Author.ID, err.Error()), m)
	}
}

func printPlaylist(b *Bot, m *discordgo.MessageCreate) {
	if _, ok := b.textChannelLookup[m.ChannelID]; !ok {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
		}).Error("Failed to find controller from channel id")
		return
	}
	b.reply(fmt.Sprintf("<@%s> - **Current Playlist**\n\n%s", m.Author.ID, b.textChannelLookup[m.ChannelID].player.playlist), m)
}
