package handlers

import (
	"fmt"
	"go-spotify-track-preview-bot/spotify_api"
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

	playlistTextWithUrls := getTrackListWithLinks(playlist.Tracks, numberOfTracks, playlist.Total)
	if len(playlistTextWithUrls) < maxMessageLength {
		sb.WriteString(playlistTextWithUrls)
	} else {
		sb.WriteString(getTrackListWithoutLinks(playlist.Tracks, numberOfTracks, playlist.Total))
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("⬇️ [Listen on Spotify](%v)", playlist.PlaylistUrl))

	return c.Send(&tb.Photo{
		File:    tb.FromURL(playlist.PlaylistPicUrl),
		Caption: sb.String(),
	}, tb.ModeMarkdown)
}

func getTrackListWithLinks(tracks []*spotify_api.Track, maxNumber int, total int) string {
	sb := strings.Builder{}
	for i, track := range tracks {
		if track.TrackUrl == "" {
			continue
		}
		if !(i < maxNumber) {
			sb.WriteString(fmt.Sprintf("• _and %d more..._\n", total-maxNumber))
			break
		}

		sb.WriteString(fmt.Sprintf(
			"• %v - [%v](%v) (%v)\n",
			track.Artists,
			track.Title,
			track.TrackUrl,
			track.Duration,
		))
	}
	return sb.String()
}

func getTrackListWithoutLinks(tracks []*spotify_api.Track, maxNumber int, total int) string {
	sb := strings.Builder{}
	for i, track := range tracks {
		if track.TrackUrl == "" {
			continue
		}
		if !(i < maxNumber) {
			sb.WriteString(fmt.Sprintf("• _and %d more..._\n", total-maxNumber))
			break
		}

		sb.WriteString(fmt.Sprintf(
			"• %v - %v (%v)\n",
			track.Artists,
			track.Title,
			track.Duration,
		))
	}
	return sb.String()
}
