package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// BotConfig is used to the the bot specific configuration.
type BotConfig struct {
	Volume                 float64 `json:"volume"`
	YtDlPath               string  `json:"ytdl_path"`
	SaveVideos             bool    `json:"save_videos"`
	CacheDir               string  `json:"cache_dir"`
	UsePlaylist            bool    `json:"use_playlist"`
	PlaylistPath           string  `json:"playlist_path"`
	AutoPause              bool    `json:"auto_pause"`
	DeleteMessages         bool    `json:"delete_messages"`
	DeleteInvokingMessages bool    `json:"delete_invoking_messages"`
	NowPlayingMentions     bool    `json:"now_playing_mentions"`
	SkipsRequired          float64 `json:"skips_required"`
	SkipRatio              float64 `json:"skip_ratio"`
}

// Guilds stores channel information for each guild/server your bot connects too
type Guilds struct {
	BindToTextChannels   []string `json:"bind_to_text_channels"`
	AutoJoinVoiceChannel string   `json:"auto_join_voice_channel"`
}

// Config is used to store the application configuration.
type Config struct {
	BotToken      string    `json:"bot_token"`
	OwnerID       string    `json:"owner_id"`
	GoogleAPIKey  string    `json:"google_api_key"`
	CommandPrefix string    `json:"command_prefix"`
	Bot           BotConfig `json:"bot"`
	Guilds        []Guilds  `json:"guilds"`
}

var (
	defaultBotConfig = BotConfig{
		Volume:                 0.35,
		SaveVideos:             true,
		CacheDir:               "video_cache",
		UsePlaylist:            true,
		PlaylistPath:           "conf/playlist.json",
		AutoPause:              true,
		DeleteMessages:         false,
		DeleteInvokingMessages: false,
		NowPlayingMentions:     true,
		SkipsRequired:          4,
		SkipRatio:              0.5,
	}
	defaultConfig = Config{
		CommandPrefix: "!",
		Bot:           defaultBotConfig,
		Guilds:        []Guilds{},
	}
)

// LoadConfig takes a string for a filename and attempts to load it and
// unmarshal the json inside. Also, auth tokens are attempted to be found as
// environment variables. If successful, returns a pointer to a Config object,
// otherwise returns an error.
func LoadConfig(filename string) (*Config, error) {
	configContents, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		return nil, err
	}
	conf := defaultConfig
	json.Unmarshal(configContents, &conf)
	if os.Getenv("DISCORD_BOT_TOKEN") != "" {
		conf.BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	}
	if os.Getenv("DISCORD_BOT_OWNERID") != "" {
		conf.OwnerID = os.Getenv("DISCORD_BOT_OWNERID")
	}
	if os.Getenv("GOOGLE_API_KEY") != "" {
		conf.GoogleAPIKey = os.Getenv("GOOGLE_API_KEY")
	}
	return &conf, err
}

// DumpConfigFormat will write out a sample config with the default values. It
// be written to the path of the filename string supplied. Returns an error if
// one is encountered.
func DumpConfigFormat(filename string) error {
	jsonConf, _ := json.MarshalIndent(defaultConfig, "", "    ")
	err := ioutil.WriteFile(filepath.FromSlash(filename), jsonConf, 0644)
	if err != nil {
		return err
	}
	return nil
}
