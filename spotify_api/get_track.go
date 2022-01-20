package spotify_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Track struct {
	Artists   string
	Title     string
	Duration  string
	AlbumName string
	AlbumUrl  string
	Year      string

	AlbumPicUrl string
	TrackUrl    string
}

type trackInfoResponse struct {
	Album struct {
		Name   string `json:"name"`
		Images []struct {
			Url string `json:"url"`
		} `json:"images"`
		ReleaseDate string `json:"release_date"`
		Urls        struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	} `json:"album"`

	Artists []struct {
		Name string `json:"name"`
		Urls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	} `json:"artists"`

	Urls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`

	Uri        string `json:"uri"`
	Name       string `json:"name"`
	DurationMs int    `json:"duration_ms"`
}

func getUrl(trackId string) string {
	return fmt.Sprintf(
		"https://api.spotify.com/v1/tracks/%v",
		trackId,
	)
}

func parseDuration(ms int) string {
	msInMinute := 1000 * 60
	seconds := ms % msInMinute / 1000
	minutes := ms / msInMinute

	return fmt.Sprintf("%v:%v", minutes, seconds)
}

func parseArtists(rawTrack *trackInfoResponse) string {
	sb := strings.Builder{}
	for i, artist := range rawTrack.Artists {
		sb.WriteString(fmt.Sprintf("[%v](%v)", artist.Name, artist.Urls.Spotify))
		if i != len(rawTrack.Artists)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (c *Client) GetTrack(trackId string) (*Track, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", getUrl(trackId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	rawTrack := &trackInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(rawTrack)

	track := &Track{
		Artists:   parseArtists(rawTrack),
		Title:     rawTrack.Name,
		Duration:  parseDuration(rawTrack.DurationMs),
		AlbumName: rawTrack.Album.Name,
		Year:      rawTrack.Album.ReleaseDate,

		AlbumUrl:    rawTrack.Album.Urls.Spotify,
		AlbumPicUrl: rawTrack.Album.Images[0].Url,
		TrackUrl:    rawTrack.Urls.Spotify,
	}

	return track, err
}
