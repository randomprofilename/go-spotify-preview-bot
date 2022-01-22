package handlers

import (
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v3"
)

func (mh *MessageHandler) HandlePlaylist(c tb.Context, playlistId string) error {
	const maxMessageLength = 1024

	playlist, err := mh.spotifyClient.GetPlaylist(playlistId)
	if err != nil {
		return err
	}

	sb := strings.Builder{}

	if playlist.Owner != "" {
		sb.WriteString(fmt.Sprintf("🎶 *%v* by _%v_\n", playlist.Name, playlist.Owner))
	} else {
		sb.WriteString(fmt.Sprintf("🎶 *%v* \n", playlist.Name))
	}

	if playlist.Description != "" {
		sb.WriteString(fmt.Sprintf("💬 _%v_ \n", playlist.Description))
	}

	sb.WriteString("\n")
	const numberOfTracks = 6

	for i, track := range playlist.Tracks {
		if track.TrackUrl == "" {
			continue
		}
		if !(i < numberOfTracks) {
			sb.WriteString(fmt.Sprintf("• _and %d more..._\n", playlist.Total-numberOfTracks))
			break
		}

		newTrackText := fmt.Sprintf(
			"• %v - %v (%v)\n",
			track.Artists,
			track.Title,
			track.Duration,
		)

		if len(newTrackText)+sb.Len() < maxMessageLength {
			sb.WriteString(newTrackText)
		}
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("⬇️ [Listen on Spotify](%v)", playlist.PlaylistUrl))

	return c.Send(&tb.Photo{
		File:    tb.FromURL(playlist.PlaylistPicUrl),
		Caption: sb.String(),
	}, tb.ModeMarkdown)
}
