package gitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	RestURL   = "https://api.gitter.im/v1/"
	StreamURL = "https://stream.gitter.im/v1/"
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

func NewGitter(token string) (*Gitter, error) {
	// USERID
	req, err := http.NewRequest("GET", RestURL+"user", nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create GET request for Gitter bot user: %v", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve Gitter bot user: %v", err)
	}
	defer resp.Body.Close()

	var users []User
	if err = json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("Could not decode Gitter bot user: %v", err)
	}
	if len(users) < 1 {
		return nil, fmt.Errorf("Gitter user data is empty")
	}
	user := users[0]
	log.Printf("User: %+v\n", user)

	// ROOMS
	req, err = http.NewRequest("GET", RestURL+"rooms", nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create GET request for Gitter rooms: %v", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve Gitter rooms: %v", err)
	}
	defer resp2.Body.Close()

	var rooms []Room
	if err = json.NewDecoder(resp2.Body).Decode(&rooms); err != nil {
		return nil, fmt.Errorf("Could not decode room data: %v", err)
	}

	return &Gitter{token: token, user: user, rooms: rooms}, nil
}

func (g *Gitter) Initialize() {
	for _, room := range g.rooms {
		go func(room Room) error {
			for {
				req, err := http.NewRequest("GET", StreamURL+room.Id+"/chatMessages", nil)
				if err != nil {
					return fmt.Errorf("Could not create GET request for Gitter msgs: %v", err)
				}
				req.Header.Add("Accept", "application/json")
				req.Header.Add("Authorization", "Bearer "+g.token)

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
