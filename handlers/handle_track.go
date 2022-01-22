package handlers

import (
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v3"
)

func (mh *MessageHandler) HandleTrack(c tb.Context, trackId string) error {
	track, err := mh.spotifyClient.GetTrack(trackId)
	if err != nil {
		return err
	}
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf(
		"🎧 %v - %v (%v) \n",
		track.Artists,
		track.Title,
		track.Duration,
	))

	sb.WriteString(fmt.Sprintf(
		"💿 [%v](%v) (%v) \n\n",
		track.AlbumName,
		track.AlbumUrl,
		track.Year,
	))

	sb.WriteString(fmt.Sprintf(
		"⬇️ [Listen on Spotify](%v)",
		track.TrackUrl,
	))

	c.Send(&tb.Photo{
		File:    tb.FromURL(track.AlbumPicUrl),
		Caption: sb.String(),
	}, tb.ModeMarkdown)

	return nil
}
