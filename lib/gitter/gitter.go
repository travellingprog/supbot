// TODO: Add documentation
// TODO: Be able to handle multiple rooms
package gitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gophergala2016/supbot/lib/hal"
)

var (
	RestURL   = "https://api.gitter.im/v1"
	StreamURL = "https://stream.gitter.im/v1"
)

type Message struct {
	Text     string `json:"text"`
	Mentions []struct {
		UserId string
	} `json:"mentions"`
}

type Room struct {
	Id, Name string
}

type User struct {
	Id, Username string
}

type Gitter struct {
	token string
	user  User
	rooms []Room
}

func NewGitter(token string) (g *Gitter, err error) {
	var users []User
	if err = get(RestURL+"/user", token, &users, "current user"); err != nil {
		return nil, err
	}
	if len(users) < 1 {
		err = fmt.Errorf("Gitter user data is empty")
		return nil, err
	}
	user := users[0]

	var rooms []Room
	if err = get(RestURL+"/rooms", token, &rooms, "rooms"); err != nil {
		return nil, err
	}

	g = &Gitter{token: token, user: user, rooms: rooms}
	return g, nil
}

// Start begins fetching messages for the Gitter room, and outputs them to the console
func (g *Gitter) Start() {
	msgCh := make(chan Message)
	errCh := make(chan error)

	if len(g.rooms) < 1 {
		return
	}

	supBot := hal.New(g)

	go func() {
		for {
			g.processMsgs(supBot, msgCh)
		}
	}()
	go func() {
		for {
			log.Println(<-errCh)
		}
	}()

	log.Println("Listening for messages on Gitter...")
	for {
		g.getRoomMsgs(g.rooms[0], msgCh, errCh)
	}
}

// Write is given the output from sup and writes it to the chat room.
func (g *Gitter) Write(o []byte) (n int, err error) {
	url := RestURL + "/rooms/" + g.rooms[0].Id + "/chatMessages"

	body := new(bytes.Buffer)
	newMsg := Message{Text: string(o)}
	if err := json.NewEncoder(body).Encode(newMsg); err != nil {
		return 0, fmt.Errorf("Could not encode supbot output: %v", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return 0, fmt.Errorf("Could not create POST request to Gitter: %v", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+g.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Could not POST Gitter message: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Gitter POST request returned status %d", resp.StatusCode)
	}
	return len(o), nil
}

func (g *Gitter) getRoomMsgs(room Room, msgCh chan Message, errCh chan error) {
	msgURL := StreamURL + "/rooms/" + room.Id + "/chatMessages"
	var msg Message
	if err := get(msgURL, g.token, &msg, "chat messages"); err != nil {
		errCh <- err
		return
	}
	msgCh <- msg
}

// processMsgs takes one Message from msgCh that mentions the gitter bot, and
// sends it to our instance of Hal.
func (g *Gitter) processMsgs(supBot io.Writer, msgCh chan Message) {
	msg := <-msgCh
	if g.wasMentioned(msg) {
		msg.Text = strings.Replace(msg.Text, "@"+g.user.Username, "", -1)
		msg.Text = strings.TrimSpace(msg.Text)
		supBot.Write([]byte(msg.Text))
	}
}

// wasMentioned checks if our gitter bot was mentioned in the message
func (g *Gitter) wasMentioned(msg Message) bool {
	for _, mention := range msg.Mentions {
		if mention.UserId == g.user.Id {
			return true
		}
	}
	return false
}

func get(path string, token string, target interface{}, descr string) error {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return fmt.Errorf("Could not create GET request for Gitter %s: %v", descr, err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not GET Gitter %s: %v", descr, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Gitter GET request for %s returned status %d", descr, resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("Could not decode Gitter %s: %v", descr, err)
	}
	return nil
}
