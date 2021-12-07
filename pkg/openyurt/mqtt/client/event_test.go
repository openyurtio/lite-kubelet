package client

import "testing"

func TestGetPreTopic(t *testing.T) {
	e := &events{namespace: "default"}
	t.Logf(e.GetPreTopic())
}
