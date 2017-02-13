package utils

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type BotConfig struct {
	Volume                 float64 `json:"volume"`
	YtDlPath               string  `json:"ytdl_path"`
	SaveVideos             bool    `json:"save_videos"`
	CacheDir               string  `json:"cache_dir"`
	UsePlaylist            bool    `json:"use_playlist"`
	AutoPause              bool    `json:"auto_pause"`
	DeleteMessages         bool    `json:"delete_messages"`
	DeleteInvokingMessages bool    `json:"delete_invoking_messages"`
	NowPlayingMentions     bool    `json:"now_playing_mentions"`
	SkipsRequired          float64 `json:"skips_required"`
	SkipRatio              float64 `json:"skip_ratio"`
}

type Config struct {
	BotToken              string    `json:"bot_token"`
	OwnerId               string    `json:"owner_id"`
	GoogleApiKey          string    `json:"google_api_key"`
	BindToTextChannels    []string  `json:"bind_to_text_channels"`
	AutoJoinVoiceChannels []string  `json:"auto_join_voice_channels"`
	CommandPrefix         string    `json:"command_prefix"`
	Bot                   BotConfig `json:"bot"`
}

var (
	defaultBotConfig = BotConfig{
		Volume:                 0.35,
		SaveVideos:             true,
		CacheDir:               "video_cache",
		UsePlaylist:            true,
		AutoPause:              true,
		DeleteMessages:         false,
		DeleteInvokingMessages: false,
		NowPlayingMentions:     true,
		SkipsRequired:          4,
		SkipRatio:              0.5,
	}
	defaultConfig = Config{
		CommandPrefix:         "!",
		BindToTextChannels:    []string{"none"},
		AutoJoinVoiceChannels: []string{"none"},
		Bot: defaultBotConfig,
	}
)

func ParseConfigFile(filename string) (*Config, error) {
	configContents, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		return nil, err
	}
	conf := defaultConfig
	json.Unmarshal(configContents, &conf)
	return &conf, err
}

func DumpConfigFormat(filename string) error {
	jsonConf, _ := json.MarshalIndent(defaultConfig, "", "    ")
	err := ioutil.WriteFile(filepath.FromSlash(filename), jsonConf, 0644)
	if err != nil {
		return err
	}
	return nil
}
