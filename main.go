package main

import (
	"fmt"
	"strings"
)

const SLACK_API_KEY = "xxx"
const SPOTIFY_CLIENT_AUTH = "xxx"
const SPOTIFY_REFRESH_TOKEN = "xxx"

var spotify Spotify

func handler(s *SlackBot, message MessageJSON) {
	if strings.HasPrefix(message.Text, s.user) {
		fmt.Println("Message Received:", message.Text)

		switch {
		case strings.Contains(message.Text, "current"):
			s.Send(`_Currently Playing "`+spotify.Current()+`"_`, message.Channel)
		case strings.Contains(message.Text, "play"):
			s.Send(`_Playing Spotify..._`, message.Channel)
		case strings.Contains(message.Text, "pause"):
			s.Send(`_Pausing Spotify..._`, message.Channel)
		case strings.Contains(message.Text, "skip"):
			s.Send(`_Skipping this song..._`, message.Channel)
		case strings.Contains(message.Text, "last"):
			s.Send(`_Playing last song..._`, message.Channel)
		case strings.Contains(message.Text, "restart"):
			s.Send(`_Restarting this song..._`, message.Channel)
		case strings.Contains(message.Text, "joke"):
			s.Send(Joke{}.Get(), message.Channel)
		}
	}
}

func main() {
	// log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	spotify = Spotify{
		client:  SPOTIFY_CLIENT_AUTH,
		refresh: SPOTIFY_REFRESH_TOKEN,
	}

	spotify.Connect()

	slackbot := SlackBot{}
	slackbot.Connect("https://slack.com/api/rtm.start?token=" + SLACK_API_KEY)
	slackbot.Subscribe(handler)
}
