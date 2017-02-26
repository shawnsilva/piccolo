package piccolo

import (
	"github.com/shawnsilva/piccolo/utils"
)

type (
	playlistEntry struct {
		title   string
		videoID string
	}

	requestEntry struct {
		requester string
		song      playlistEntry
	}

	requestQueue struct {
		*utils.Queue
	}
)

func newRequestQueue() *requestQueue {
	q := &requestQueue{utils.NewQueue()}
	return q
}

func (q requestQueue) addSong(requester string, id string, title string) {
	q.Push(requestEntry{
		requester: requester,
		song: playlistEntry{
			title:   title,
			videoID: id,
		},
	})
}

func (q requestQueue) nextSong() *requestEntry {
	for {
		if q.Length() <= 0 {
			break
		}
		song, ok := q.Pop().(requestEntry)
		if !ok {
			continue
		}
		return &song
	}
	return nil
}
