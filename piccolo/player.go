package piccolo

import (
	"github.com/shawnsilva/piccolo/utils"
)

type (
	player struct {
		conf     *utils.Config
		playlist *playlist
	}
)

func newPlayer(confpointer *utils.Config) *player {
	p := &player{conf: confpointer}
	p.playlist = newPlaylist(p.conf.Bot.UsePlaylist, p.conf.Bot.PlaylistPath)
	return p
}
