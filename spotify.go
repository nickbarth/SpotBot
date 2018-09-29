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

type SongJSON struct {
	Item struct {
		Artists []ArtistJSON `json:"artists"`
		Name    string       `json:"name"`
	} `json:"item"`
}

func openBrowser(url string) {
	var err error
	err = exec.Command("open", url).Start()

	if err != nil {
		log.Fatal(err)
	}
}

func request(method string, address string, header map[string]string, data map[string]string) string {
	client := &http.Client{}

	params := url.Values{}
	for key, val := range data {
		params.Add(key, val)
	}

	paramsEncoded := strings.NewReader(params.Encode())

	req, err := http.NewRequest(method, address, paramsEncoded)

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
	token  string
	client string
}

func (s *Spotify) getTokenFromRefresh(code string) TokenJSON {
	var token TokenJSON

	header := map[string]string{
		"Authorization": "Basic " + s.client,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	data := map[string]string{
		"refresh_token": code,
		"grant_type":    "refresh_token",
	}

	body := request("POST", "https://accounts.spotify.com/api/token", header, data)

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

	data := map[string]string{
		"code":         code,
		"grant_type":   "authorization_code",
		"redirect_uri": "http://localhost/",
	}

	body := request("POST", "https://accounts.spotify.com/api/token", header, data)

	fmt.Println(string(body))
	json.Unmarshal([]byte(body), &token)

	return token
}

func (s *Spotify) run(method string, endpoint string) string {
	header := map[string]string{
		"Authorization": "Bearer " + s.token,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	return request(method, endpoint, header, nil)
}

func (s *Spotify) Connect(code string) {
	token := s.getTokenFromRefresh(code)
	s.token = token.Code
}

func (s *Spotify) Play() {
	s.run("POST", "https://api.spotify.com/v1/me/player/play")
}

func (s *Spotify) Pause() {
	s.run("POST", "https://api.spotify.com/v1/me/player/pause")
}

func (s *Spotify) Skip() {
	s.run("POST", "https://api.spotify.com/v1/me/player/next")
}

func (s *Spotify) Last() {
	s.run("POST", "https://api.spotify.com/v1/me/player/previous")
}

func (s *Spotify) Volume(volume string) {
	s.run("PUT", "https://api.spotify.com/v1/me/player/volume?volume_percent="+volume)
}

func (s *Spotify) Restart() {
	s.run("PUT", "https://api.spotify.com/v1/me/player/seek?position_ms=0")
}

func (s *Spotify) Current() {
	var song SongJSON

	body := s.run("GET", "https://api.spotify.com/v1/me/player/currently-playing")
	json.Unmarshal([]byte(body), &song)

	fmt.Println(song.Item.Name, " - ", song.Item.Artists[0].Name)
}
