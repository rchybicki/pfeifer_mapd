package maps

import (
	"testing"

	"capnproto.org/go/capnp/v3"
	"pfeifer.dev/mapd/cereal/offline"
)

func tileWithWaySize(t *testing.T, size capnp.ObjectSize) []byte {
	t.Helper()
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	root, err := offline.NewRootOffline(seg)
	if err != nil {
		t.Fatal(err)
	}
	ways, err := capnp.NewCompositeList(seg, size, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := root.SetWays(offline.Way_List(ways)); err != nil {
		t.Fatal(err)
	}
	data, err := msg.MarshalPacked()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestReadOfflineRejectsOnlyTheOldPracticalSpeedLayout(t *testing.T) {
	cases := []struct {
		name   string
		size   capnp.ObjectSize
		loaded bool
	}{
		{"v2.0.6 based fork", capnp.ObjectSize{DataSize: 96, PointerCount: 4}, false},
		{"upstream v2.0.6", capnp.ObjectSize{DataSize: 72, PointerCount: 4}, true},
		{"upstream v2.4.0", capnp.ObjectSize{DataSize: 80, PointerCount: 7}, true},
		{"current", capnp.ObjectSize{DataSize: 104, PointerCount: 7}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tile := ReadOffline(tileWithWaySize(t, tc.size))
			if tile.Loaded != tc.loaded {
				t.Fatalf("Loaded = %v, want %v", tile.Loaded, tc.loaded)
			}
		})
	}
}
