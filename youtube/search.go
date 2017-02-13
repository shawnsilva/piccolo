package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shawnsilva/piccolo/log"
)

func (yt Manager) createSearchUrl(searchString string) (*string, error) {
	searchUrl, err := url.Parse("https://www.googleapis.com/youtube/v3/search")
	if err != nil {
		return nil, err
	}
	searchParameters := url.Values{}
	searchParameters.Add("part", "snippet")
	searchParameters.Add("q", searchString)
	searchParameters.Add("key", yt.ApiKey)

	searchUrl.RawQuery = searchParameters.Encode()
	searchStr := searchUrl.String()
	return &searchStr, nil
}

func (yt Manager) Search(searchStr string) (YoutubeSearchListResponse, error) {
	var searchResponse YoutubeSearchListResponse
	searchUrl, err := yt.createSearchUrl(searchStr)
	if err != nil {
		return searchResponse, err
	}
	resp, err := http.Get(*searchUrl)
	if err != nil {
		log.Printf("[WARN] Error searching: %s", err)
		return searchResponse, err
	}
	if resp.StatusCode != 200 {
		log.Printf("[WARN] Search failed with status: %s", resp.Status)
		return searchResponse, fmt.Errorf("Got a bad http response: %s", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&searchResponse)

	return searchResponse, nil
}

func (yt Manager) SearchFirstResult(searchStr string) (YoutubeSearchResult, error) {
	var searchResult YoutubeSearchResult
	searchResponseList, err := yt.Search(searchStr)
	if err != nil {
		return searchResult, err
	}
	if len(searchResponseList.Items) == 0 {
		log.Printf("[INFO] Search returned no results: %s", searchStr)
		return searchResult, fmt.Errorf("Search returned no results: %s", searchStr)
	}
	searchResult = searchResponseList.Items[0]
	return searchResult, nil
}
