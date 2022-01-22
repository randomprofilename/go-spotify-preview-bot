package main

import (
	"go-spotify-track-preview-bot/handlers"
	"go-spotify-track-preview-bot/spotify_api"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/telebot.v3"
)

type appConfig struct {
	token               string
	port                string
	webhookUrl          string
	spotifyClientId     string
	spotifyClientSecret string
}

func createConfig() appConfig {
	if err := godotenv.Load(); err != nil {
		log.Print("File .env doesn't exist")
	} else {
		log.Print("File .env was loaded")
	}

	envToken, ok := os.LookupEnv("TG_TOKEN")
	if !ok {
		log.Fatal("Provide token in env variable 'TG_TOKEN'")
	}

	envPort, _ := os.LookupEnv("PORT")
	envWebhookUrl, _ := os.LookupEnv("WEBHOOK_URL")

	envSpotifyClientId, ok := os.LookupEnv("SPOTIFY_CLIENT_ID")
	if !ok {
		log.Fatal("Provide token in env variable 'SPOTIFY_CLIENT_ID'")
	}

	envSpotifyClientSecret, ok := os.LookupEnv("SPOTIFY_CLIENT_SECRET")
	if !ok {
		log.Fatal("Provide token in env variable 'SPOTIFY_CLIENT_SECRET'")
	}

	return appConfig{
		port:                envPort,
		token:               envToken,
		webhookUrl:          envWebhookUrl,
		spotifyClientId:     envSpotifyClientId,
		spotifyClientSecret: envSpotifyClientSecret,
	}
}

func main() {
	config := createConfig()

	var poller tb.Poller

	if config.webhookUrl != "" {
		poller = &tb.Webhook{
			Listen:   ":" + config.port,
			Endpoint: &tb.WebhookEndpoint{PublicURL: config.webhookUrl},
		}
	} else {
		poller = &tb.LongPoller{Timeout: 10 * time.Second}
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  config.token,
		Poller: poller,
	})
	if config.webhookUrl == "" {
		b.RemoveWebhook()
	}

	rand.Seed(time.Now().Unix())

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print("Connected...")

	spotifyClient := spotify_api.NewClient(config.spotifyClientId, config.spotifyClientSecret)
	err = spotifyClient.UpdateToken()
	if err != nil {
		log.Fatal(err)
	}
	if config.webhookUrl == "" {
		b.RemoveWebhook()
	}

	mh := handlers.NewMessageHandler(spotifyClient)

	mh.Register(b)

	b.Start()
	log.Print("Started...")
}
