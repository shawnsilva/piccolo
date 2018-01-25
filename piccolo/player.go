package piccolo

import (
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

		currentSongTitle string

		dg *discordgo.Session

		lock         *sync.Mutex
		downloadLock *sync.Mutex
		// chWork       <-chan struct{}
		// chWorkBackup <-chan struct{}
		// chControl    chan struct{}
		// wg           sync.WaitGroup
	}
)

func newPlayer(confpointer *utils.Config, guildID string, voiceChID string, youtube *youtube.Manager, downloadLock *sync.Mutex, discordSession *discordgo.Session) *player {
	p := &player{conf: confpointer, guildID: guildID, voiceChannelID: voiceChID, yt: youtube, dg: discordSession}
	p.playlist = newPlaylist(p.conf.Bot.UsePlaylist, p.conf.Bot.PlaylistPath)
	p.lock = &sync.Mutex{}
	p.downloadLock = downloadLock
	return p
}

func (p *player) Shutdown() error {
	if p.stream != nil {
		p.stream.SetPaused(true)
	}
	// p.Quit()
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
	// p.start()
	go p.playLoop()
	return nil
}

func (p *player) playLoop() {
	for {
		songPath, err := p.getNextSongPath()
		// Download the next song in the background
		go p.downloadNextSong()
		if err == nil {
			reader, err := os.Open(filepath.FromSlash(songPath))
			if err != nil {
				continue
			}
			decoder := dca.NewDecoder(reader)
			p.vc.Speaking(true)
			done := make(chan error)
			p.stream = dca.NewStream(decoder, p.vc, done)
			p.updateStatus()
			streamErr := <-done
			if streamErr != nil && streamErr != io.EOF {
				// Handle the error
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
		p.dg.UpdateStatus(0, fmt.Sprintf("❚❚ %s", p.currentSongTitle))
	} else {
		p.dg.UpdateStatus(0, p.currentSongTitle)
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

// func (p *player) playLoop() {
// 	defer p.wg.Done()
//
// 	for {
// 		select {
// 		case <-p.chWork:
// 			songPath, err := p.getNextSongPath()
// 			// Download the next song in the background
// 			go p.downloadNextSong()
// 			if err == nil {
// 				reader, err := os.Open(filepath.FromSlash(songPath))
// 				if err != nil {
// 					continue
// 				}
// 				decoder := dca.NewDecoder(reader)
// 				done := make(chan error)
// 				p.stream = dca.NewStream(decoder, p.vc, done)
// 				streamErr := <-done
// 				if streamErr != nil && streamErr != io.EOF {
// 					// Handle the error
// 				}
// 			}
// 			// Check if the next song is downloaded, if not block until it is. Catches
// 			// additions to the request queue.
// 			p.downloadNextSong()
// 		case _, ok := <-p.chControl:
// 			if ok {
// 				continue
// 			}
// 			return
// 		}
// 	}
// }

// func (p *player) start() {
// 	ch := make(chan struct{})
// 	close(ch)
// 	p.chWork = ch
// 	p.chWorkBackup = ch
//
// 	p.chControl = make(chan struct{})
//
// 	p.wg = sync.WaitGroup{}
// 	p.wg.Add(1)
//
// 	go p.playLoop()
// }
//
// func (p *player) Pause() {
// 	p.chWork = nil
// 	p.chControl <- struct{}{}
// }
//
// func (p *player) Play() {
// 	p.chWork = p.chWorkBackup
// 	p.chControl <- struct{}{}
// }
//
// func (p *player) Quit() {
// 	p.chWork = nil
// 	close(p.chControl)
// }
//
// func (p *player) Wait() {
// 	p.wg.Wait()
// }

func (p *player) downloadNextSong() {
	p.downloadLock.Lock()
	defer p.downloadLock.Unlock()
	var songID string
	nextSong := p.playlist.peekNextSong()
	playlistSong, ok := nextSong.(*PlaylistEntry)
	if !ok {
		requestSong, reqOk := nextSong.(*requestEntry)
		if !reqOk {
			log.WithFields(log.Fields{
				"song": nextSong,
			}).Error("peekNextSong is invalid")
			return
		}
		songID = requestSong.song.VideoID
	} else {
		songID = playlistSong.VideoID
	}
	songFilePath := path.Join(p.yt.YTCacheDir, songID+".dca")
	if _, err := os.Stat(filepath.FromSlash(songFilePath)); os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"song": filepath.FromSlash(songFilePath),
		}).Debug("Downloading song")
		p.yt.DownloadDCAAudio(songID)
	} else {
		log.WithFields(log.Fields{
			"song": filepath.FromSlash(songFilePath),
		}).Debug("Song already downloaded")
	}
}

func (p *player) getNextSongPath() (string, error) {
	var songID string
	var songTile string
	nextSong := p.playlist.nextSong()
	playlistSong, ok := nextSong.(*PlaylistEntry)
	if !ok {
		requestSong, reqOk := nextSong.(*requestEntry)
		if !reqOk {
			log.WithFields(log.Fields{
				"song": nextSong,
			}).Error("NextSong is invalid")
			p.currentSongTitle = ""
			return "", fmt.Errorf("NextSong is invalid: %s", nextSong)
		}
		songID = requestSong.song.VideoID
		songTile = requestSong.song.Title
	} else {
		songID = playlistSong.VideoID
		songTile = playlistSong.Title
	}
	songFilePath := path.Join(p.yt.YTCacheDir, songID+".dca")
	if _, err := os.Stat(filepath.FromSlash(songFilePath)); err == nil {
		p.currentSongTitle = songTile
		return songFilePath, nil
	}
	log.WithFields(log.Fields{
		"songpath": songFilePath,
	}).Error("Song is not on disk")
	p.currentSongTitle = ""
	return "", fmt.Errorf("Song is not on disk: %s", songFilePath)
}
