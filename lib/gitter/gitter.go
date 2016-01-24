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

func (g *Gitter) Initialize() {
	for _, room := range g.rooms {
		go func(room Room) error {
			for {
				req, err := http.NewRequest("GET", StreamURL+"/rooms/"+room.Id+"/chatMessages", nil)
				if err != nil {
					return fmt.Errorf("Could not create GET request for Gitter msgs: %v", err)
				}
				req.Header.Add("Accept", "application/json")
				req.Header.Add("Authorization", "Bearer "+g.token)

				// Long-polling here
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return fmt.Errorf("Could not retrieve Gitter messages: %v", err)
				}
				defer resp.Body.Close()

				var msgs []Message
				if err = json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
					return fmt.Errorf("Could not decode Gitter messages data: %v", err)
				}

				log.Printf("msgs: %+v\n", msgs)
			}
		}(room)
	}
}

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
