package maps

//go:generate osmium add-locations-to-ways testdata/speed_limits.osm -o testdata/speed_limits.osm.pbf --overwrite

import (
	"math"
	"testing"

	m "pfeifer.dev/mapd/math"
	ms "pfeifer.dev/mapd/settings"
)

func useOfflineSettings(t *testing.T, settings OfflineSettings) {
	t.Helper()
	original := DEFAULT_SETTINGS
	DEFAULT_SETTINGS = settings
	t.Cleanup(func() { DEFAULT_SETTINGS = original })
}

func setConditionalSpeedLimitControl(t *testing.T, enabled bool) {
	t.Helper()
	original := ms.Settings.ConditionalSpeedLimitControlEnabled
	ms.Settings.ConditionalSpeedLimitControlEnabled = enabled
	t.Cleanup(func() { ms.Settings.ConditionalSpeedLimitControlEnabled = original })
}

func generateFixtureWays(t *testing.T) map[int64]Way {
	t.Helper()
	settings := OfflineSettings{
		Box:             m.Box{MinPos: m.NewPosition(51, 17), MaxPos: m.NewPosition(51.25, 17.25)},
		OutputDirectory: t.TempDir(),
		InputFile:       "testdata/speed_limits.osm.pbf",
		Overlap:         0.001,
	}
	GenerateOffline(settings)
	useOfflineSettings(t, settings)

	tile, err := FindWaysAroundPosition(m.NewPosition(51.1, 17.05))
	if err != nil || !tile.Loaded {
		t.Fatalf("generated tile did not load: %v", err)
	}
	ways := map[int64]Way{}
	for i := range tile.Ways.Len() {
		way := tile.Ways.At(i)
		ways[way.Id()] = way
	}
	return ways
}

func TestGeneratedWaySpeedLimits(t *testing.T) {
	ways := generateFixtureWays(t)
	for _, tc := range []struct {
		name                    string
		wayID                   int64
		conditionalControl      bool
		forwardKph, backwardKph float64
	}{
		{"posted speed only", 1001, false, 50, 50},
		{"practical speed replaces the posted speed", 1002, false, 40, 40},
		{"posted speed for the direction beats practical speed for the way", 1003, false, 70, 60},
		{"practical speed for the direction beats every other speed", 1004, false, 50, 60},
		{"reversed one-way keeps practical speed with its direction of travel", 1005, false, 30, 50},
		{"practical speed applies while conditional control is off", 1006, false, 40, 40},
		{"applying conditional speed beats practical speed", 1006, true, 30, 30},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setConditionalSpeedLimitControl(t, tc.conditionalControl)
			way, found := ways[tc.wayID]
			if !found {
				t.Fatalf("way %d missing from the generated tile", tc.wayID)
			}
			for _, direction := range []struct {
				isForward bool
				kph       float64
			}{{true, tc.forwardKph}, {false, tc.backwardKph}} {
				got := way.EffectiveMaxSpeed(direction.isForward)
				if want := direction.kph * ms.KPH_TO_MS; math.Abs(got-want) > 1e-3 {
					t.Errorf("forward=%t: got %.3f m/s, want %.3f m/s (%.0f kph)", direction.isForward, got, want, direction.kph)
				}
			}
		})
	}
}

func TestFindWaysAroundPositionMissingTile(t *testing.T) {
	useOfflineSettings(t, OfflineSettings{OutputDirectory: t.TempDir()})

	pos := m.NewPosition(51.1, 17.05)
	tile, _ := FindWaysAroundPosition(pos)
	if tile.Loaded {
		t.Fatal("expected a missing tile to load nothing")
	}
	box := tile.Box()
	if !box.PosInside(pos) {
		t.Fatalf("expected the unloaded tile box %+v to contain %+v so the caller waits before retrying", box, pos)
	}
}
