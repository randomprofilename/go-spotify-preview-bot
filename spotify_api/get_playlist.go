package spotify_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Playlist struct {
	Tracks      []*Track
	Name        string
	Description string

	PlaylistPicUrl string
	PlaylistUrl    string
	Total          int
}

type playlistInfoResponse struct {
	Tracks struct {
		Items []struct {
			Track trackInfoResponse `json:"track"`
		} `json:"items"`
		Total int `json:"total"`
	} `json:"tracks"`

	Urls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`

	Images []struct {
		Url string `json:"url"`
	} `json:"images"`

	Name        string `json:"name"`
	Description string `json:"description"`
}

func getPlaylistUrl(playlistId string) string {
	return fmt.Sprintf(
		"https://api.spotify.com/v1/playlists/%v",
		playlistId,
	)
}

func (c *Client) GetPlaylist(playlistId string) (*Playlist, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", getPlaylistUrl(playlistId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	rawPlaylist := &playlistInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(rawPlaylist)

	pic := ""
	if len(rawPlaylist.Images) > 0 {
		pic = rawPlaylist.Images[0].Url
	}

	tracks := make([]*Track, 0, len(rawPlaylist.Tracks.Items))
	for _, rawTrack := range rawPlaylist.Tracks.Items {
		track := parseTrack(&rawTrack.Track, pic)
		track.Artists = parseArtists(&rawTrack.Track, false)

		tracks = append(tracks, track)
	}

	return &Playlist{
		Tracks: tracks,

		Name:           rawPlaylist.Name,
		Description:    rawPlaylist.Description,
		PlaylistPicUrl: pic,
		PlaylistUrl:    rawPlaylist.Urls.Spotify,
		Total:          rawPlaylist.Tracks.Total,
	}, err
}
