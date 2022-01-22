package spotify_api

import (
	"fmt"
	"strings"
)

func parseDuration(ms int) string {
	msInMinute := 1000 * 60
	seconds := ms % msInMinute / 1000
	minutes := ms / msInMinute
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func parseArtists(rawTrack *trackInfoResponse, withLinks bool) string {
	sb := strings.Builder{}
	for i, artist := range rawTrack.Artists {
		if withLinks {
			sb.WriteString(fmt.Sprintf("[%v](%v)", artist.Name, artist.Urls.Spotify))
		} else {
			sb.WriteString(fmt.Sprintf("%v", artist.Name))
		}

		if i != len(rawTrack.Artists)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
