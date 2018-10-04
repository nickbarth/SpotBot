package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

type RequestDataType int

const (
	RequestTypeString RequestDataType = iota
	RequestTypeValues
)

type RequestData struct {
	rtype  RequestDataType
	values url.Values
	text   string
}

type TokenJSON struct {
	Code    string `json:"access_token"`
	Type    string `json:"token_type"`
	Expires int    `json:"expires_in"`
	Refresh string `json:"refresh_token"`
	Scope   string `json:"scope"`
	Error   string `json:"error"`
	ErrMsg  string `json:"error_description"`
}

type ArtistJSON struct {
	Name string `json:"name"`
}

type TracksJSON struct {
	Tracks struct {
		Items []TrackJSON `json:"items"`
	} `json:"tracks"`
}

type TrackJSON struct {
	Artists []ArtistJSON `json:"artists"`
	Name    string       `json:"name"`
	ID      string       `json:"id"`
}

type SongJSON struct {
	Context struct {
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"context"`
	Item struct {
		Artists []ArtistJSON `json:"artists"`
		Name    string       `json:"name"`
		ID      string       `json:"id"`
	} `json:"item"`
}

func openBrowser(url string) {
	var err error
	err = exec.Command("open", url).Start()

	if err != nil {
		log.Fatal(err)
	}
}

func request(method string, address string, header map[string]string, data *RequestData) string {
	client := &http.Client{}
	var err error
	var req *http.Request

	if data == nil {
		req, err = http.NewRequest(method, address, nil)
	} else {
		switch data.rtype {
		case RequestTypeString:
			encodedText := strings.NewReader(data.text)
			req, err = http.NewRequest(method, address, encodedText)
		case RequestTypeValues:
			encodedValues := strings.NewReader(data.values.Encode())
			req, err = http.NewRequest(method, address, encodedValues)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	for key, val := range header {
		req.Header.Add(key, val)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

type Spotify struct {
	token    string
	client   string
	refresh  string
	playlist string
}

func (s *Spotify) getTokenFromRefresh(code string) TokenJSON {
	var token TokenJSON

	header := map[string]string{
		"Authorization": "Basic " + s.client,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	data := RequestData{
		rtype: RequestTypeValues,
		values: url.Values{
			"refresh_token": {code},
			"grant_type":    {"refresh_token"},
		},
	}

	body := request("POST", "https://accounts.spotify.com/api/token", header, &data)

	fmt.Println(string(body))
	json.Unmarshal([]byte(body), &token)

	return token
}

func (s *Spotify) getToken(code string) TokenJSON {
	var token TokenJSON

	header := map[string]string{
		"Authorization": "Basic " + s.client,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	data := RequestData{
		rtype: RequestTypeValues,
		values: url.Values{
			"code":         {code},
			"grant_type":   {"authorization_code"},
			"redirect_uri": {"http://localhost/"},
		},
	}

	body := request("POST", "https://accounts.spotify.com/api/token", header, &data)

	fmt.Println(string(body))
	json.Unmarshal([]byte(body), &token)

	return token
}

func (s *Spotify) run(method string, endpoint string, data *RequestData) string {
	header := map[string]string{
		"Authorization": "Bearer " + s.token,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	if data != nil && data.rtype == RequestTypeString {
		header["Content-Type"] = "application/json"
	}

	return request(method, endpoint, header, data)
}

func (s *Spotify) Connect(playlist string) {
	token := s.getTokenFromRefresh(s.refresh)
	s.token = token.Code
	s.playlist = playlist
}

func (s *Spotify) Search(term string) *TrackJSON {
	var track TracksJSON
	p := url.Values{"type": {"track"}, "q": {term}}
	body := s.run("GET", "https://api.spotify.com/v1/search?"+p.Encode(), nil)
	json.Unmarshal([]byte(body), &track)

	if len(track.Tracks.Items) == 0 {
		return nil
	}

	return &track.Tracks.Items[0]
}

func (s *Spotify) Play(songID string) {
	data := RequestData{
		rtype: RequestTypeString,
		text:  "{\"context_uri\":\"spotify:album:5ht7ItJgpBH7W6vJ5BqpPr\",\"offset\":{\"position\":5},\"position_ms\":0}",
	}

	a := s.run("PUT", "https://api.spotify.com/v1/me/player/play", &data)
	fmt.Println(a)
	fmt.Println(s.playlist)
}

func (s *Spotify) Pause() {
	s.run("POST", "https://api.spotify.com/v1/me/player/pause", nil)
}

func (s *Spotify) Resume() {
	s.run("POST", "https://api.spotify.com/v1/me/player/play", nil)
}

func (s *Spotify) Skip() {
	s.run("POST", "https://api.spotify.com/v1/me/player/next", nil)
}

func (s *Spotify) Last() {
	s.run("POST", "https://api.spotify.com/v1/me/player/previous", nil)
}

func (s *Spotify) Volume(volume string) {
	s.run("PUT", "https://api.spotify.com/v1/me/player/volume?volume_percent="+volume, nil)
}

func (s *Spotify) Restart() {
	s.run("PUT", "https://api.spotify.com/v1/me/player/seek?position_ms=0", nil)
}

func (s *Spotify) Current() string {
	var song SongJSON

	body := s.run("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	json.Unmarshal([]byte(body), &song)

	return song.Item.Artists[0].Name + " - " + song.Item.Name
}
