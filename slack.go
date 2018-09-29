package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

type StreamJSON struct {
	URL string `json:"url"`
}

type MessageJSON struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

func getWebSocket(url string) string {
	var stream StreamJSON

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &stream)
	return stream.URL
}

type SlackBot struct {
	ws *websocket.Conn
}

func (s *SlackBot) Connect(url string) {
	url = getWebSocket(url)
	ws, err := websocket.Dial(url, "", url)

	if err != nil {
		log.Fatal(err)
	}

	s.ws = ws
}

func (s *SlackBot) receive(channel chan MessageJSON) {
	for {
		var data string
		var message MessageJSON

		err := websocket.Message.Receive(s.ws, &data)

		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal([]byte(data), &message)

		if message.Type == "message" && message.Text != "" {
			channel <- message
		}
	}
}

func (s *SlackBot) Send(message string, channel string) {
	resp := new(MessageJSON)
	resp.Type = "message"
	resp.Text = message
	resp.Channel = channel

	err := websocket.JSON.Send(s.ws, resp)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *SlackBot) Subscribe(fn func(*SlackBot, MessageJSON)) {
	messages := make(chan MessageJSON)
	go s.receive(messages)

	for {
		select {
		case message := <-messages:
			fn(s, message)
		}
	}
}
