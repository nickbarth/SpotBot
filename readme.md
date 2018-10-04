# SpotBot
A Slack Spotify Bot in Go

### Commands
```
> @spotbot current
Currently Playing "Artist - Song"

> @spotbot play
Playing Spotify...

> @spotbot pause
Pausing Spotify...

> @spotbot skip
Skipping this song...

> @spotbot play
Playing last song...

> @spotbot restart
Restarting this song...

> @spotbot joke
What did the grape do when he got stepped on? He let out a little wine.
```

### Setup

```bash
# Slack App - https://api.slack.com/apps
# Spotify App - https://developer.spotify.com/dashboard/applications
SPOTIFY_CLIENT_AUTH - <CLIENT_ID>:<CLIENT_SECRET> -> Base64 -> XXXX=
https://accounts.spotify.com/en/authorize?response_type=code&client_id=XXXX&redirect_uri=http:%2F%2Flocalhost%2F&scope=user-modify-playback-state%20playlist-modify-public%20playlist-modify-private
http://localhost/#code=XXXX
curl -H "Authorization: Basic XXX=" -d grant_type=authorization_code -d code=XXX -d redirect_uri=http%3A%2F%2Flocalhost%2F https://accounts.spotify.com/api/token
```

### License
WTFPL &copy; 2018 Nick Barth
