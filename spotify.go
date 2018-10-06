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
	"time"
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

type UserJSON struct {
	ID   string `json:"id"`
	Name string `json:"display_name"`
}

type TrackJSON struct {
	Artists []ArtistJSON `json:"artists"`
	Name    string       `json:"name"`
	ID      string       `json:"id"`
	URI     string       `json:"uri"`
	User    UserJSON     `json:"added_by"`
}

func (t TrackJSON) Title() string {
	return t.Artists[0].Name + " - " + t.Name
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
		User  UserJSON  `json:"added_by"`
	} `json:"items"`
}

type ErrorJSON struct {
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}

type apiError struct {
	status int
	prob   string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%d - %s", e.status, e.prob)
}

func openBrowser(url string) {
	var err error
	err = exec.Command("open", url).Start()

	if err != nil {
		log.Fatal(err)
	}
}

func request(method string, address string, header http.Header, data *RequestData) string {
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

	req.Header = header

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
	timeout  time.Time
}

func (s *Spotify) getTokenFromRefresh(code string) TokenJSON {
	var token TokenJSON

	header := http.Header{
		"Authorization": {"Basic " + s.client},
		"Content-Type":  {"application/x-www-form-urlencoded"},
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

	header := http.Header{
		"Authorization": {"Basic " + s.client},
		"Content-Type":  {"application/x-www-form-urlencoded"},
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

func (s *Spotify) run(method string, endpoint string, data *RequestData) (string, error) {
	if s.token == "" {
		return "", &apiError{401, "Authorization Header Required"}
	}

	if time.Now().After(s.timeout) {
		s.Connect()
	}

	header := http.Header{
		"Authorization": {"Bearer " + s.token},
		"Content-Type":  {"application/x-www-form-urlencoded"},
	}

	if data != nil && data.rtype == RequestTypeString {
		header.Set("Content-Type", "application/json")
	}

	var spotErr ErrorJSON
	body := request(method, endpoint, header, data)
	err := json.Unmarshal([]byte(body), &spotErr)

	if err == nil && spotErr.Error.Status != 0 {
		return "", &apiError{spotErr.Error.Status, spotErr.Error.Message}
	}

	return body, nil
}

func (s *Spotify) Connect() {
	token := s.getTokenFromRefresh(s.refresh)
	s.token = token.Code
	s.timeout = time.Now().Add(time.Second * time.Duration(token.Expires))
}

func (s *Spotify) Search(term string) (*TrackJSON, error) {
	var search SearchJSON

	p := url.Values{"type": {"track"}, "q": {term}}
	body, err := s.run("GET", "https://api.spotify.com/v1/search?"+p.Encode(), nil)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &search)

	if err != nil {
		return nil, err
	}

	if len(search.Search.Tracks) == 0 {
		return nil, err
	}

	return &search.Search.Tracks[0], nil
}

func (s *Spotify) PlayAdd(uid string) error {
	err := s.AddUnique(uid)

	if err != nil {
		return err
	}

	index, err := s.Index(uid)

	if err != nil {
		return err
	}

	data := RequestData{
		rtype: RequestTypeString,
		text:  `{"context_uri":"spotify:playlist:` + s.playlist + `","offset":{"position":` + strconv.Itoa(index) + `},"position_ms":0}`,
	}

	query := ""
	if s.device != "" {
		query = "device_id=" + s.device
	}

	_, err = s.run("PUT", "https://api.spotify.com/v1/me/player/play?"+query, &data)
	return err
}

func (s *Spotify) PlaySong(uri string) error {
	data := RequestData{
		rtype: RequestTypeString,
		text:  `{"uris":["` + uri + `"],"offset":{"position":0},"position_ms":0}`,
	}

	_, err := s.run("PUT", "https://api.spotify.com/v1/me/player/play", &data)
	return err
}

func (s *Spotify) Pause() error {
	_, err := s.run("PUT", "https://api.spotify.com/v1/me/player/pause", nil)
	return err
}

func (s *Spotify) Resume() error {
	_, err := s.run("PUT", "https://api.spotify.com/v1/me/player/play", nil)
	return err
}

func (s *Spotify) Skip() error {
	_, err := s.run("POST", "https://api.spotify.com/v1/me/player/next", nil)
	return err
}

func (s *Spotify) User(id string) (*UserJSON, error) {
	var user UserJSON
	body, err := s.run("GET", "https://api.spotify.com/v1/users/"+id, nil)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Spotify) Last() error {
	_, err := s.run("POST", "https://api.spotify.com/v1/me/player/previous", nil)
	return err
}

func (s *Spotify) Add(uri string) error {
	p := url.Values{"uris": {uri}, "position": {"0"}}
	_, err := s.run("POST", "https://api.spotify.com/v1/playlists/"+s.playlist+"/tracks?"+p.Encode(), nil)
	return err
}

func (s *Spotify) Remove(uri string) error {
	data := RequestData{
		rtype: RequestTypeString,
		text:  `{"tracks":[{"uri":"` + uri + `"}]}`,
	}
	_, err := s.run("DELETE", "https://api.spotify.com/v1/playlists/"+s.playlist+"/tracks", &data)
	return err
}

func (s *Spotify) Tracks() ([]TrackJSON, error) {
	var playlist PlaylistJSON
	tracks := []TrackJSON{}

	p := url.Values{"fields": {"total,next,items(added_by,track(id,name,uri))"}}
	body, _ := s.run("GET", "https://api.spotify.com/v1/playlists/"+s.playlist+"/tracks?"+p.Encode(), nil)

	fmt.Println(body)

	err := json.Unmarshal([]byte(body), &playlist)

	if err != nil {
		return nil, err
	}

	for _, item := range playlist.Items {
		item.Track.User = item.User
		// fmt.Println(item.Track.User, item.User)
		tracks = append(tracks, item.Track)
	}

	return tracks, nil
}

func (s *Spotify) Contains(uri string) (bool, error) {
	tracks, err := s.Tracks()

	if err != nil {
		return false, err
	}

	for _, track := range tracks {
		if track.URI == uri {
			return true, nil
		}
	}

	return false, nil
}

func (s *Spotify) Index(uri string) (int, error) {
	tracks, err := s.Tracks()

	if err != nil {
		return -1, err
	}

	for index, track := range tracks {
		fmt.Println(uri, "|", track.Name)
		if track.Name == uri {
			return index, nil
		}
	}

	fmt.Println("not found")
	return -1, nil
}

func (s *Spotify) Blame(uri string) (*UserJSON, error) {
	tracks, err := s.Tracks()

	if err != nil {
		return nil, err
	}

	for _, track := range tracks {
		// fmt.Println(track.User)
		if track.Title() == uri {
			fmt.Println("matched")
			user, err := s.User(track.User.ID)
			//fmt.Println(track.User.ID)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
	}

	return nil, nil
}

func (s *Spotify) AddUnique(uri string) error {
	added, err := s.Contains(uri)

	if err != nil {
		return err
	}

	if !added {
		err = s.Add(uri)
		return err
	}

	return nil
}

func (s *Spotify) Volume(volume string) error {
	_, err := s.run("PUT", "https://api.spotify.com/v1/me/player/volume?volume_percent="+volume, nil)
	return err
}

func (s *Spotify) Restart() error {
	_, err := s.run("PUT", "https://api.spotify.com/v1/me/player/seek?position_ms=0", nil)
	return err
}

func (s *Spotify) Devices() ([]DeviceJSON, error) {
	var devices DevicesJSON
	body, err := s.run("GET", "https://api.spotify.com/v1/me/player/devices", nil)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &devices)

	if err != nil {
		return nil, err
	}

	return devices.Items, nil
}

func (s *Spotify) Current() (*TrackJSON, error) {
	var song CurrentJSON

	body, err := s.run("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)

	if err != nil {
		return nil, err
	}

	if body == "" {
		return &TrackJSON{Name: "Nothing", Artists: []ArtistJSON{
			ArtistJSON{Name: "No One"},
		}}, nil
	}

	err = json.Unmarshal([]byte(body), &song)

	if err != nil {
		return nil, err
	}

	return &song.Track, nil
}
