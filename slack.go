package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	id       int
	ws       *websocket.Conn
	commands map[string]func(string) string
}

type Ping struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Time int32  `json:"time"`
}

type Pong struct {
	Reply int    `json:"reply_to"`
	Type  string `json:"type"`
	Time  int32  `json:"time"`
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
	s.id, _ = strconv.Atoi(stream.Self.ID)
}

func (s *SlackBot) ping() {
	ping := new(Ping)
	ping.Id = s.id
	ping.Type = "ping"

	for {
		ping.Time = int32(time.Now().Unix())
		err := websocket.JSON.Send(s.ws, ping)

		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(60 * time.Second)
	}
}

func (s *SlackBot) receive(channel chan MessageJSON) {
	for {
		var data string
		var message MessageJSON

		err := websocket.Message.Receive(s.ws, &data)

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(data), &message)

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
	go s.ping()

	messages := make(chan MessageJSON)
	go s.receive(messages)

	for {
		select {
		case message := <-messages:
			if strings.HasPrefix(message.Text, s.user) {

				if message.Type != "pong" {
					fmt.Printf("Message Received: %+v\n", message)
				}

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
