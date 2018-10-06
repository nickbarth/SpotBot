package main

import (
	"fmt"
	"log"
)

const SLACK_API_KEY = ""
const SPOTIFY_CLIENT_AUTH = ""
const SPOTIFY_REFRESH_TOKEN = ""
const PLAYLIST_ID = ""
const SPOTIFY_DEVICE = ""

var spotify Spotify

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	spotify = Spotify{
		client:   SPOTIFY_CLIENT_AUTH,
		refresh:  SPOTIFY_REFRESH_TOKEN,
		device:   SPOTIFY_DEVICE,
		playlist: PLAYLIST_ID,
	}

	spotify.Connect()

	slackbot := SlackBot{}
	slackbot.Connect("https://slack.com/api/rtm.start?token=" + SLACK_API_KEY)

	slackbot.Command("default", func(args string) string {
		return `_I'm sorry. I'm afraid I can't do that._`
	})

	slackbot.Command("current", func(args string) string {
		current, err := spotify.Current()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Currently Playing "` + current.Title() + `"_`
	})

	slackbot.Command("play", func(name string) string {
		if name == "" {
			err := spotify.Resume()
			if err != nil {
				return `_Error ` + err.Error() + `_`
			}
			return `_Resuming Spotify..._`
		}

		song, err := spotify.Search(name)

		if err != nil {
			return `_Error ` + err.Error() + `_`
		}

		if song == nil {
			return `_Song "` + name + `" Not Found._`
		} else {
			err = spotify.PlayAdd(song.URI)

			if err != nil {
				return `_Error ` + err.Error() + `_`
			}

			return `_Now Playing "` + song.Title() + `"_`
		}
	})

	slackbot.Command("add", func(name string) string {
		if name == "" {
			return `_No Song Specified._`
		}

		song, err := spotify.Search(name)

		if err != nil {
			return `_Error ` + err.Error() + `_`
		} else if song == nil {
			return `_Song "` + name + `" Not Found._`
		} else {
			err = spotify.AddUnique(song.URI)

			if err != nil {
				return `_Error ` + err.Error() + `_`
			}

			return `_Song "` + song.Title() + `" Was Added._`
		}
	})

	slackbot.Command("remove", func(name string) string {
		if name == "" {
			return `_No Song Specified._`
		}

		song, err := spotify.Search(name)

		if err != nil {
			return `_Error ` + err.Error() + `_`
		} else if song == nil {
			return `_Song "` + name + `" Not Found._`
		} else {
			err = spotify.Remove(song.URI)

			if err != nil {
				return `_Error ` + err.Error() + `_`
			}

			current, err := spotify.Current()
			if err != nil {
				return `_Error ` + err.Error() + `_`
			}

			if current.ID == song.ID {
				err := spotify.Skip()
				if err != nil {
					return `_Error ` + err.Error() + `_`
				}
			}

			return `_Song "` + song.Title() + `" Was Removed._`
		}
	})

	slackbot.Command("pause", func(args string) string {
		err := spotify.Pause()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Pausing Spotify..._`
	})

	slackbot.Command("resume", func(args string) string {
		err := spotify.Resume()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Resuming Spotify..._`
	})

	slackbot.Command("next", func(args string) string {
		err := spotify.Skip()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Skipping Song..._`
	})

	slackbot.Command("last", func(args string) string {
		err := spotify.Last()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Playing Previous..._`
	})

	slackbot.Command("restart", func(args string) string {
		err := spotify.Restart()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Restarting Song..._`
	})

	slackbot.Command("joke", func(args string) string {
		return Joke{}.Get()
	})

	fmt.Println("Listening...")
	slackbot.Listen()
}
