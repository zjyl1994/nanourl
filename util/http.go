package util

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/zjyl1994/nanourl/vars"
)

func HttpDownload(remoteUrl, localPath string, timeout time.Duration) error {
	hc := http.Client{Timeout: timeout}
	resp, err := hc.Get(remoteUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func DownloadIfNotExist(remoteUrl, localPath string) error {
	if FileExist(localPath) {
		return nil
	}
	return HttpDownload(remoteUrl, localPath, vars.DEFAULT_DOWNLOAD_TIMEOUT)
}
