package handlers

import (
	"fmt"
	"go-spotify-track-preview-bot/spotify_api"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

func Register(b *tb.Bot, spotifyClient *spotify_api.Client) {
	b.Handle(tb.OnText, func(m tb.Message) {
		log.Println("Got a message")
		trackId, err := spotifyClient.ParseTrackIdFromUrl(m.Text)
		if err != nil || trackId == "" {
			fmt.Println(err.Error())
			return
		}

		track, err := spotifyClient.GetTrack(trackId)
		if err != nil {
			fmt.Println(err.Error())
			return
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

		b.Send(m.Chat, &tb.Photo{
			File:    tb.FromURL(track.AlbumPicUrl),
			Caption: text,
		}, tb.ModeMarkdown)
	})
}
