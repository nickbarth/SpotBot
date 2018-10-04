package main

import (
	"fmt"
	"log"
	"strings"
)

const SLACK_API_KEY = ""
const SPOTIFY_CLIENT_AUTH = ""
const SPOTIFY_REFRESH_TOKEN = ""
const SPOTIFY_DEVICE = ""
const PLAYLIST_ID = ""

var spotify Spotify

func handler(s *SlackBot, message MessageJSON) {
	if strings.HasPrefix(message.Text, s.user) {
		fmt.Println("Message Received:", message.Text)

		args := strings.Split(message.Text, " ")
		command := args[1]

		switch command {
		case "current":
			s.Send(`_Currently Playing "`+spotify.Current()+`"_`, message.Channel)
		case "play":
			song := spotify.Search(strings.Join(args[2:], " "))
			if song == nil {
				s.Send(`_Song "`+strings.Join(args[2:], " ")+`" Not Found._`, message.Channel)
			} else {
				name := song.Artists[0].Name + " - " + song.Name
				s.Send(`_Now Playing "`+name+`"_`, message.Channel)
			}
		case "add":
			song := spotify.Search(strings.Join(args[2:], " "))
			if song == nil {
				s.Send(`_Song "`+strings.Join(args[2:], " ")+`" Not Found._`, message.Channel)
			} else {
				name := song.Artists[0].Name + " - " + song.Name
				s.Send(`_Song "`+name+`" Was Added._`, message.Channel)
			}
		case "pause":
			s.Send(`_Pausing Spotify..._`, message.Channel)
			spotify.Pause()
		case "resume":
			s.Send(`_Resuming Spotify..._`, message.Channel)
			spotify.Resume()
		case "skip":
			s.Send(`_Skipping Song..._`, message.Channel)
			spotify.Skip()
		case "last":
			s.Send(`_Playing Previous..._`, message.Channel)
			spotify.Last()
		case "restart":
			s.Send(`_Restarting Song..._`, message.Channel)
			spotify.Restart()
		case "joke":
			s.Send(Joke{}.Get(), message.Channel)
		}
	}
}

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	spotify = Spotify{
		client:   SPOTIFY_CLIENT_AUTH,
		refresh:  SPOTIFY_REFRESH_TOKEN,
		device:   SPOTIFY_DEVICE,
		playlist: PLAYLIST_ID,
	}

	fmt.Println(spotify.Current())

	// spotify.Pause()
	// slackbot := SlackBot{}
	// slackbot.Connect("https://slack.com/api/rtm.start?token=" + SLACK_API_KEY)
	// slackbot.Subscribe(handler)
}
