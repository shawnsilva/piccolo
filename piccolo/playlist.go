package piccolo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/jatgam/goutils"
	"github.com/jatgam/goutils/log"
)

type (
	// PlaylistEntry is an individual song in the playlist
	PlaylistEntry struct {
		Requester        *discordgo.User `json:"-"`
		RequestChannelID string          `json:"-"`
		Title            string          `json:"title"`
		VideoID          string          `json:"videoID"`
		TrackDuration    time.Duration   `json:"duration,omitempty"`
	}

	playlist struct {
		requestQueue          *goutils.Queue
		list                  *goutils.DoubleLinkedList
		current               *goutils.Node
		usePlaylist           bool
		playlistPath          string
		readWritePlaylistLock *sync.Mutex
	}

	// PlaylistJSON is used to handled marshalling and unmarshalling a playlist
	// to a file on disk
	PlaylistJSON struct {
		Entries []PlaylistEntry `json:"entries"`
	}
)

func newPlaylist(usePlaylist bool, playlistPath string) *playlist {
	p := &playlist{requestQueue: goutils.NewQueue(), list: goutils.NewDoubleLinkedList(),
		usePlaylist: usePlaylist, playlistPath: playlistPath}
	p.readWritePlaylistLock = &sync.Mutex{}
	p.loadPlaylist()
	p.current = p.list.First()

	return p
}

func (p *playlist) loadPlaylist() error {
	if p.usePlaylist {
		p.readWritePlaylistLock.Lock()
		defer p.readWritePlaylistLock.Unlock()
		playlistFileContents, err := ioutil.ReadFile(filepath.FromSlash(p.playlistPath))
		if err != nil {
			log.WithFields(log.Fields{
				"file":  filepath.FromSlash(p.playlistPath),
				"error": err,
			}).Error("Failed to open playlist file to read")
			return fmt.Errorf("Couldn't read the playlist file")
		}
		var filePlaylist = PlaylistJSON{}
		jsonErr := json.Unmarshal(playlistFileContents, &filePlaylist)
		if jsonErr != nil {
			log.WithFields(log.Fields{
				"file":  filepath.FromSlash(p.playlistPath),
				"error": err,
			}).Error("Failed to decode playlist json")
			return fmt.Errorf("Couldn't decode the playlist file")
		}
		for _, entry := range filePlaylist.Entries {
			p.list.InsertEnd(goutils.NewNode(entry.VideoID, entry))
		}
	} else {
		log.Debug("Attempted to load a playlist when use is disabled in config file.")
		return fmt.Errorf("Using a playlist is currently disabled via the config file")
	}
	log.Debug("Loaded playlist")
	return nil
}

func (p *playlist) savePlaylist() error {
	if p.usePlaylist {
		p.readWritePlaylistLock.Lock()
		defer p.readWritePlaylistLock.Unlock()
		currentPlaylist := &PlaylistJSON{Entries: []PlaylistEntry{}}
		currentNode := p.list.First()
		for {
			if currentNode == nil {
				break
			}
			_, value := currentNode.GetData()
			if song, ok := value.(PlaylistEntry); ok {
				currentPlaylist.Entries = append(currentPlaylist.Entries, song)
			}
			currentNode = currentNode.Next()
		}
		jsonPlaylist, _ := json.MarshalIndent(currentPlaylist, "", "    ")
		err := ioutil.WriteFile(filepath.FromSlash(p.playlistPath), jsonPlaylist, 0644)
		if err != nil {
			log.WithFields(log.Fields{
				"file":  filepath.FromSlash(p.playlistPath),
				"error": err,
			}).Error("Failed to write playlist json file")
			return fmt.Errorf("Error saving playlist file")
		}
	} else {
		log.Debug("Attempted to save a playlist when use is disabled in config file.")
		return fmt.Errorf("Using a playlist is currently disabled via the config file")
	}
	log.Debug("Saved playlist")
	return nil
}

func (p *playlist) String() string {
	return p.printPlaylist("")
}

func (p *playlist) printPlaylist(currentVideoID string) string {
	var queueString string
	var playlistString string
	queueItem := p.requestQueue.First()
	count := 1
	if queueItem == nil {
		queueString = "    Empty"
	} else {
		for {
			if queueItem == nil {
				break
			}
			song, ok := queueItem.Data().(PlaylistEntry)
			if !ok {
				continue
			}
			queueString = queueString + fmt.Sprintf("    %d. %s - Requester: %s\n",
				count, song.Title, song.Requester.Username)
			queueItem = queueItem.Next()
			count++
		}
	}
	count = 1
	if p.usePlaylist {
		currentNode := p.list.First()
		for {
			if currentNode == nil {
				break
			}
			_, value := currentNode.GetData()
			if song, ok := value.(PlaylistEntry); ok {
				if currentVideoID == song.VideoID {
					playlistString = playlistString + fmt.Sprintf("→   %d. %s\n",
						count, song.Title)
				} else {
					playlistString = playlistString + fmt.Sprintf("    %d. %s\n",
						count, song.Title)
				}
				count++
			}
			currentNode = currentNode.Next()
		}
	} else {
		playlistString = "    Disabled"
	}

	return fmt.Sprintf("**Request Queue:**\n```%s```\n**Playlist:**\n```%s```",
		queueString, playlistString)
}

func (p *playlist) addSong(requester *discordgo.User, channelID string, id string, title string) {
	p.requestQueue.Push(PlaylistEntry{
		Requester:        requester,
		RequestChannelID: channelID,
		Title:            title,
		VideoID:          id,
	})
}

func (p *playlist) addRequestedSong(pEntry *PlaylistEntry) {
	p.requestQueue.Push(PlaylistEntry{
		Requester:        pEntry.Requester,
		RequestChannelID: pEntry.RequestChannelID,
		Title:            pEntry.Title,
		VideoID:          pEntry.VideoID,
		TrackDuration:    pEntry.TrackDuration,
	})
}

func (p *playlist) nextSong() *PlaylistEntry {
	for {
		if p.requestQueue.Length() <= 0 {
			if p.list.Length() <= 0 {
				break
			} else {
				if p.current == nil {
					if p.list.First() != nil {
						p.current = p.list.First()
					} else {
						break
					}
				}
				_, songData := p.current.GetData()
				p.current = p.current.Next()
				song, ok := songData.(PlaylistEntry)
				if !ok {
					continue
				}
				return &song
			}
		}
		song, ok := p.requestQueue.Pop().(PlaylistEntry)
		if !ok {
			continue
		}
		return &song
	}
	return nil
}

func (p *playlist) peekNextSong() *PlaylistEntry {
	for {
		if p.requestQueue.Length() <= 0 {
			if p.list.Length() <= 0 {
				break
			} else {
				if p.current == nil {
					if p.list.First() != nil {
						p.current = p.list.First()
					} else {
						break
					}
				}
				_, songData := p.current.GetData()
				song, ok := songData.(PlaylistEntry)
				if !ok {
					continue
				}
				return &song
			}
		}
		song, ok := p.requestQueue.Look().(PlaylistEntry)
		if !ok {
			continue
		}
		return &song
	}
	return nil
}

func (p *playlist) modifyPlaylistEntry(videoID string, dataContents PlaylistEntry) error {
	p.readWritePlaylistLock.Lock()
	defer p.readWritePlaylistLock.Unlock()
	playlistEntry := p.list.Find(videoID)
	if playlistEntry != nil {
		playlistEntry.SetData(videoID, dataContents)
	} else {
		return fmt.Errorf("Couldn't find playlist entry to modify: %s", videoID)
	}
	return nil
}

func (p *playlist) removePlaylistEntry(videoID string) {
	p.readWritePlaylistLock.Lock()
	defer p.readWritePlaylistLock.Unlock()
	p.list.Delete(videoID)
}
