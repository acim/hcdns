package hcdns_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"go.acim.net/hcdns"
)

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	value := bytes.NewBufferString(`"2022-05-27 06:21:24.324 +0000 UTC"`)
	want := "2022-05-27T06:21:24Z"

	var got hcdns.Time

	if err := json.Unmarshal(value.Bytes(), &got); err != nil {
		t.Fatal(err)
	}

	if got.String() != want {
		t.Errorf("UnmarshalJSON(%s)=%s; want %s", value, got, want)
	}
}
