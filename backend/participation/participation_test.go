package participation_test

import (
	"testing"

	"github.com/sneat-co/ext-eventius/backend/participation"
)

// TestLevelOrdinal verifies the ascending ordering and that all five levels
// report a valid, strictly-increasing ordinal.
func TestLevelOrdinal(t *testing.T) {
	ordered := []participation.Level{
		participation.LevelNo,
		participation.LevelUnlikely,
		participation.LevelMaybe,
		participation.LevelLikely,
		participation.LevelYes,
	}
	for i := 0; i < len(ordered)-1; i++ {
		lo := ordered[i].Ordinal()
		hi := ordered[i+1].Ordinal()
		if lo >= hi {
			t.Errorf("level %q ordinal %d must be less than %q ordinal %d", ordered[i], lo, ordered[i+1], hi)
		}
	}
}

// TestLevelOrdinalUnknown verifies that an unrecognised level returns -1.
func TestLevelOrdinalUnknown(t *testing.T) {
	if got := participation.Level("unknown").Ordinal(); got != -1 {
		t.Errorf("unknown level ordinal: got %d, want -1", got)
	}
}

// TestLevelIsValid covers every valid constant and an invalid value.
func TestLevelIsValid(t *testing.T) {
	valid := []participation.Level{
		participation.LevelNo,
		participation.LevelUnlikely,
		participation.LevelMaybe,
		participation.LevelLikely,
		participation.LevelYes,
	}
	for _, l := range valid {
		if !l.IsValid() {
			t.Errorf("level %q should be valid", l)
		}
	}
	if participation.Level("going").IsValid() {
		t.Error("level 'going' should not be valid")
	}
}

// TestLevelAtLeast is a comprehensive table test.
func TestLevelAtLeast(t *testing.T) {
	tests := []struct {
		level participation.Level
		other participation.Level
		want  bool
	}{
		{participation.LevelYes, participation.LevelYes, true},
		{participation.LevelYes, participation.LevelLikely, true},
		{participation.LevelYes, participation.LevelMaybe, true},
		{participation.LevelYes, participation.LevelUnlikely, true},
		{participation.LevelYes, participation.LevelNo, true},
		{participation.LevelLikely, participation.LevelYes, false},
		{participation.LevelLikely, participation.LevelLikely, true},
		{participation.LevelMaybe, participation.LevelLikely, false},
		{participation.LevelMaybe, participation.LevelMaybe, true},
		{participation.LevelUnlikely, participation.LevelMaybe, false},
		{participation.LevelNo, participation.LevelUnlikely, false},
		{participation.LevelNo, participation.LevelNo, true},
	}
	for _, tt := range tests {
		got := tt.level.AtLeast(tt.other)
		if got != tt.want {
			t.Errorf("%q.AtLeast(%q): got %v, want %v", tt.level, tt.other, got, tt.want)
		}
	}
}

// TestLevelCoarseMapping verifies the total mapping of Level → Coarse.
func TestLevelCoarseMapping(t *testing.T) {
	tests := []struct {
		level participation.Level
		want  participation.Coarse
	}{
		{participation.LevelYes, participation.CoarseYes},
		{participation.LevelLikely, participation.CoarseYes},
		{participation.LevelMaybe, participation.CoarseMaybe},
		{participation.LevelUnlikely, participation.CoarseNo},
		{participation.LevelNo, participation.CoarseNo},
	}
	for _, tt := range tests {
		if got := tt.level.Coarse(); got != tt.want {
			t.Errorf("Level(%q).Coarse() = %q, want %q", tt.level, got, tt.want)
		}
	}
}

// TestCoarseIsValid covers all three valid values and an invalid one.
func TestCoarseIsValid(t *testing.T) {
	for _, c := range []participation.Coarse{
		participation.CoarseYes,
		participation.CoarseNo,
		participation.CoarseMaybe,
	} {
		if !c.IsValid() {
			t.Errorf("coarse %q should be valid", c)
		}
	}
	if participation.Coarse("going").IsValid() {
		t.Error("coarse 'going' should not be valid")
	}
}

// TestCoarseLevelMapping verifies the total mapping of Coarse → Level.
func TestCoarseLevelMapping(t *testing.T) {
	tests := []struct {
		coarse participation.Coarse
		want   participation.Level
	}{
		{participation.CoarseYes, participation.LevelYes},
		{participation.CoarseNo, participation.LevelNo},
		{participation.CoarseMaybe, participation.LevelMaybe},
	}
	for _, tt := range tests {
		if got := tt.coarse.Level(); got != tt.want {
			t.Errorf("Coarse(%q).Level() = %q, want %q", tt.coarse, got, tt.want)
		}
	}
}

// TestRoundTripCoarseLevelCoarse verifies the identity:
// Coarse → Level → Coarse is the identity function for all three values.
func TestRoundTripCoarseLevelCoarse(t *testing.T) {
	for _, c := range []participation.Coarse{
		participation.CoarseYes,
		participation.CoarseNo,
		participation.CoarseMaybe,
	} {
		if got := c.Level().Coarse(); got != c {
			t.Errorf("round-trip Coarse(%q).Level().Coarse() = %q, want %q", c, got, c)
		}
	}
}

// TestParseLevel covers valid and invalid inputs.
func TestParseLevel(t *testing.T) {
	valid := []string{"no", "unlikely", "maybe", "likely", "yes"}
	for _, s := range valid {
		l, err := participation.ParseLevel(s)
		if err != nil {
			t.Errorf("ParseLevel(%q) unexpected error: %v", s, err)
		}
		if string(l) != s {
			t.Errorf("ParseLevel(%q) = %q, want %q", s, l, s)
		}
	}
	if _, err := participation.ParseLevel("going"); err == nil {
		t.Error("ParseLevel('going') should return an error")
	}
	if _, err := participation.ParseLevel(""); err == nil {
		t.Error("ParseLevel('') should return an error")
	}
}

// TestParseCoarse covers valid and invalid inputs.
func TestParseCoarse(t *testing.T) {
	valid := []string{"yes", "no", "maybe"}
	for _, s := range valid {
		c, err := participation.ParseCoarse(s)
		if err != nil {
			t.Errorf("ParseCoarse(%q) unexpected error: %v", s, err)
		}
		if string(c) != s {
			t.Errorf("ParseCoarse(%q) = %q, want %q", s, c, s)
		}
	}
	if _, err := participation.ParseCoarse("going"); err == nil {
		t.Error("ParseCoarse('going') should return an error")
	}
	if _, err := participation.ParseCoarse(""); err == nil {
		t.Error("ParseCoarse('') should return an error")
	}
}

// TestAvailabilityConversions verifies the total mapping of Availability ↔ Coarse.
func TestAvailabilityConversions(t *testing.T) {
	tests := []struct {
		avail  participation.Availability
		coarse participation.Coarse
	}{
		{participation.AvailabilityAvailable, participation.CoarseYes},
		{participation.AvailabilityUnavailable, participation.CoarseNo},
		{participation.AvailabilityMaybe, participation.CoarseMaybe},
	}
	for _, tt := range tests {
		if got := tt.avail.Coarse(); got != tt.coarse {
			t.Errorf("Availability(%q).Coarse() = %q, want %q", tt.avail, got, tt.coarse)
		}
		if got := participation.AvailabilityFromCoarse(tt.coarse); got != tt.avail {
			t.Errorf("AvailabilityFromCoarse(%q) = %q, want %q", tt.coarse, got, tt.avail)
		}
	}
}

// TestAvailabilityIsValid covers all values.
func TestAvailabilityIsValid(t *testing.T) {
	for _, a := range []participation.Availability{
		participation.AvailabilityAvailable,
		participation.AvailabilityUnavailable,
		participation.AvailabilityMaybe,
	} {
		if !a.IsValid() {
			t.Errorf("availability %q should be valid", a)
		}
	}
	if participation.Availability("yes").IsValid() {
		t.Error("availability 'yes' should not be valid")
	}
}

// TestGameResponseConversions verifies the total mapping of GameResponse ↔ Coarse.
// Covers gameboard/backend/gameboard/game_invite.go:82-88 (going/maybe/out scale).
func TestGameResponseConversions(t *testing.T) {
	tests := []struct {
		game   participation.GameResponse
		coarse participation.Coarse
	}{
		{participation.GameResponseGoing, participation.CoarseYes},
		{participation.GameResponseMaybe, participation.CoarseMaybe},
		{participation.GameResponseOut, participation.CoarseNo},
	}
	for _, tt := range tests {
		if got := tt.game.Coarse(); got != tt.coarse {
			t.Errorf("GameResponse(%q).Coarse() = %q, want %q", tt.game, got, tt.coarse)
		}
		if got := participation.GameResponseFromCoarse(tt.coarse); got != tt.game {
			t.Errorf("GameResponseFromCoarse(%q) = %q, want %q", tt.coarse, got, tt.game)
		}
	}
}

// TestGameResponseIsValid covers all values.
func TestGameResponseIsValid(t *testing.T) {
	for _, g := range []participation.GameResponse{
		participation.GameResponseGoing,
		participation.GameResponseMaybe,
		participation.GameResponseOut,
	} {
		if !g.IsValid() {
			t.Errorf("game response %q should be valid", g)
		}
	}
	if participation.GameResponse("yes").IsValid() {
		t.Error("game response 'yes' should not be valid")
	}
}

// TestRoundTripAvailabilityCoarseAvailability verifies the identity:
// Availability → Coarse → Availability is the identity for all three values.
func TestRoundTripAvailabilityCoarseAvailability(t *testing.T) {
	for _, a := range []participation.Availability{
		participation.AvailabilityAvailable,
		participation.AvailabilityUnavailable,
		participation.AvailabilityMaybe,
	} {
		if got := participation.AvailabilityFromCoarse(a.Coarse()); got != a {
			t.Errorf("round-trip Availability(%q).Coarse() → AvailabilityFromCoarse = %q, want %q", a, got, a)
		}
	}
}

// TestRoundTripGameResponseCoarseGameResponse verifies the identity:
// GameResponse → Coarse → GameResponse is the identity for all three values.
func TestRoundTripGameResponseCoarseGameResponse(t *testing.T) {
	for _, g := range []participation.GameResponse{
		participation.GameResponseGoing,
		participation.GameResponseMaybe,
		participation.GameResponseOut,
	} {
		if got := participation.GameResponseFromCoarse(g.Coarse()); got != g {
			t.Errorf("round-trip GameResponse(%q).Coarse() → GameResponseFromCoarse = %q, want %q", g, got, g)
		}
	}
}
