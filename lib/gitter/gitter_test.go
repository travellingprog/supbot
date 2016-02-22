package gitter

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	Token  = "SOME_TOKEN"
	RoomId = "SOME_ROOM"
)

func TestMain(m *testing.M) {
	server := setMockServer()
	defer server.Close()
	RestURL = server.URL
	StreamURL = server.URL
	os.Exit(m.Run())
}

func setMockServer() *httptest.Server {
	// Responses taken from real calls to Gitter API
	userResponse := []byte(`[{"id":"56a3f554e610378809bddc9c","username":"supgitter","displayName":"supgitter","url":"/supgitter","avatarUrlSmall":"https://avatars0.githubusercontent.com/u/16857436?v=3&s=60","avatarUrlMedium":"https://avatars0.githubusercontent.com/u/16857436?v=3&s=128","v":1,"gv":"3"}]`)
	roomsResponse := []byte(`[{"id":"SOME_ROOM","name":"travellingprog/supbot","topic":"","uri":"travellingprog/supbot","oneToOne":false,"userCount":2,"unreadItems":0,"mentions":0,"lastAccessTime":"2016-01-24T15:51:30.513Z","lurk":false,"activity":true,"url":"/travellingprog/supbot","githubType":"REPO","security":"PUBLIC","noindex":false,"tags":[],"roomMember":true}]`)
	msgsResponse := []byte(`[{"id":"56a4306f8fbaf4220af8f817","text":"@travellingprog what's good?","html":"<span data-link-type=\"mention\" data-screen-name=\"travellingprog\" class=\"mention\">@travellingprog</span> what&#39;s good?","sent":"2016-01-24T02:01:19.407Z","fromUser":{"id":"56a3f554e610378809bddc9c","username":"supgitter","displayName":"supgitter","url":"/supgitter","avatarUrlSmall":"https://avatars0.githubusercontent.com/u/16857436?v=3&s=60","avatarUrlMedium":"https://avatars0.githubusercontent.com/u/16857436?v=3&s=128","v":1,"gv":"3"},"unread":false,"readBy":1,"urls":[],"mentions":[{"screenName":"travellingprog","userId":"56a3eab0e610378809bddb7d","userIds":[]}],"issues":[],"meta":[],"v":1},{"id":"56a50f97eaf741c118d49b13","text":"@supgitter local ping","html":"<span data-link-type=\"mention\" data-screen-name=\"supgitter\" class=\"mention\">@supgitter</span> local ping","sent":"2016-01-24T17:53:27.758Z","fromUser":{"id":"56a3eab0e610378809bddb7d","username":"travellingprog","displayName":"Erick Cardenas-Mendez","url":"/travellingprog","avatarUrlSmall":"https://avatars2.githubusercontent.com/u/3519160?v=3&s=60","avatarUrlMedium":"https://avatars2.githubusercontent.com/u/3519160?v=3&s=128","gv":"3"},"unread":false,"readBy":1,"urls":[],"mentions":[{"screenName":"supgitter","userId":"56a3f554e610378809bddc9c","userIds":[]}],"issues":[],"meta":[],"v":1}]`)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/user":
			w.Write(userResponse)
		case "/rooms":
			w.Write(roomsResponse)
		case "/rooms/" + RoomId + "/chatMessages":
			w.Write(msgsResponse)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestGet(t *testing.T) {
	var err error

	var users []User
	if err = get(RestURL+"/user", Token, &users, "current user"); err != nil {
		t.Error(err)
	}
	t.Logf("users: %+v\n", users)

	var rooms []Room
	if err = get(RestURL+"/rooms", Token, &rooms, "rooms"); err != nil {
		t.Error(err)
	}
	t.Logf("rooms: %+v\n", rooms)

	var msgs []Message
	if err = get(StreamURL+"/rooms/"+RoomId+"/chatMessages", Token, &msgs, "chat messages"); err != nil {
		t.Error(err)
	}
	t.Logf("msgs: %+v\n", msgs)
}

func TestNewGitter(t *testing.T) {
	gitter, err := NewGitter(Token)
	if err != nil {
		t.Error(err)
	}
	t.Logf("gitter: %+v\n", gitter)
}

func TestStart(t *testing.T) {
	// Nothing to check here, this calls other methods
	// that each have their own test
	t.SkipNow()
}

func TestWrite(t *testing.T) {
	gitter, _ := NewGitter(Token)
	supOutput := []byte("Hello World")
	if _, err := gitter.Write(supOutput); err != nil {
		t.Error(err)
	}
}

func TestGetRoomMsgs(t *testing.T) {
	gitter, _ := NewGitter(Token)

	// Test by passing correct room
	goodRoom := Room{Id: RoomId}
	msgCh := make(chan Message, 10)
	errCh := make(chan error, 1)

	gitter.getRoomMsgs(goodRoom, msgCh, errCh)
	select {
	case err := <-errCh:
		t.Error(err)
	case msg := <-msgCh:
		t.Logf("msg: %+v\n", msg)
	}

	// Test by passing wrong room
	badRoom := Room{Id: "NOT_" + RoomId}
	msgCh2 := make(chan Message, 10)
	errCh2 := make(chan error, 1)

	gitter.getRoomMsgs(badRoom, msgCh2, errCh2)
	select {
	case msg := <-msgCh2:
		t.Errorf("Should not have received msg: %+v\n", msg)
	case err := <-errCh2:
		t.Logf("Received correct error: %v\n", err)
	}
}

func TestProcessErrs(t *testing.T) {
	gitter, _ := NewGitter(Token)
	w := new(bytes.Buffer)
	errCh := make(chan error, 1)

	input := "Some error"
	errCh <- fmt.Errorf(input)
	gitter.processErrs(w, errCh)

	output := string(w.Bytes())
	output = strings.SplitN(output, " - ", 2)[1]
	output = strings.TrimSpace(output)
	if output != input {
		t.Errorf("Expected '%s', Received '%s'\n", input, output)
	}
}

func TestProcessMsgs(t *testing.T) {
	gitter, _ := NewGitter(Token)
	supBot := new(bytes.Buffer)
	msgCh := make(chan Message, 1)

	msg := new(Message)
	input := "Hello sup"
	msg.Text = "@" + gitter.user.Username + " " + input
	msg.Mentions = append(msg.Mentions, struct{ UserId string }{gitter.user.Id})
	msgCh <- *msg

	gitter.processMsgs(supBot, msgCh)

	output := string(supBot.Bytes())
	if output != input {
		t.Errorf("Expected '%s', Received '%s'\n", input, output)
	}
}

func TestWasMentioned(t *testing.T) {
	gitter, _ := NewGitter(Token)
	userID := gitter.user.Id

	// test when mentioned
	msg1 := new(Message)
	msg1.Mentions = append(msg1.Mentions, struct{ UserId string }{userID})
	t.Logf("msg1: %+v\n", msg1)
	if !gitter.wasMentioned(*msg1) {
		t.Error("Should have been mentioned.")
	}

	// test when NOT mentioned
	msg2 := new(Message)
	t.Logf("msg2: %+v\n", msg2)
	if gitter.wasMentioned(*msg2) {
		t.Error("Should NOT have been mentioned.")
	}
}
