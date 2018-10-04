package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
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

type SearchJSON struct {
	Search struct {
		Tracks []TrackJSON `json:"items"`
	} `json:"tracks"`
}

type ArtistJSON struct {
	Name string `json:"name"`
}

type ContextJSON struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type TrackJSON struct {
	Artists []ArtistJSON `json:"artists"`
	Name    string       `json:"name"`
	ID      string       `json:"id"`
	URI     string       `json:"uri"`
}

type DeviceJSON struct {
	ID     string `json:"id"`
	Active bool   `json:"is_active"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type DevicesJSON struct {
	Items []DeviceJSON `json:"devices"`
}

type CurrentJSON struct {
	Context ContextJSON `json:"context"`
	Track   TrackJSON   `json:"item"`
}

type PlaylistJSON struct {
	Items []struct {
		Track TrackJSON `json:"track"`
	} `json:"items"`
}

type ErrorJSON struct {
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
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

	switch {
	case data == nil:
		req, err = http.NewRequest(method, address, nil)
	case data.rtype == RequestTypeString:
		encodedText := strings.NewReader(data.text)
		req, err = http.NewRequest(method, address, encodedText)
	case data.rtype == RequestTypeValues:
		encodedValues := strings.NewReader(data.values.Encode())
		req, err = http.NewRequest(method, address, encodedValues)
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
	device   string
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

	var spotErr ErrorJSON
	body := request(method, endpoint, header, data)
	err := json.Unmarshal([]byte(body), &spotErr)

	if err == nil && spotErr.Error.Status != 0 {
		log.Fatal(spotErr.Error.Status, " - ", spotErr.Error.Message)
	}

	return body
}

func (s *Spotify) Connect() {
	token := s.getTokenFromRefresh(s.refresh)
	s.token = token.Code
}

func (s *Spotify) Search(term string) *TrackJSON {
	var search SearchJSON

	p := url.Values{"type": {"track"}, "q": {term}}
	body := s.run("GET", "https://api.spotify.com/v1/search?"+p.Encode(), nil)
	json.Unmarshal([]byte(body), &search)

	if len(search.Search.Tracks) == 0 {
		return nil
	}

	return &search.Search.Tracks[0]
}

func (s *Spotify) Play(uid string) {
	s.AddUnique(uid)
	index := strconv.Itoa(s.Index(uid))

	data := RequestData{
		rtype: RequestTypeString,
		text:  `{"context_uri":"spotify:playlist:` + s.playlist + `","offset":{"position":` + index + `},"position_ms":0}`,
	}

	query := ""
	if s.device != "" {
		query = "device_id=" + s.device
	}

	s.run("PUT", "https://api.spotify.com/v1/me/player/play?"+query, &data)
}

func (s *Spotify) PlaySong(uri string) {
	data := RequestData{
		rtype: RequestTypeString,
		text:  `{"uris":["` + uri + `"],"offset":{"position":0},"position_ms":0}`,
	}

	s.run("PUT", "https://api.spotify.com/v1/me/player/play", &data)
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

func (s *Spotify) Add(uri string) {
	p := url.Values{"uris": {uri}, "position": {"0"}}
	s.run("POST", "https://api.spotify.com/v1/playlists/"+s.playlist+"/tracks?"+p.Encode(), nil)
}

func (s *Spotify) Tracks() []TrackJSON {
	var playlist PlaylistJSON
	tracks := []TrackJSON{}

	p := url.Values{"fields": {"items(track(id,name,uri))"}}
	body := s.run("GET", "https://api.spotify.com/v1/playlists/"+s.playlist+"/tracks?"+p.Encode(), nil)
	err := json.Unmarshal([]byte(body), &playlist)

	if err != nil {
		log.Fatal(err)
	}

	for _, item := range playlist.Items {
		tracks = append(tracks, item.Track)
	}

	return tracks
}

func (s *Spotify) Contains(uri string) bool {
	for _, track := range s.Tracks() {
		if track.URI == uri {
			return true
		}
	}

	return false
}

func (s *Spotify) Index(uri string) int {
	for index, track := range s.Tracks() {
		if track.URI == uri {
			return index
		}
	}

	return -1
}

func (s *Spotify) AddUnique(uri string) {
	if !s.Contains(uri) {
		s.Add(uri)
	}
}

func (s *Spotify) Volume(volume string) {
	s.run("PUT", "https://api.spotify.com/v1/me/player/volume?volume_percent="+volume, nil)
}

func (s *Spotify) Restart() {
	s.run("PUT", "https://api.spotify.com/v1/me/player/seek?position_ms=0", nil)
}

func (s *Spotify) Devices() []DeviceJSON {
	var devices DevicesJSON
	body := s.run("GET", "https://api.spotify.com/v1/me/player/devices", nil)
	err := json.Unmarshal([]byte(body), &devices)

	if err != nil {
		log.Fatal(err)
	}

	return devices.Items
}

func (s *Spotify) Current() string {
	var song CurrentJSON

	body := s.run("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)

	if body == "" {
		return "Nothing - No One"
	}

	err := json.Unmarshal([]byte(body), &song)

	if err != nil {
		log.Fatal(err)
	}

	return song.Track.Artists[0].Name + " - " + song.Track.Name
}
