package gunfish

import (
	"encoding/json"
	"testing"

	"github.com/kayac/Gunfish/apns"
)

var (
	jstr = `{"aps":{"alert":"hoge","badge":1,"sound":"default"},"mio":"hoge","uid":"hoge"}`
)

func TestUnmarshal(t *testing.T) {
	var payload apns.Payload

	err := json.Unmarshal([]byte(jstr), &payload)
	if err != nil {
		t.Error(err)
	}

	pjson, err := payload.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	if string(pjson) != jstr {
		t.Errorf("Expected %s, but got %s", jstr, pjson)
	}
}
