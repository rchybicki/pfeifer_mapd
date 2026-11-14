package maps

import (
	"testing"

	ms "pfeifer.dev/mapd/settings"
)

func TestApplyMaxSpeedOverrideUpdatesDirectionalSpeeds(t *testing.T) {
	way := TmpWay{
		MaxSpeed:                  1,
		MaxSpeedForward:           2,
		MaxSpeedBackward:          3,
		MaxSpeedPractical:         4,
		MaxSpeedPracticalForward:  5,
		MaxSpeedPracticalBackward: 6,
	}

	applyMaxSpeedOverride(1167942324, &way)

	expected := 45 * ms.KPH_TO_MS
	if way.MaxSpeed != expected {
		t.Fatalf("expected max speed %f, got %f", expected, way.MaxSpeed)
	}
	if way.MaxSpeedForward != expected {
		t.Fatalf("expected forward max speed %f, got %f", expected, way.MaxSpeedForward)
	}
	if way.MaxSpeedBackward != expected {
		t.Fatalf("expected backward max speed %f, got %f", expected, way.MaxSpeedBackward)
	}
	if way.MaxSpeedPractical != expected {
		t.Fatalf("expected practical max speed %f, got %f", expected, way.MaxSpeedPractical)
	}
	if way.MaxSpeedPracticalForward != expected {
		t.Fatalf("expected forward practical max speed %f, got %f", expected, way.MaxSpeedPracticalForward)
	}
	if way.MaxSpeedPracticalBackward != expected {
		t.Fatalf("expected backward practical max speed %f, got %f", expected, way.MaxSpeedPracticalBackward)
	}
}

func TestApplyMaxSpeedOverrideDoesNotCreateDirectionalSpeeds(t *testing.T) {
	way := TmpWay{MaxSpeed: 1}

	applyMaxSpeedOverride(1167942324, &way)

	expected := 45 * ms.KPH_TO_MS
	if way.MaxSpeed != expected {
		t.Fatalf("expected max speed %f, got %f", expected, way.MaxSpeed)
	}
	if way.MaxSpeedForward != 0 {
		t.Fatalf("expected forward max speed to stay unset, got %f", way.MaxSpeedForward)
	}
	if way.MaxSpeedBackward != 0 {
		t.Fatalf("expected backward max speed to stay unset, got %f", way.MaxSpeedBackward)
	}
	if way.MaxSpeedPractical != expected {
		t.Fatalf("expected practical max speed %f, got %f", expected, way.MaxSpeedPractical)
	}
	if way.MaxSpeedPracticalForward != 0 {
		t.Fatalf("expected forward practical max speed to stay unset, got %f", way.MaxSpeedPracticalForward)
	}
	if way.MaxSpeedPracticalBackward != 0 {
		t.Fatalf("expected backward practical max speed to stay unset, got %f", way.MaxSpeedPracticalBackward)
	}
}

func TestApplyMaxSpeedOverrideNoopWithoutMatch(t *testing.T) {
	way := TmpWay{
		MaxSpeed:         1,
		MaxSpeedForward:  2,
		MaxSpeedBackward: 3,
	}

	applyMaxSpeedOverride(-1, &way)

	if way.MaxSpeed != 1 {
		t.Fatalf("expected max speed to remain unchanged, got %f", way.MaxSpeed)
	}
	if way.MaxSpeedForward != 2 {
		t.Fatalf("expected forward max speed to remain unchanged, got %f", way.MaxSpeedForward)
	}
	if way.MaxSpeedBackward != 3 {
		t.Fatalf("expected backward max speed to remain unchanged, got %f", way.MaxSpeedBackward)
	}
	if way.MaxSpeedPractical != 0 {
		t.Fatalf("expected practical max speed to remain unchanged, got %f", way.MaxSpeedPractical)
	}
}
