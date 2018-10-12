package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const SLACK_API_KEY = ""
const SPOTIFY_CLIENT_AUTH = ""
const SPOTIFY_REFRESH_TOKEN = ""
const PLAYLIST_ID = ""
const SPOTIFY_DEVICE = ""

var spotify Spotify

func main() {
	rand.Seed(time.Now().Unix())
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

	hello := func(args string) string {
		hour := time.Now().Hour()
		switch {
		case hour < 12:
			return `_Good Morning!_`
		case hour == 12:
			return `_Happy Lunch Time!_`
		case hour > 12 && hour < 17:
			return `_Good Afternoon!_`
		default:
			return `_Good Evening!_`
		}
	}
	slackbot.Command("hi", hello)
	slackbot.Command("hello", hello)

	bye := func(args string) string {
		return `_Goodbye!_`
	}
	slackbot.Command("bye", bye)
	slackbot.Command("goodbye", bye)

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
			err = spotify.PlaylistPlay(song.URI)

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

	slackbot.Command("blame", func(name string) string {
		var song *TrackJSON
		var err error

		if name == "" {
			song, err = spotify.Current()
		} else {
			song, err = spotify.Search(name)
		}

		if err != nil {
			return `_Error ` + err.Error() + `_`
		}

		if song == nil {
			return `_Song "` + name + `" Not Found._`
		}

		user, err := spotify.Blame(song.URI)

		if err != nil {
			return `_Error ` + err.Error() + `_`
		}

		if user.Name == "" {
			user.Name = "uknown"
		}

		return `_"` + song.Title() + `" was added by ` + user.Name + `._`
	})

	slackbot.Command("remove", func(name string) string {
		var song *TrackJSON
		var err error

		if name == "" {
			song, err = spotify.Current()
		} else {
			song, err = spotify.Search(name)
		}

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

	pause := func(args string) string {
		err := spotify.Pause()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Pausing Spotify..._`
	}
	slackbot.Command("pause", pause)
	slackbot.Command("stop", pause)

	slackbot.Command("resume", func(args string) string {
		err := spotify.Resume()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Resuming Spotify..._`
	})

	slackbot.Command("setup", func(args string) string {
		err := spotify.Shuffle()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}

		err = spotify.Repeat()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}

		return `_Shuffle and Repeat Enabled._`
	})

	slackbot.Command("shuffle", func(args string) string {
		err := spotify.ShufflePlaylist()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}

		return `_I shuffled your playlist._`
	})

	next := func(args string) string {
		err := spotify.Skip()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Skipping Song..._`
	}
	slackbot.Command("next", next)
	slackbot.Command("skip", next)

	previous := func(args string) string {
		err := spotify.Last()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Playing Previous..._`
	}
	slackbot.Command("previous", previous)
	slackbot.Command("last", previous)

	slackbot.Command("restart", func(args string) string {
		err := spotify.Restart()
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Restarting Song..._`
	})

	slackbot.Command("volume", func(args string) string {
		err := spotify.Volume(args)
		if err != nil {
			return `_Error ` + err.Error() + `_`
		}
		return `_Setting Volume..._`
	})

	slackbot.Command("joke", func(args string) string {
		return Joke{}.Get()
	})

	fmt.Println("Listening...")
	slackbot.Listen()
}
