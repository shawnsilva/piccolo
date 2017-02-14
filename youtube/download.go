package youtube

import (
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

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
