package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type StreamJSON struct {
	URL  string `json:"url"`
	Self struct {
		ID string `json:"id"`
	} `json:"self"`
}

type MessageJSON struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

type SlackBot struct {
	user     string
	ws       *websocket.Conn
	commands map[string]func(string) string
}

func (s *SlackBot) Connect(url string) {
	var stream StreamJSON

	s.commands = map[string]func(string) string{}

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

	ws, err := websocket.Dial(stream.URL, "", stream.URL)

	if err != nil {
		log.Fatal(err)
	}

	s.ws = ws
	s.user = "<@" + stream.Self.ID + ">"
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

func (s *SlackBot) Command(command string, fn func(string) string) {
	s.commands[command] = fn
}

func (s *SlackBot) Listen() {
	messages := make(chan MessageJSON)
	go s.receive(messages)

	for {
		select {
		case message := <-messages:
			if strings.HasPrefix(message.Text, s.user) {
				fmt.Println("Message Received:", message.Text)

				args := strings.Split(message.Text, " ")

				if len(args) < 2 {
					s.Send(s.commands["default"](""), message.Channel)
					return
				}

				command := args[1]

				if handler, ok := s.commands[command]; ok {
					s.Send(handler(strings.Join(args[2:], " ")), message.Channel)
				} else {
					s.Send(s.commands["default"](""), message.Channel)
				}
			}
		}
	}
}
