// Package botcontract defines the persistence-free Eventius port used by bot
// controllers and other conversational clients.
package botcontract

import (
	"context"
	"time"
)

// Event is the presentation-safe event view required by the legacy Eventius
// host adapter. It is a view model, not a persisted aggregate or DAL record.
type Event struct {
	ID              string
	Title           string
	Start           time.Time
	Location        string
	MinParticipants int
}

// Rollup is the presentation-safe attendance summary for an event.
type Rollup struct {
	Adults        int
	Children      int
	OpenAttention int
}

// EventWithRollup joins an event to the summary displayed by the bot.
type EventWithRollup struct {
	Event
	Rollup Rollup
}

// AttentionItem is the presentation-safe host-attention view.
type AttentionItem struct {
	EventID    string
	Kind       string
	Detail     string
	At         time.Time
	ResolvedAt *time.Time
}

// CreateEventSpec contains the fields an Eventius conversational client may
// create. It carries user intent, never a database request envelope.
type CreateEventSpec struct {
	Title           string
	Start           time.Time
	Location        string
	Vertical        string
	Label           string
	MinParticipants int
}

// Service is retained for the pre-vertical compatibility adapter. New bot
// controllers use EventService from canonical_service.go.
type Service interface {
	CreateEvent(ctx context.Context, userID, spaceID string, spec CreateEventSpec) (eventID string, err error)
	ListUpcoming(ctx context.Context, spaceID string, horizon time.Duration) ([]EventWithRollup, error)
	ListOpenAttention(ctx context.Context, spaceID string) (items []AttentionItem, events map[string]Event, err error)
	IssueOpenLink(ctx context.Context, spaceID, eventID string) (url string, err error)
}
