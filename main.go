package main

import (
	"fmt"
	"strings"
)

const SLACK_API_KEY = ""
const SPOTIFY_CLIENT_AUTH = ""

func handler(s *SlackBot, message MessageJSON) {
	fmt.Println("Message Received:", message.Text)
	// s.Send(message.Text, message.Channel)
	if strings.HasPrefix(message.Text, "/spotify") {
		fmt.Println("yay")
	}
}

func main() {
	// log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	slackbot := SlackBot{}
	slackbot.Connect("https://slack.com/api/rtm.start?token=" + SLACK_API_KEY)
	slackbot.Subscribe(handler)
  
  code := ""
  
  spotify := Spotify{
    client: SPOTIFY_CLIENT_AUTH,
  }
  
  spotify.Connect(code)
  spotify.Current()
}
