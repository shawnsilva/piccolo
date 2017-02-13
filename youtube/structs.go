package youtube

type (
	Manager struct {
		ApiKey     string
		YtDlPath   string
		YTCacheDir string
	}

	thumbnailInfo struct {
		URL    string  `json:"url"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}

	YoutubeSearchResult struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind       string `json:"kind"`
			VideoID    string `json:"videoId"`
			ChannelID  string `json:"channelId"`
			PlaylistID string `json:"playlistId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt string `json:"publishedAt"`
			ChannelID   string `json:"channelId"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Thumbnails  struct {
				Default thumbnailInfo `json:"default"`
				Medium  thumbnailInfo `json:"medium"`
				High    thumbnailInfo `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string `json:"channelTitle"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	}

	YoutubeSearchListResponse struct {
		Kind          string `json:"kind"`
		Etag          string `json:"etag"`
		NextPageToken string `json:"nextPageToken"`
		PrevPageToken string `json:"prevPageToken"`
		RegionCode    string `json:"regionCode`
		PageInfo      struct {
			TotalResults   float64 `json:"totalResults"`
			ResultsPerPage float64 `json:"resutlsPerPage"`
		} `json:"pageInfo"`
		Items []YoutubeSearchResult `json:"items"`
	}

	VideoFormatInfo map[string]string
	YoutubeVideo    struct {
		ID      string
		Title   string
		Url     string
		Formats []VideoFormatInfo
	}
)
