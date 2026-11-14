package maps

import (
	"testing"

	"capnproto.org/go/capnp/v3"
	"pfeifer.dev/mapd/cereal/offline"
)

type wayMaxSpeeds struct {
	maxSpeed, forward, backward                    float64
	practical, practicalForward, practicalBackward float64
}

func newTestWay(t *testing.T, speeds wayMaxSpeeds) Way {
	t.Helper()
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("failed to create capnp message: %v", err)
	}
	raw, err := offline.NewRootWay(seg)
	if err != nil {
		t.Fatalf("failed to create offline way: %v", err)
	}
	raw.SetMaxSpeed(speeds.maxSpeed)
	raw.SetMaxSpeedForward(speeds.forward)
	raw.SetMaxSpeedBackward(speeds.backward)
	raw.SetMaxSpeedPractical(speeds.practical)
	raw.SetMaxSpeedPracticalForward(speeds.practicalForward)
	raw.SetMaxSpeedPracticalBackward(speeds.practicalBackward)
	return NewWay(raw)
}

func TestMaxSpeedForDirectionPrecedence(t *testing.T) {
	for _, tc := range []struct {
		name      string
		speeds    wayMaxSpeeds
		isForward bool
		want      float64
	}{
		{"practical forward wins", wayMaxSpeeds{maxSpeed: 25, practical: 20, forward: 22, practicalForward: 19}, true, 19},
		{"practical backward wins", wayMaxSpeeds{maxSpeed: 25, practical: 20, backward: 22, practicalBackward: 18}, false, 18},
		{"posted directional speed beats shared practical speed", wayMaxSpeeds{maxSpeed: 25, practical: 20, forward: 22}, true, 22},
		{"shared practical speed beats posted speed", wayMaxSpeeds{maxSpeed: 25, practical: 20}, true, 20},
		{"practical speed for the other direction is ignored", wayMaxSpeeds{maxSpeed: 25, practicalBackward: 18}, true, 25},
		{"posted speed is the fallback", wayMaxSpeeds{maxSpeed: 25}, true, 25},
	} {
		t.Run(tc.name, func(t *testing.T) {
			way := newTestWay(t, tc.speeds)
			if got := way.EffectiveMaxSpeed(tc.isForward); got != tc.want {
				t.Fatalf("got %f, want %f", got, tc.want)
			}
		})
	}
}
