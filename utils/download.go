package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/shawnsilva/piccolo/log"
)

func HttpDownloadToString(url string, desc string) (string, error) {
	log.Printf("[INFO] Downloading: %s", desc)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] Error Downloading %s: %s", desc, err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[ERROR] Error downloading %s with status: %s", desc, resp.Status)
		return "", fmt.Errorf("Got a bad http response: %s", resp.Status)
	}
	httpData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Error reading data for %s: %s", desc, err)
		return "", err
	}
	httpString := string(httpData)
	return httpString, nil
}