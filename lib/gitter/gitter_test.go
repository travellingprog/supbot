package gitter

import (
	"testing"
)

var (
	token = "4b409f3d662592192095055ac603eaf106b0b92b"
)

func TestNewGitter(t *testing.T) {
	gitter, err := NewGitter(token)
	if err != nil {
		t.Error(err)
	}

	gitter.Initialize()
}

// supgitter@gmail.com
// suppressly

// https://github.com/supgitter
// suppressly1

// Gitter token: 4b409f3d662592192095055ac603eaf106b0b92b

// Room: 56a3f2eae610378809bddc50

// User: 56a3f554e610378809bddc9c

// curl -i -H "Accept: application/json" -H "Authorization: Bearer 4b409f3d662592192095055ac603eaf106b0b92b" "https://stream.gitter.im/v1/rooms/56a3f2eae610378809bddc50/chatMessages"

// {
//     "id": "56a50f97eaf741c118d49b13",
//     "text": "@supgitter local ping",
//     "html": "<span data-link-type=\"mention\" data-screen-name=\"supgitter\" class=\"mention\">@supgitter</span> local ping",
//     "sent": "2016-01-24T17:53:27.758Z",
//     "fromUser": {
//         "id": "56a3eab0e610378809bddb7d",
//         "username": "travellingprog",
//         "displayName": "Erick Cardenas-Mendez",
//         "url": "/travellingprog",
//         "avatarUrlSmall": "https://avatars2.githubusercontent.com/u/3519160?v=3&s=60",
//         "avatarUrlMedium": "https://avatars2.githubusercontent.com/u/3519160?v=3&s=128",
//         "gv": "3"
//     },
//     "unread": true,
//     "readBy": 0,
//     "urls": [],
//     "mentions": [
//         {
//             "screenName": "supgitter",
//             "userId": "56a3f554e610378809bddc9c",
//             "userIds": []
//         }
//     ],
//     "issues": [],
//     "meta": [],
//     "v": 1
// }
