package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shawnsilva/piccolo/log"
)

func (yt Manager) createSearchURL(searchString string) (*string, error) {
	searchURL, err := url.Parse("https://www.googleapis.com/youtube/v3/search")
	if err != nil {
		return nil, err
	}
	searchParameters := url.Values{}
	searchParameters.Add("part", "snippet")
	searchParameters.Add("q", searchString)
	searchParameters.Add("key", yt.APIKey)

	searchURL.RawQuery = searchParameters.Encode()
	searchStr := searchURL.String()
	return &searchStr, nil
}

// Search takes a string input and searches youtube for results. Search returns a
// YoutubeSearchListResponse with the results.
func (yt Manager) Search(searchStr string) (SearchListResponse, error) {
	var searchResponse SearchListResponse
	searchURL, err := yt.createSearchURL(searchStr)
	if err != nil {
		return searchResponse, err
	}
	resp, err := http.Get(*searchURL)
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

// SearchFirstResult takes a string input and searches youtube, returning only
// the first result in a YoutubeSearchResult
func (yt Manager) SearchFirstResult(searchStr string) (SearchResult, error) {
	var searchResult SearchResult
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
