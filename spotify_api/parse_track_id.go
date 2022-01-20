package spotify_api

import (
	"fmt"
	"strings"
)

func (c *Client) ParseTrackIdFromUrl(str string) (string, error) {
	const spotifySongPrefix = "https://open.spotify.com/track/"

	if !strings.HasPrefix(str, spotifySongPrefix) {
		return "", nil
	}

	trackId := strings.Split(strings.ReplaceAll(str, spotifySongPrefix, ""), "?")
	if len(trackId) < 1 {
		return "", fmt.Errorf("invalid string: %v", str)
	}

	return trackId[0], nil
}
