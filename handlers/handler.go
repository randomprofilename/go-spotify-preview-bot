package handlers

import (
	"go-spotify-track-preview-bot/metrics"
	"go-spotify-track-preview-bot/spotify_api"
	"log"

	tb "gopkg.in/telebot.v3"
)

type MessageHandler struct {
	spotifyClient  *spotify_api.Client
	metricsHandler *metrics.Metrics
}

func NewMessageHandler(spotifyClient *spotify_api.Client, metrics *metrics.Metrics) *MessageHandler {
	return &MessageHandler{
		spotifyClient:  spotifyClient,
		metricsHandler: metrics,
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
			mh.metricsHandler.CountRequests(metrics.GetTrackEvent)
			return mh.HandleTrack(c, trackId)
		}

		playlistId, err := mh.spotifyClient.ParsePlaylistIdFromUrl(c.Text())
		if err != nil {
			return err
		}

		if playlistId != "" {
			mh.metricsHandler.CountRequests(metrics.GetPlaylistEvent)
			return mh.HandlePlaylist(c, playlistId)
		}

		return
	})
}
