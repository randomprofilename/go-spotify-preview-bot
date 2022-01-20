package handlers

import (
	"fmt"
	"go-spotify-track-preview-bot/spotify_api"
	"log"

	tb "gopkg.in/telebot.v3"
)

func Register(b *tb.Bot, spotifyClient *spotify_api.Client) {
	b.Handle(tb.OnText, func(c tb.Context) (err error) {
		log.Println("Got a message")
		trackId, err := spotifyClient.ParseTrackIdFromUrl(c.Text())
		if err != nil || trackId == "" {
			return err
		}

		track, err := spotifyClient.GetTrack(trackId)
		if err != nil {
			return err
		}

		text := fmt.Sprintf(
			"🎧 %v - %v (%v) \n💿 [%v](%v) (%v) \n⬇️ [Listen on Spotify](%v)",
			track.Artists,
			track.Title,
			track.Duration,
			track.AlbumName,
			track.AlbumUrl,
			track.Year,
			track.TrackUrl,
		)

		err = c.Send(&tb.Photo{
			File:    tb.FromURL(track.AlbumPicUrl),
			Caption: text,
		}, tb.ModeMarkdown)

		return
	})
}
