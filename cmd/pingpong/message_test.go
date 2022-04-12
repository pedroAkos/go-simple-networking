package pingpong

import (
	"testing"
)

func TestSerialization(t *testing.T) {
	ping := Ping{}
	pong := Pong{}

	p, _ := ping.Deserialize(nil)
	if p.String() != ping.String() {
		t.Errorf("p.String() = %v; want %v", p, ping)
	}
	p, _ = pong.Deserialize(nil)
	if p.String() != pong.String() {
		t.Errorf("p.String() = %v; want %v", p, pong)
	}
}
