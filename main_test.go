package main_test

import (
	"fmt"
	. "github.com/adminibar/container-updater"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var example = `
 {
  "push_data": {
    "pushed_at": 1407181354,
    "images": [
      "imagehash1",
      "imagehash2",
      "imagehash3"
    ],
    "pusher": "advanderveer"
  },
  "callback_url": "https:\/\/registry.hub.docker.com\/u\/adminibar\/api\/hook\/213cbbf33hfje1g0j11i4c55afiba0afbe3i0\/",
  "repository": {
    "status": "Active",
    "description": "",
    "is_trusted": true,
    "full_description": "api\n===\n\n\nTODO:\n=====\n- ensure indexes for user collection",
    "repo_url": "https:\/\/registry.hub.docker.com\/u\/adminibar\/api\/",
    "owner": "adminibar",
    "is_official": false,
    "is_private": false,
    "name": "api",
    "namespace": "adminibar",
    "star_count": 0,
    "comment_count": 0,
    "date_created": 1406993033,
    "dockerfile": "FROM google\/golang:1.3\n\n#we usegodep\nRUN go get github.com\/tools\/godep\nRUN mkdir -p \/gopath\/src\/github.com\/adminibar\/api\n\nWORKDIR \/gopath\/src\/github.com\/adminibar\/api\nADD . \/gopath\/src\/github.com\/adminibar\/api\/\nRUN godep go build -o \/gopath\/bin\/adminibarapi\n\nEXPOSE 8000\nCMD []\nENTRYPOINT [\"\/gopath\/bin\/adminibarapi\"]",
    "repo_name": "adminibar\/api"
  }
}
`

func TestDefaultScenario(t *testing.T) {
	r, _ := http.NewRequest("POST", fmt.Sprint("/"), strings.NewReader(example))
	rec := httptest.NewRecorder()

	HookHandler(rec, r)

	if rec.Code != 200 {
		t.Error(rec.Body.String())
	}
}

func TestArgParsing(t *testing.T) {
	os.Args = append(os.Args, "adminibar/api@0.1.2", "adminibar/ui@latest")

	spec, err := CreateSpec()
	if err != nil {
		t.Error(err)
	}

	if len(spec) != 2 {
		t.Error("expected two entries in spec")
	}

	if spec[0].TagPattern != "0.1.2" || spec[0].RepoName != "adminibar/api" {
		t.Error("Expected test to have correct tag")
	}

	if spec[1].TagPattern != "latest" || spec[1].RepoName != "adminibar/ui" {
		t.Error("Expected test to have correct tag")
	}
}
