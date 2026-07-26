// Package facade4eventius defines the platform-neutral Eventius application
// contract used by Telegram, web, mobile, and future delivery surfaces.
//
// Eventius owns event lifecycle, authorization, and event-oriented queries.
// Calendarius owns the canonical Happening record and Invitus owns canonical
// Invite and InviteResponse records. Implementations compose those capabilities;
// consumers never call either persistence implementation directly.
package facade4eventius

import (
	"context"
	"time"

	"github.com/sneat-co/ext-eventius/backend/participation"
)

// EventStatus is the Eventius lifecycle state projected from the canonical
// Calendarius Happening.
type EventStatus string

const (
	EventStatusActive    EventStatus = "active"
	EventStatusArchived  EventStatus = "archived"
	EventStatusCancelled EventStatus = "cancelled"
)

// ScheduleState is derived from an event Happening's planning slot. It is a
// read-model value, not separately persisted state.
type ScheduleState string

const (
	// ScheduleStatePlanning means that date or time is still undecided.
	ScheduleStatePlanning ScheduleState = "planning"
	// ScheduleStateScheduled means that both date and start time are known.
	ScheduleStateScheduled ScheduleState = "scheduled"
)

// Event is the persistence-neutral Eventius projection of a Calendarius
// Happening whose kind is "event". Event ID and Happening ID are the same
// canonical identity.
type Event struct {
	ID          string      `json:"id"`
	SpaceID     string      `json:"spaceID"`
	Title       string      `json:"title"`
	Date        string      `json:"date,omitempty"`
	Time        string      `json:"time,omitempty"`
	Location    string      `json:"location,omitempty"`
	Description string      `json:"description,omitempty"`
	Status      EventStatus `json:"status"`
	CreatedBy   string      `json:"createdBy"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// ScheduleState derives whether an event is ready to appear on a calendar.
func (e Event) ScheduleState() ScheduleState {
	if e.Date != "" && e.Time != "" {
		return ScheduleStateScheduled
	}
	return ScheduleStatePlanning
}

// CreateEventRequest is user intent, not a storage envelope. Date, time,
// location, and description are independently optional.
type CreateEventRequest struct {
	RequestID   string `json:"requestID"`
	SpaceID     string `json:"spaceID"`
	Title       string `json:"title"`
	Date        string `json:"date,omitempty"`
	Time        string `json:"time,omitempty"`
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateEventRequest uses pointers so omitted and explicitly cleared planning
// fields remain distinct.
type UpdateEventRequest struct {
	RequestID   string  `json:"requestID"`
	SpaceID     string  `json:"spaceID"`
	EventID     string  `json:"eventID"`
	Title       *string `json:"title,omitempty"`
	Date        *string `json:"date,omitempty"`
	Time        *string `json:"time,omitempty"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
}

// InviteKind distinguishes a recipient-bound invitation from a shareable one.
type InviteKind string

const (
	InviteKindPersonal InviteKind = "personal"
	InviteKindOpen     InviteKind = "open"
)

// InviteStatus is the invitation lifecycle projected from Invitus.
type InviteStatus string

const (
	InviteStatusActive  InviteStatus = "active"
	InviteStatusRevoked InviteStatus = "revoked"
	InviteStatusExpired InviteStatus = "expired"
)

// Invite is Eventius's event-oriented view of an Invitus invitation. Its ID is
// never overloaded with the Event ID.
type Invite struct {
	ID            string       `json:"id"`
	SpaceID       string       `json:"spaceID"`
	EventID       string       `json:"eventID"`
	Kind          InviteKind   `json:"kind"`
	Status        InviteStatus `json:"status"`
	CreatedBy     string       `json:"createdBy"`
	InviteeUserID string       `json:"inviteeUserID,omitempty"`
	CreatedAt     time.Time    `json:"createdAt"`
	ExpiresAt     time.Time    `json:"expiresAt,omitempty"`
	GoingCount    int          `json:"goingCount,omitempty"`
	MaybeCount    int          `json:"maybeCount,omitempty"`
	DeclinedCount int          `json:"declinedCount,omitempty"`
}

// CreateInviteRequest creates either a personal invitation (InviteeUserID set)
// or an open invitation (InviteeUserID empty). RequestID makes retries
// idempotent.
type CreateInviteRequest struct {
	RequestID     string `json:"requestID"`
	SpaceID       string `json:"spaceID"`
	EventID       string `json:"eventID"`
	InviteeUserID string `json:"inviteeUserID,omitempty"`
}

// Response is Eventius's event-oriented view of a canonical Invitus
// InviteResponse.
type Response struct {
	ID          string               `json:"id"`
	InviteID    string               `json:"inviteID"`
	EventID     string               `json:"eventID"`
	UserID      string               `json:"userID"`
	Answer      participation.Coarse `json:"answer"`
	SubmittedAt time.Time            `json:"submittedAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
}

// InvitationContext contains everything a client needs to render an invitation
// and the current user's existing answer, if any.
type InvitationContext struct {
	Invite   Invite    `json:"invite"`
	Event    Event     `json:"event"`
	Response *Response `json:"response,omitempty"`
}

// RespondRequest records or changes the authenticated user's event answer.
// A replay with the same RequestID is idempotent; a later RequestID may change
// the answer without creating a duplicate response.
type RespondRequest struct {
	RequestID string               `json:"requestID"`
	InviteID  string               `json:"inviteID"`
	Answer    participation.Coarse `json:"answer"`
}

// Service is the reusable Eventius application facade. Implementations own
// authorization, orchestration, validation, idempotency, and observability.
type Service interface {
	CreateEvent(ctx context.Context, actorUserID string, request CreateEventRequest) (Event, error)
	UpdateEvent(ctx context.Context, actorUserID string, request UpdateEventRequest) (Event, error)
	GetEvent(ctx context.Context, actorUserID, spaceID, eventID string) (Event, error)
	ListEvents(ctx context.Context, actorUserID, spaceID string) ([]Event, error)

	CreateInvite(ctx context.Context, actorUserID string, request CreateInviteRequest) (Invite, error)
	GetInvitation(ctx context.Context, actorUserID, inviteID string) (InvitationContext, error)
	Respond(ctx context.Context, actorUserID string, request RespondRequest) (Response, error)
	ListResponses(ctx context.Context, actorUserID, spaceID, eventID string) ([]Response, error)
}
