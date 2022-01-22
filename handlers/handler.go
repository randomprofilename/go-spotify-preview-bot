package handlers

import (
	"go-spotify-track-preview-bot/spotify_api"
	"log"

	tb "gopkg.in/telebot.v3"
)

type MessageHandler struct {
	spotifyClient *spotify_api.Client
}

func NewMessageHandler(spotifyClient *spotify_api.Client) *MessageHandler {
	return &MessageHandler{
		spotifyClient: spotifyClient,
	}
}

func (mh *MessageHandler) Register(b *tb.Bot) {
	b.Handle(tb.OnText, func(c tb.Context) (err error) {
		log.Println("Got a message")

		trackId, err := mh.spotifyClient.ParseTrackIdFromUrl(c.Text())
		if err != nil {
			return err
		}

		if trackId != "" {
			return mh.HandleTrack(c, trackId)
		}

		playlistId, err := mh.spotifyClient.ParsePlaylistIdFromUrl(c.Text())
		if err != nil {
			return err
		}

		if playlistId != "" {
			return mh.HandlePlaylist(c, playlistId)
		}

		return
	})
}
