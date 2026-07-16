// Package participation is the shared participation-response vocabulary for the
// Eventius occasion family (ToGethered, RSVP.express verticals, Eventius/Invitus).
//
// It is the single source of truth for three previously-diverging response
// scales, satisfying the "types needed by several extensions → ext-<id>/backend
// contract module" row of the extension-backend-architecture decision ladder
// (sneat-specs/standards/extension-backend-architecture.md).
//
// # Consumers
//
//   - togethered/backend — stores IntentDbo.Likelihood as a [Level]; the
//     Notifies/co-presence policy is ToGethered-owned (see NON-GOALS).
//   - rsvp-express verticals — personal-events (yes/no/maybe), practice-training
//     (available/unavailable/maybe as [Availability]).
//   - eventius/backend — [Coarse] replaces the local RsvpStatus enum.
//   - gameboard/backend — [GameResponse] replaces the local going/maybe/out enum.
//   - invitus mapping — maps accepted/declined → [Coarse]; that mapping stays in
//     sneat-core-modules, which imports this vocabulary.
//
// # NON-GOALS
//
// This package intentionally omits:
//
//   - Notification policy (who gets notified at which level): that is extension
//     policy owned by togethered/backend; this module must never grow predicates
//     such as NotifiesFollowers().
//   - Co-presence eligibility: likewise extension policy; co-presence at "likely"
//     is a ToGethered domain rule, not a vocabulary rule.
//   - Phrase or natural-language parsing: phrase→level mapping ("probably"→likely,
//     "should make it"→likely) is a conversational-layer concern that stays in
//     sneat-bots. Use [ParseLevel] and [ParseCoarse] only for strict wire-value
//     round-trips.
//
// # Scale changes
//
// Adding, removing, or reordering the values in [Level] or [Coarse] requires an
// explicit owner decision (alexander.trakhimenok@gmail.com) because every stored
// value and every consumer is affected.
package participation

import "fmt"

// Level is the five-point graded likelihood scale used in ToGethered intents.
// Values are stored as-is in Firestore; never change the string representations
// without a coordinated migration.
//
// Ordering (ascending): [LevelNo] < [LevelUnlikely] < [LevelMaybe] < [LevelLikely] < [LevelYes].
type Level string

const (
	// LevelNo — the participant will not attend.
	LevelNo Level = "no"
	// LevelUnlikely — the participant probably will not attend.
	LevelUnlikely Level = "unlikely"
	// LevelMaybe — the participant is undecided.
	LevelMaybe Level = "maybe"
	// LevelLikely — the participant will probably attend.
	LevelLikely Level = "likely"
	// LevelYes — the participant will attend.
	LevelYes Level = "yes"
)

// levelOrdinal maps each Level to its position in the ascending scale.
var levelOrdinal = map[Level]int{
	LevelNo:       0,
	LevelUnlikely: 1,
	LevelMaybe:    2,
	LevelLikely:   3,
	LevelYes:      4,
}

// Ordinal returns the integer position of the level in the ascending scale
// (no=0, unlikely=1, maybe=2, likely=3, yes=4).
// Returns -1 for an unrecognised value.
func (l Level) Ordinal() int {
	if ord, ok := levelOrdinal[l]; ok {
		return ord
	}
	return -1
}

// IsValid reports whether the level is one of the five known values.
func (l Level) IsValid() bool {
	_, ok := levelOrdinal[l]
	return ok
}

// AtLeast reports whether the level is greater than or equal to other in the
// ascending scale. Both values must be valid; panics otherwise.
func (l Level) AtLeast(other Level) bool {
	lo := l.Ordinal()
	oo := other.Ordinal()
	if lo < 0 {
		panic(fmt.Sprintf("participation: invalid Level %q", l))
	}
	if oo < 0 {
		panic(fmt.Sprintf("participation: invalid Level %q", other))
	}
	return lo >= oo
}

// Coarse returns the coarse three-point answer for this level:
//
//   - yes, likely → [CoarseYes]
//   - maybe       → [CoarseMaybe]
//   - unlikely, no → [CoarseNo]
//
// This is a lossy mapping; use [Coarse.Level] for the lossless embedding.
func (l Level) Coarse() Coarse {
	switch l {
	case LevelYes, LevelLikely:
		return CoarseYes
	case LevelMaybe:
		return CoarseMaybe
	default:
		return CoarseNo
	}
}

// ParseLevel parses a wire value into a [Level]. It accepts only the five
// canonical stored values; it does not interpret natural-language phrases.
// Returns an error for any unrecognised input.
func ParseLevel(s string) (Level, error) {
	l := Level(s)
	if !l.IsValid() {
		return "", fmt.Errorf("participation: unknown level %q", s)
	}
	return l, nil
}

// Coarse is the three-point quick-answer scale shared by:
//   - eventius/backend RsvpStatus (yes/no/maybe)
//   - rsvp-express personal-events (yes/no/maybe)
//   - invitus InviteResponse mapping (accepted→yes, declined→no)
//
// [CoarseYes], [CoarseNo], [CoarseMaybe] are the only valid values.
// See also [Availability] and [GameResponse] for alias renderings.
type Coarse string

const (
	// CoarseYes — the participant is attending (or has accepted).
	CoarseYes Coarse = "yes"
	// CoarseNo — the participant is not attending (or has declined).
	CoarseNo Coarse = "no"
	// CoarseMaybe — the participant is undecided.
	CoarseMaybe Coarse = "maybe"
)

// IsValid reports whether the coarse value is one of the three known values.
func (c Coarse) IsValid() bool {
	return c == CoarseYes || c == CoarseNo || c == CoarseMaybe
}

// Level returns the lossless embedding of the coarse answer into the graded scale:
//
//   - [CoarseYes]   → [LevelYes]
//   - [CoarseNo]    → [LevelNo]
//   - [CoarseMaybe] → [LevelMaybe]
//
// The reverse mapping [Level.Coarse] is lossy because two levels map to yes.
func (c Coarse) Level() Level {
	switch c {
	case CoarseYes:
		return LevelYes
	case CoarseNo:
		return LevelNo
	default:
		return LevelMaybe
	}
}

// ParseCoarse parses a wire value into a [Coarse]. It accepts only the three
// canonical values; it does not interpret natural-language phrases.
// Returns an error for any unrecognised input.
func ParseCoarse(s string) (Coarse, error) {
	c := Coarse(s)
	if !c.IsValid() {
		return "", fmt.Errorf("participation: unknown coarse value %q", s)
	}
	return c, nil
}

// Availability is an alias rendering of [Coarse] for the practice-training
// vertical (rsvp-express/practice-training) and any other availability-framed
// occasion. It is NOT a third independent scale; it maps directly onto [Coarse]:
//
//	[AvailabilityAvailable]   ≡ [CoarseYes]
//	[AvailabilityUnavailable] ≡ [CoarseNo]
//	[AvailabilityMaybe]       ≡ [CoarseMaybe]
//
// When stored, prefer the [Coarse] values; use [Availability] for UI labels and
// API surfaces that speak "availability" terminology.
type Availability string

const (
	// AvailabilityAvailable — the member is available for the session.
	AvailabilityAvailable Availability = "available"
	// AvailabilityUnavailable — the member is not available.
	AvailabilityUnavailable Availability = "unavailable"
	// AvailabilityMaybe — the member is unsure.
	AvailabilityMaybe Availability = "maybe"
)

// IsValid reports whether the availability value is one of the three known values.
func (a Availability) IsValid() bool {
	return a == AvailabilityAvailable || a == AvailabilityUnavailable || a == AvailabilityMaybe
}

// Coarse converts the availability label to its canonical [Coarse] equivalent.
func (a Availability) Coarse() Coarse {
	switch a {
	case AvailabilityAvailable:
		return CoarseYes
	case AvailabilityUnavailable:
		return CoarseNo
	default:
		return CoarseMaybe
	}
}

// AvailabilityFromCoarse converts a [Coarse] value to its [Availability] label.
func AvailabilityFromCoarse(c Coarse) Availability {
	switch c {
	case CoarseYes:
		return AvailabilityAvailable
	case CoarseNo:
		return AvailabilityUnavailable
	default:
		return AvailabilityMaybe
	}
}

// GameResponse is an alias rendering of [Coarse] for the gameboard/backend
// game-attendance RSVP (gameboard/backend/gameboard/game_invite.go:82-88),
// where the stored values are "going", "maybe", "out". It is NOT a third
// independent scale; it maps directly onto [Coarse]:
//
//	[GameResponseGoing] ≡ [CoarseYes]
//	[GameResponseMaybe] ≡ [CoarseMaybe]
//	[GameResponseOut]   ≡ [CoarseNo]
//
// When stored in new code, prefer [Coarse] values. Use [GameResponse] for
// backward-compatible read/write of existing gameboard records.
type GameResponse string

const (
	// GameResponseGoing — the player is attending the game.
	GameResponseGoing GameResponse = "going"
	// GameResponseMaybe — the player is unsure.
	GameResponseMaybe GameResponse = "maybe"
	// GameResponseOut — the player is not attending.
	GameResponseOut GameResponse = "out"
)

// IsValid reports whether the game response is one of the three known values.
func (g GameResponse) IsValid() bool {
	return g == GameResponseGoing || g == GameResponseMaybe || g == GameResponseOut
}

// Coarse converts the game response label to its canonical [Coarse] equivalent.
func (g GameResponse) Coarse() Coarse {
	switch g {
	case GameResponseGoing:
		return CoarseYes
	case GameResponseOut:
		return CoarseNo
	default:
		return CoarseMaybe
	}
}

// GameResponseFromCoarse converts a [Coarse] value to its [GameResponse] label.
func GameResponseFromCoarse(c Coarse) GameResponse {
	switch c {
	case CoarseYes:
		return GameResponseGoing
	case CoarseNo:
		return GameResponseOut
	default:
		return GameResponseMaybe
	}
}
