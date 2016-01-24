package gitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// supgitter@gmail.com
// suppressly

// https://github.com/supgitter
// suppressly1

// Gitter token: 4b409f3d662592192095055ac603eaf106b0b92b

var (
	RestURL   = "https://api.gitter.im/v1/"
	StreamURL = "https://stream.gitter.im/v1/"
)

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
		go func(room Room) {
			for {
				// req, err = http.NewRequest("GET", RestURL + "rooms", nil)
				// if err != nil {
				// 	return nil, fmt.Errorf("Could not create GET request for Gitter rooms: %v", err)
				// }
				// req.Header.Add("Accept", "application/json")
				// req.Header.Add("Authorization", "Bearer "+token)

				resp, err := http.Get(StreamURL + room.Id + "/chatMessages")
				if err != nil {
					log.Fatalf("Error getting message, room %s: %v", room.Name, err)
				}
				log.Println("%+v", resp)

			}
		}(room)
	}
}

// {
//     "id": "56a42b486b6468374a0926a4",
//     "text": "sup man",
//     "html": "sup man",
//     "sent": "2016-01-24T01:39:20.224Z",
//     "fromUser": {
//         "id": "56a3f554e610378809bddc9c",
//         "username": "supgitter",
//         "displayName": "supgitter",
//         "url": "/supgitter",
//         "avatarUrlSmall": "https://avatars0.githubusercontent.com/u/16857436?v=3&s=60",
//         "avatarUrlMedium": "https://avatars0.githubusercontent.com/u/16857436?v=3&s=128",
//         "v": 1,
//         "gv": "3"
//     },
//     "unread": true,
//     "readBy": 0,
//     "urls": [],
//     "mentions": [],
//     "issues": [],
//     "meta": [],
//     "v": 1
// }
