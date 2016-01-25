package gitter

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	token  = "4b409f3d662592192095055ac603eaf106b0b92b"
	roomId = "56a3f2eae610378809bddc50"
)

func TestMain(m *testing.M) {
	server := setMockServer()
	defer server.Close()
	RestURL = server.URL
	os.Exit(m.Run())
}

func setMockServer() *httptest.Server {
	// Response taken from real calls to Gitter API
	userResponse := []byte(`[{"id":"56a3f554e610378809bddc9c","username":"supgitter","displayName":"supgitter","url":"/supgitter","avatarUrlSmall":"https://avatars0.githubusercontent.com/u/16857436?v=3&s=60","avatarUrlMedium":"https://avatars0.githubusercontent.com/u/16857436?v=3&s=128","v":1,"gv":"3"}]`)
	roomsResponse := []byte(`[{"id":"56a3f2eae610378809bddc50","name":"travellingprog/supbot","topic":"","uri":"travellingprog/supbot","oneToOne":false,"userCount":2,"unreadItems":0,"mentions":0,"lastAccessTime":"2016-01-24T15:51:30.513Z","lurk":false,"activity":true,"url":"/travellingprog/supbot","githubType":"REPO","security":"PUBLIC","noindex":false,"tags":[],"roomMember":true}]`)
	msgsResponse := []byte(`[{"id":"56a4306f8fbaf4220af8f817","text":"@travellingprog what's good?","html":"<span data-link-type=\"mention\" data-screen-name=\"travellingprog\" class=\"mention\">@travellingprog</span> what&#39;s good?","sent":"2016-01-24T02:01:19.407Z","fromUser":{"id":"56a3f554e610378809bddc9c","username":"supgitter","displayName":"supgitter","url":"/supgitter","avatarUrlSmall":"https://avatars0.githubusercontent.com/u/16857436?v=3&s=60","avatarUrlMedium":"https://avatars0.githubusercontent.com/u/16857436?v=3&s=128","v":1,"gv":"3"},"unread":false,"readBy":1,"urls":[],"mentions":[{"screenName":"travellingprog","userId":"56a3eab0e610378809bddb7d","userIds":[]}],"issues":[],"meta":[],"v":1},{"id":"56a50f97eaf741c118d49b13","text":"@supgitter local ping","html":"<span data-link-type=\"mention\" data-screen-name=\"supgitter\" class=\"mention\">@supgitter</span> local ping","sent":"2016-01-24T17:53:27.758Z","fromUser":{"id":"56a3eab0e610378809bddb7d","username":"travellingprog","displayName":"Erick Cardenas-Mendez","url":"/travellingprog","avatarUrlSmall":"https://avatars2.githubusercontent.com/u/3519160?v=3&s=60","avatarUrlMedium":"https://avatars2.githubusercontent.com/u/3519160?v=3&s=128","gv":"3"},"unread":false,"readBy":1,"urls":[],"mentions":[{"screenName":"supgitter","userId":"56a3f554e610378809bddc9c","userIds":[]}],"issues":[],"meta":[],"v":1}]`)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/user":
			w.Write(userResponse)
		case "/rooms":
			w.Write(roomsResponse)
		case "/rooms/" + roomId + "/chatMessages":
			w.Write(msgsResponse)
		}
	}

	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestGet(t *testing.T) {
	var err error

	var users []User
	if err = Get(RestURL+"/user", token, &users, "current user"); err != nil {
		t.Error(err)
	}
	t.Logf("users: %+v\n", users)

	var rooms []Room
	if err = Get(RestURL+"/rooms", token, &rooms, "rooms"); err != nil {
		t.Error(err)
	}
	t.Logf("rooms: %+v\n", rooms)

	var msgs []Message
	if err = Get(RestURL+"/rooms/"+roomId+"/chatMessages?limit=2", token, &msgs, "chat messages"); err != nil {
		t.Error(err)
	}
	t.Logf("msgs: %+v\n", msgs)
}

func TestNewGitter(t *testing.T) {
	gitter, err := NewGitter(token)
	if err != nil {
		t.Error(err)
	}
	t.Logf("gitter: %+v\n", gitter)
}

func TestGetMessages(t *testing.T) {
	gitter, err := NewGitter(token)
	if err != nil {
		t.Error(err)
	}
	t.Logf("gitter: %+v\n", gitter)
	// TODO: Test GetRoomMsgs()
}
