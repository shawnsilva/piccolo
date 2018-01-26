package piccolo

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"

	"github.com/shawnsilva/piccolo/log"
	"github.com/shawnsilva/piccolo/utils"
	"github.com/shawnsilva/piccolo/youtube"
)

type (
	player struct {
		conf           *utils.Config
		playlist       *playlist
		guildID        string
		voiceChannelID string
		vc             *discordgo.VoiceConnection
		stream         *dca.StreamingSession
		yt             *youtube.Manager

		streamDoneChan chan error

		currentSong *songAndPath

		dg *discordgo.Session

		downloadLock *sync.Mutex
	}

	songAndPath struct {
		fsPath string
		*PlaylistEntry
	}
)

var errShutdown = errors.New("SHUTDOWN")

func newPlayer(confpointer *utils.Config, guildID string, voiceChID string, youtube *youtube.Manager, downloadLock *sync.Mutex, discordSession *discordgo.Session) *player {
	p := &player{conf: confpointer, guildID: guildID, voiceChannelID: voiceChID, yt: youtube, dg: discordSession}
	p.playlist = newPlaylist(p.conf.Bot.UsePlaylist, p.conf.Bot.PlaylistPath)
	p.downloadLock = downloadLock
	p.streamDoneChan = make(chan error)
	return p
}

func (p *player) Shutdown() error {
	if p.stream != nil {
		p.stream.SetPaused(true)
	}
	p.downloadLock.Lock()
	p.streamDoneChan <- errShutdown
	p.stream = nil
	if p.vc != nil {
		p.vc.Speaking(false)
		err := p.vc.Disconnect()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *player) JoinVoiceChannel() error {
	vc, err := p.dg.ChannelVoiceJoin(p.guildID, p.voiceChannelID, false, true)
	if err != nil {
		return err
	}
	p.vc = vc
	p.downloadNextSong()

	go p.playLoop()
	return nil
}

func (p *player) playLoop() {
	for {
		nextSong, err := p.getNextSongPath()
		// Download the next song in the background
		go p.downloadNextSong()
		if err == nil {
			reader, err := os.Open(filepath.FromSlash(nextSong.fsPath))
			if err != nil {
				continue
			}
			p.currentSong = nextSong
			decoder := dca.NewDecoder(reader)
			p.vc.Speaking(true)

			p.stream = dca.NewStream(decoder, p.vc, p.streamDoneChan)
			p.updateStatus()
			if nextSong.Requester != nil && nextSong.RequestChannelID != "" {
				// Message requester their song is playing
				message := fmt.Sprintf("<@%s> - Your song is now playing: **%s**", nextSong.Requester.ID, nextSong.Title)
				msg, err := p.dg.ChannelMessageSend(nextSong.RequestChannelID, message)
				if err != nil {
					log.WithFields(log.Fields{
						"msg":   msg,
						"error": err,
					}).Error("Failed to send message about request now playing")
				}
			}
			streamErr := <-p.streamDoneChan
			if streamErr == errShutdown {
				return
			}
			if streamErr != nil && streamErr != io.EOF {
				// Handle the error
				log.WithFields(log.Fields{
					"error": streamErr,
				}).Error("Error streaming song")
			}
			p.vc.Speaking(false)
		}
		// Check if the next song is downloaded, if not block until it is. Catches
		// additions to the request queue.
		p.downloadNextSong()
	}
}

func (p *player) updateStatus() {
	if p.stream == nil {
		p.dg.UpdateStatus(0, "Bot Stopped")
		return
	}
	if p.stream.Paused() {
		p.dg.UpdateStatus(0, fmt.Sprintf("❚❚ %s", p.currentSong.Title))
	} else {
		p.dg.UpdateStatus(0, p.currentSong.Title)
	}
}

func (p *player) Pause() {
	if p.stream != nil {
		p.stream.SetPaused(true)
		p.updateStatus()
	}
}

func (p *player) Play() {
	if p.stream != nil {
		p.stream.SetPaused(false)
		p.updateStatus()
	}
}

func (p *player) downloadNextSong() {
	p.downloadLock.Lock()
	defer p.downloadLock.Unlock()

	nextSong := p.playlist.peekNextSong()
	if nextSong == nil {
		log.Warn("Can't download next song, playlist is empty!")
		return
	}

	songFilePath := path.Join(p.yt.YTCacheDir, nextSong.VideoID+".dca")
	if _, err := os.Stat(filepath.FromSlash(songFilePath)); os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"song": filepath.FromSlash(songFilePath),
		}).Debug("Downloading song")
		p.yt.DownloadDCAAudio(nextSong.VideoID)
	} else {
		log.WithFields(log.Fields{
			"song": filepath.FromSlash(songFilePath),
		}).Debug("Song already downloaded")
	}
}

func (p *player) getNextSongPath() (*songAndPath, error) {
	nextSong := p.playlist.nextSong()
	if nextSong == nil {
		log.Warn("Can't get next song path, playlist is empty!")
	}

	songFilePath := path.Join(p.yt.YTCacheDir, nextSong.VideoID+".dca")
	if _, err := os.Stat(filepath.FromSlash(songFilePath)); err == nil {
		return &songAndPath{songFilePath, nextSong}, nil
	}
	log.WithFields(log.Fields{
		"songpath": songFilePath,
	}).Error("Song is not on disk")
	return nil, fmt.Errorf("Song is not on disk: %s", songFilePath)
}
