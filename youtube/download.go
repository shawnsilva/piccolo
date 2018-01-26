package youtube

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

// DownloadDCAAudio takes a youtube video id, downloads the audio and then
// converts the song to DCA format to be compatible with discordgo.
func (yt Manager) DownloadDCAAudio(videoID string) (string, error) {
	cacheDir := filepath.ToSlash(yt.YTCacheDir)
	outputFilePath := path.Join(cacheDir, "/", videoID+".dca")

	if _, err := os.Stat(filepath.FromSlash(cacheDir)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.FromSlash(cacheDir), os.ModeDir)
		if err != nil {
			// Handle the error
			fmt.Println(err)
		}
	}

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.Application = "audio"
	options.Volume = 125

	videoInfo, err := ytdl.GetVideoInfo(videoID)
	if err != nil {
		// Handle the error
		fmt.Println(err)
	}

	format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := videoInfo.GetDownloadURL(format)
	if err != nil {
		// Handle the error
		fmt.Println(err)
	}

	encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
	if err != nil {
		// Handle the error
		fmt.Println(err)
	}
	defer encodingSession.Cleanup()

	output, err := os.Create(filepath.FromSlash(outputFilePath))
	if err != nil {
		// Handle the error
		fmt.Println(err)
	}

	io.Copy(output, encodingSession)

	return filepath.FromSlash(outputFilePath), nil
}
