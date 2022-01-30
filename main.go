package main

import (
	"context"
	"go-spotify-track-preview-bot/handlers"
	"go-spotify-track-preview-bot/metrics"
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

	PromConfig metrics.PromConfig
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

	envPromURI, _ := os.LookupEnv("PROM_REMOTE_WRITE_URI")
	envPromUsername, _ := os.LookupEnv("PROM_USERNAME")
	envPromPassword, _ := os.LookupEnv("PROM_PASSWORD")

	return appConfig{
		port:                envPort,
		token:               envToken,
		webhookUrl:          envWebhookUrl,
		spotifyClientId:     envSpotifyClientId,
		spotifyClientSecret: envSpotifyClientSecret,

		PromConfig: metrics.PromConfig{
			PrometheusRemoteWriteURI:      envPromURI,
			PrometheusRemoteWriteUsername: envPromUsername,
			PrometheusRemoteWritePassword: envPromPassword,
		},
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := createConfig()

	var poller tb.Poller

	if config.webhookUrl != "" {
		log.Print("Connecting via webhook")
		poller = &tb.Webhook{
			Listen:   ":" + config.port,
			Endpoint: &tb.WebhookEndpoint{PublicURL: config.webhookUrl},
		}
	} else {
		log.Print("Connecting via longpoller")
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

	metricsHandler := metrics.InitMetrics(ctx, config.PromConfig)
	mh := handlers.NewMessageHandler(spotifyClient, metricsHandler)

	mh.Register(b)
	log.Print("Started...")
	b.Start()
}
