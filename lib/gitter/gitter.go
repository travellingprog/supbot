// TODO: Add documentation
package gitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	RestURL   = "https://api.gitter.im/v1"
	StreamURL = "https://stream.gitter.im/v1"
)

type Message struct {
	Text     string
	Mentions []struct {
		UserId string
	}
}

type Room struct {
	Id, Name string
}

type User struct {
	Id string
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

func (g *Gitter) Initialize(done chan bool) {
	// TODO: Test this somehow
	msgCh := make(chan Message)
	errCh := make(chan error)
	defer close(msgCh)
	defer close(errCh)

	if len(g.rooms) < 1 {
		return
	}

	// TODO: Handle multiple Gitter rooms
	// - create io.Writer for each room
	// - create Hal for each room (will probably require change in hal.go)
	go g.GetRoomMsgs(g.rooms[0], msgCh, errCh, done)

	for {
		select {
		case <-done:
			return
		case msg := <-msgCh:
			// TODO: check if mentioned
			// TODO: pass to Hal
			log.Printf("msg: %+v\n", msg)
		}
	}
}

func (g *Gitter) GetRoomMsgs(room Room, msgCh chan Message, errCh chan error, done chan bool) {
	msgURL := StreamURL + "/rooms/" + room.Id + "/chatMessages"

	for {
		select {
		case <-done:
			return
		default:
		}

		var msgs []Message
		if err := get(msgURL, g.token, &msgs, "chat messages"); err != nil {
			errCh <- fmt.Errorf("Error processing chat messages: %v\n", err)
			continue
		}

		for _, msg := range msgs {
			select {
			case <-done:
				return
			case msgCh <- msg:
			}
		}
	}
}

// TODO: Make an io.Writer for Hal

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
