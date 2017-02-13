package youtube

import (
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

func (yt Manager) DownloadAudio(videoId string) (string, error) {
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

	getFilename := exec.Command(filepath.FromSlash(ytdlfullpath), "--get-filename", "--no-call-home", "-o", outputFileTemplate, "-f", "m4a", fmt.Sprintf("http://youtube.com/watch?v=%s", videoId))
	outputFilepath, err := getFilename.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to determine output path for: %s: %s", videoId, err)
	}

	download := exec.Command(filepath.FromSlash(ytdlfullpath), "--no-call-home", "-o", outputFileTemplate, "-f", "m4a", fmt.Sprintf("http://youtube.com/watch?v=%s", videoId))
	err = download.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to download %s: %s", videoId, err)
	}

	return string(outputFilepath), nil
}
