package youtube

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

// DownloadDCAAudio takes a youtube video id, downloads the audio and then
// converts the song to DCA format to be compatible with discordgo.
func (yt Manager) DownloadDCAAudio(videoID string) (string, time.Duration, error) {
	cacheDir := filepath.ToSlash(yt.YTCacheDir)
	outputFilePath := path.Join(cacheDir, "/", videoID+".dca")

	if _, err := os.Stat(filepath.FromSlash(cacheDir)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.FromSlash(cacheDir), os.ModeDir)
		if err != nil {
			return "", time.Duration(0), err
		}
	}

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.Application = "audio"
	options.Volume = 125

	videoInfo, err := ytdl.GetVideoInfo(videoID)
	if err != nil {
		return "", time.Duration(0), err
	}

	formats := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)
	if len(formats) < 1 {
		return "", time.Duration(0), fmt.Errorf("Couldn't Find Audiot Formats")
	}
	format := formats[0]
	downloadURL, err := videoInfo.GetDownloadURL(format)
	if err != nil {
		return "", time.Duration(0), err
	}

	encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
	if err != nil {
		return "", time.Duration(0), err
	}
	defer encodingSession.Cleanup()

	output, err := os.Create(filepath.FromSlash(outputFilePath))
	if err != nil {
		return "", time.Duration(0), err
	}

	_, err = io.Copy(output, encodingSession)
	if err != nil {
		return "", time.Duration(0), err
	}

	return filepath.FromSlash(outputFilePath), encodingSession.Stats().Duration, nil
}
