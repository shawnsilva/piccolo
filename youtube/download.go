package youtube

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

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

// DownloadAudio takes a string of a youtube video id and invokes youtube-dl through
// os/exec to download the audio track of the video. The file is saved based on the
// location of the video_cache as supplied in the Manager config, with a folder of
// "ytdl" appended.
func (yt Manager) DownloadAudio(videoID string) (string, error) {
	cacheDir := filepath.ToSlash(yt.YTCacheDir)
	outputFileTemplate := path.Join(cacheDir, "/", `%(id)s.%(ext)s`)
	ytdlexecutable := "youtube-dl"
	ytdlpath := filepath.ToSlash(yt.YtDlPath)
	if runtime.GOOS == "windows" {
		ytdlexecutable = fmt.Sprintf("%s.exe", ytdlexecutable)
	}
	ytdlfullpath := path.Join(ytdlpath, ytdlexecutable)

	fmt.Println(outputFileTemplate)
	fmt.Println(ytdlfullpath)

	getFilename := exec.Command(filepath.FromSlash(ytdlfullpath), "--get-filename", "--no-call-home", "-o", outputFileTemplate, "-f", "m4a", fmt.Sprintf("http://youtube.com/watch?v=%s", videoID))
	outputFilepath, err := getFilename.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to determine output path for: %s: %s", videoID, err)
	}

	download := exec.Command(filepath.FromSlash(ytdlfullpath), "--no-call-home", "-o", outputFileTemplate, "-f", "m4a", fmt.Sprintf("http://youtube.com/watch?v=%s", videoID))
	err = download.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to download %s: %s", videoID, err)
	}

	return string(outputFilepath), nil
}
