package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	apikey, ok := os.LookupEnv("YOUTUBEAPIKEY")
	if !ok {
		log.Fatalln("YOUTUBEAPIKEY not set")
	}
	playlistID := os.Args[1]

	params := url.Values{
		"part":       {"snippet"},
		"maxResults": {"50"},
		"playlistId": {playlistID},
		"fields":     {"items"},
		"key":        {apikey},
	}
	u, err := url.Parse("https://www.googleapis.com/youtube/v3/playlistItems")
	if err != nil {
		panic("error parsing hardcoded API URL")
	}
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatalf("GET error: %v\n", err)
	}
	defer resp.Body.Close()

	var pl Playlist
	err = json.NewDecoder(resp.Body).Decode(&pl)
	if err != nil {
		resp.Body.Close()
		log.Fatalf("JSON decode error: %v\n", err)
	}

	output := make([]Video, 0, len(pl.Items))
	for _, item := range pl.Items {
		v := Video{
			Title: item.Snippet.Title,
			ID:    item.Snippet.ResourceID.VideoID,
		}
		output = append(output, v)
	}
	err = json.NewEncoder(os.Stdout).Encode(output)
	if err != nil {
		log.Fatalf("JSON encode error: %v\n", err)
	}
}

type Playlist struct {
	Items []struct {
		Snippet struct {
			ResourceID struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

type Video struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}
