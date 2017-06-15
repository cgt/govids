package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

type Playlist struct {
	Items []struct {
		Snippet struct {
			ResourceID struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
			Title       string `json:"title"`
			PublishedAt string `json:"publishedAt"`
		} `json:"snippet"`
	} `json:"items"`
}

type Video struct {
	Title string    `json:"title"`
	ID    string    `json:"id"`
	Date  time.Time `json:"date"`
}

func main() {
	var tag string
	flag.StringVar(&tag, "tag", "", "tag to add to all playlist items")
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	playlistID := flag.Arg(0)

	apikey, ok := os.LookupEnv("YOUTUBEAPIKEY")
	if !ok {
		log.Fatalln("YOUTUBEAPIKEY not set")
	}

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
		log.Fatalf("GET error: %v", err)
	}
	defer resp.Body.Close()

	var pl Playlist
	err = json.NewDecoder(resp.Body).Decode(&pl)
	if err != nil {
		resp.Body.Close()
		log.Fatalf("JSON decode error: %v", err)
	}

	videos := make([]Video, 0, len(pl.Items))
	for _, item := range pl.Items {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", item.Snippet.PublishedAt)
		if err != nil {
			log.Fatalf("error parsing timestamp: %v", err)
		}
		v := Video{
			Title: item.Snippet.Title,
			ID:    item.Snippet.ResourceID.VideoID,
			Date:  t,
		}
		videos = append(videos, v)
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Date.Before(videos[j].Date)
	})

	today := time.Now().Format("2006-01-02")

	for i, v := range videos {
		fmt.Printf(`{
   "date": "%s",
   "added": "%s",
   "id": "%s",
   "title": %q,
   "speakers": [
   ],
   "tags": [
`, v.Date.Format("2006-01-02"), today, v.ID, v.Title)

		if tag != "" {
			fmt.Printf("      %q\n", tag)
		}
		fmt.Print("   ]\n}")
		if i < len(videos)-1 {
			fmt.Print(",")
		}
		fmt.Println()
	}
}
