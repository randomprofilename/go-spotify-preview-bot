package spotify_api

import (
	"fmt"
	"strings"
)

func (c *Client) ParsePlaylistIdFromUrl(str string) (string, error) {
	const spotifyPlaylistPrefix = "https://open.spotify.com/playlist/"

	if !strings.HasPrefix(str, spotifyPlaylistPrefix) {
		return "", nil
	}

	playlistId := strings.Split(strings.ReplaceAll(str, spotifyPlaylistPrefix, ""), "?")
	if len(playlistId) < 1 {
		return "", fmt.Errorf("invalid string: %v", str)
	}

	return playlistId[0], nil
}
