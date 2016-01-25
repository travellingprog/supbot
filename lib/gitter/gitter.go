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
	if err = Get(RestURL+"/user", token, &users, "current user"); err != nil {
		return
	}
	if len(users) < 1 {
		err = fmt.Errorf("Gitter user data is empty")
		return
	}
	user := users[0]

	var rooms []Room
	if err = Get(RestURL+"/rooms", token, &rooms, "rooms"); err != nil {
		return
	}

	g = &Gitter{token: token, user: user, rooms: rooms}
	return
}

func (g *Gitter) Initialize(done chan struct{}) {
	// TODO: Test this somehow
	out := make(chan Message)
	defer close(out)

	if len(g.rooms) < 1 {
		// TODO: return error
	}
	// TODO: Handle multiple Gitter rooms
	// - create io.Writer for each room
	// - create Hal for each room (will probably require change in hal.go)
	//
	// for _, room := range g.rooms {
	// 	go g.GetRoomMsgs(room, out, done)
	// }
	go g.GetRoomMsgs(g.rooms[0], out, done)

	for {
		select {
		case msg := <-out:
			// TODO: check if mentioned
			// TODO: pass to Hal
			log.Printf("msg: %+v\n", msg)
		case <-done:
			return
		}
	}
}

func (g *Gitter) GetRoomMsgs(room Room, out chan Message, done chan struct{}) {
	// TODO: Test this somehow
	msgURL := StreamURL + "/rooms/" + roomId + "/chatMessages"

	for {
		var msgs []Message
		if err := Get(msgURL, g.token, &msgs, "chat messages"); err != nil {
			log.Printf("Error processing chat messages: %v\n", err)
		}

		// TODO: if len(msgs) == 0, or there was an issue with msgs,
		//       done is not checked. Need to fix that
		for _, msg := range msgs {
			select {
			case <-done:
				return
			case out <- msg:
			}
		}
	}
}

// TODO: Make an io.Writer for Hal

func Get(path string, token string, target interface{}, descr string) (err error) {
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

	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("Could not decode Gitter %s: %v", descr, err)
	}
	return
}
