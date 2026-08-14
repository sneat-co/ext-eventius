package facade4eventius

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/sneat-co/ext-eventius/backend/participation"
)

// ErrInvalidCompetiosAttendanceRequest means a caller supplied an incomplete
// or unsafe cross-product attendance request. Implementations must reject it
// before allocating any Eventius record.
var ErrInvalidCompetiosAttendanceRequest = errors.New("eventius: invalid Competios attendance request")

// CompetiosEventKey is the opaque external identity that makes Ensure calls
// idempotent. Eventius never parses it or treats it as a public URL.
type CompetiosEventKey string
type CompetiosTournamentKey string
type CompetiosCompetitionKey string
type CompetiosEntryKey string

// CompetiosRegistrationKey is the opaque external identity of one confirmed
// Competios registration. It is not an Eventius invitation ID.
type CompetiosRegistrationKey string

// CompetiosInviteeKey is the opaque identity of exactly one invitee within an
// enrolled Competios Entry. Competios owns its construction and may encode the
// Entry lifecycle revision in it; Eventius MUST treat it as an unparsed opaque
// value. Together with the Event, Tournament, Competition, Entry and
// Registration keys, it makes an attendance invitation idempotent per invitee
// lifecycle.
type CompetiosInviteeKey string

// AttendanceEventID and AttendanceInvitationID are Eventius-owned opaque
// identities. They deliberately cannot be substituted for one another.
type AttendanceEventID string
type AttendanceInvitationID string

// CalendarEventRef points at the already-created canonical Calendarius
// kind=event Happening. The first Competios integration never asks Eventius to
// create a second Event or to guess which Space should own one.
type CalendarEventRef struct {
	SpaceID     string `json:"spaceID"`
	HappeningID string `json:"happeningID"`
}

type AttendanceEventState string

const (
	AttendanceEventActive    AttendanceEventState = "active"
	AttendanceEventCancelled AttendanceEventState = "cancelled"
)

type AttendanceInvitationState string

const (
	AttendanceInvitationActive  AttendanceInvitationState = "active"
	AttendanceInvitationRevoked AttendanceInvitationState = "revoked"
)

// AttendanceResponderKind says why the bound account may answer this targeted
// invitation. It is authority input only and is deliberately omitted from the
// safe projection returned to Competios.
type AttendanceResponderKind string

const (
	AttendanceResponderAccount  AttendanceResponderKind = "account"
	AttendanceResponderGuardian AttendanceResponderKind = "guardian"
)

// AttendanceResponderRef binds RSVP submit to one authenticated account.
// AccountID is an opaque identity reference; it is never an invitation token,
// email address, contact payload, or public projection field.
type AttendanceResponderRef struct {
	Kind      AttendanceResponderKind `json:"kind"`
	AccountID string                  `json:"accountID"`
}

// EnsureAttendanceEventRequest describes the presentation-safe attendance
// event that Eventius should attach attendance to for a Competios Event.
// CalendarEvent must already identify the canonical, scheduled Calendarius
// kind=event Happening. RequestID makes a retried command idempotent;
// CompetiosEventKey prevents conflicting correlation to another Happening.
type EnsureAttendanceEventRequest struct {
	RequestID         string            `json:"requestID"`
	CompetiosEventKey CompetiosEventKey `json:"competiosEventKey"`
	CalendarEvent     CalendarEventRef  `json:"calendarEvent"`
}

// EnsureAttendanceInvitationRequest makes one attendance invitation for a
// confirmed registration and one invitee lifecycle. No RSVP token belongs to
// this cross-product contract: any token is Eventius transport authority and
// must remain private to Eventius.
type EnsureAttendanceInvitationRequest struct {
	RequestID                string                   `json:"requestID"`
	AttendanceEventID        AttendanceEventID        `json:"attendanceEventID"`
	CompetiosRegistrationKey CompetiosRegistrationKey `json:"competiosRegistrationKey"`
	CompetiosTournamentKey   CompetiosTournamentKey   `json:"competiosTournamentKey"`
	CompetiosCompetitionKey  CompetiosCompetitionKey  `json:"competiosCompetitionKey"`
	CompetiosEntryKey        CompetiosEntryKey        `json:"competiosEntryKey"`
	CompetiosInviteeKey      CompetiosInviteeKey      `json:"competiosInviteeKey"`
	Responder                AttendanceResponderRef   `json:"responder"`
}

// GetAttendanceInviteeStatusRequest identifies exactly one safe invitation
// status projection. It has no token, contact, payment, or response-detail
// fields. Eventius MUST NOT parse CompetiosInviteeKey.
type GetAttendanceInviteeStatusRequest struct {
	CompetiosEventKey        CompetiosEventKey        `json:"competiosEventKey"`
	CompetiosTournamentKey   CompetiosTournamentKey   `json:"competiosTournamentKey"`
	CompetiosCompetitionKey  CompetiosCompetitionKey  `json:"competiosCompetitionKey"`
	CompetiosEntryKey        CompetiosEntryKey        `json:"competiosEntryKey"`
	CompetiosRegistrationKey CompetiosRegistrationKey `json:"competiosRegistrationKey"`
	CompetiosInviteeKey      CompetiosInviteeKey      `json:"competiosInviteeKey"`
}

// AttendanceStatusProjection is safe to return to Competios. It intentionally
// omits invitation URLs, RSVP tokens, invitee names, contact fields, and any
// other authority-bearing data.
type AttendanceStatusProjection struct {
	CompetiosEventKey        CompetiosEventKey         `json:"competiosEventKey"`
	CompetiosRegistrationKey CompetiosRegistrationKey  `json:"competiosRegistrationKey,omitempty"`
	CompetiosTournamentKey   CompetiosTournamentKey    `json:"competiosTournamentKey,omitempty"`
	CompetiosCompetitionKey  CompetiosCompetitionKey   `json:"competiosCompetitionKey,omitempty"`
	CompetiosEntryKey        CompetiosEntryKey         `json:"competiosEntryKey,omitempty"`
	CompetiosInviteeKey      CompetiosInviteeKey       `json:"competiosInviteeKey,omitempty"`
	AttendanceEventID        AttendanceEventID         `json:"attendanceEventID"`
	AttendanceInvitationID   AttendanceInvitationID    `json:"attendanceInvitationID,omitempty"`
	EventState               AttendanceEventState      `json:"eventState"`
	InvitationState          AttendanceInvitationState `json:"invitationState,omitempty"`
	Response                 *participation.Coarse     `json:"response,omitempty"`
	RespondedAt              *time.Time                `json:"respondedAt,omitempty"`
}

func ValidateEnsureAttendanceEventRequest(value EnsureAttendanceEventRequest) error {
	if strings.TrimSpace(value.RequestID) == "" || strings.TrimSpace(string(value.CompetiosEventKey)) == "" || strings.TrimSpace(value.CalendarEvent.SpaceID) == "" || strings.TrimSpace(value.CalendarEvent.HappeningID) == "" {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

func ValidateEnsureAttendanceInvitationRequest(value EnsureAttendanceInvitationRequest) error {
	if strings.TrimSpace(value.RequestID) == "" || strings.TrimSpace(string(value.AttendanceEventID)) == "" || strings.TrimSpace(string(value.CompetiosRegistrationKey)) == "" || strings.TrimSpace(string(value.CompetiosTournamentKey)) == "" || strings.TrimSpace(string(value.CompetiosCompetitionKey)) == "" || strings.TrimSpace(string(value.CompetiosEntryKey)) == "" || strings.TrimSpace(string(value.CompetiosInviteeKey)) == "" || strings.TrimSpace(value.Responder.AccountID) == "" {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.Responder.Kind != AttendanceResponderAccount && value.Responder.Kind != AttendanceResponderGuardian {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

// ValidateGetAttendanceInviteeStatusRequest rejects a partial external tuple
// so a provider cannot choose an arbitrary invitee for a shared registration.
func ValidateGetAttendanceInviteeStatusRequest(value GetAttendanceInviteeStatusRequest) error {
	if strings.TrimSpace(string(value.CompetiosEventKey)) == "" || strings.TrimSpace(string(value.CompetiosTournamentKey)) == "" || strings.TrimSpace(string(value.CompetiosCompetitionKey)) == "" || strings.TrimSpace(string(value.CompetiosEntryKey)) == "" || strings.TrimSpace(string(value.CompetiosRegistrationKey)) == "" || strings.TrimSpace(string(value.CompetiosInviteeKey)) == "" {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

func ValidateAttendanceStatusProjection(value AttendanceStatusProjection) error {
	if strings.TrimSpace(string(value.CompetiosEventKey)) == "" || strings.TrimSpace(string(value.AttendanceEventID)) == "" {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.EventState != AttendanceEventActive && value.EventState != AttendanceEventCancelled {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.AttendanceInvitationID == "" {
		if value.CompetiosRegistrationKey != "" || value.CompetiosTournamentKey != "" || value.CompetiosCompetitionKey != "" || value.CompetiosEntryKey != "" || value.CompetiosInviteeKey != "" || value.InvitationState != "" || value.Response != nil || value.RespondedAt != nil {
			return ErrInvalidCompetiosAttendanceRequest
		}
		return nil
	}
	if strings.TrimSpace(string(value.CompetiosRegistrationKey)) == "" || strings.TrimSpace(string(value.CompetiosTournamentKey)) == "" || strings.TrimSpace(string(value.CompetiosCompetitionKey)) == "" || strings.TrimSpace(string(value.CompetiosEntryKey)) == "" || strings.TrimSpace(string(value.CompetiosInviteeKey)) == "" || (value.InvitationState != AttendanceInvitationActive && value.InvitationState != AttendanceInvitationRevoked) {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.Response != nil && !value.Response.IsValid() {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.RespondedAt != nil && value.Response == nil {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

// CompetiosAttendanceService is the narrow provider port used by Competios.
// Eventius remains the sole authority for attendance, invitations and RSVP
// submission; Competios can only ensure/revoke/cancel and query the safe
// projection. Implementations must make both Ensure operations idempotent.
//
// GetAttendanceStatus is retained for existing registration-only callers. It
// MUST fail rather than choose a status if that registration can name multiple
// invitees. New integrations that need an event-only status or one invitation
// MUST use CompetiosAttendanceInviteeStatusService.
type CompetiosAttendanceService interface {
	EnsureAttendanceEvent(ctx context.Context, servicePrincipalID string, request EnsureAttendanceEventRequest) (AttendanceStatusProjection, error)
	EnsureAttendanceInvitation(ctx context.Context, servicePrincipalID string, request EnsureAttendanceInvitationRequest) (AttendanceStatusProjection, error)
	GetAttendanceStatus(ctx context.Context, servicePrincipalID string, eventKey CompetiosEventKey, registrationKey CompetiosRegistrationKey) (AttendanceStatusProjection, error)
	RevokeAttendanceInvitation(ctx context.Context, servicePrincipalID string, invitationID AttendanceInvitationID, reason string) (AttendanceStatusProjection, error)
	CancelAttendanceEvent(ctx context.Context, servicePrincipalID string, eventID AttendanceEventID, reason string) (AttendanceStatusProjection, error)
}

// CompetiosAttendanceInviteeStatusService is an additive provider capability;
// it deliberately does not change CompetiosAttendanceService, so existing
// providers remain source-compatible. GetAttendanceEventStatus returns an
// event-only projection with no invitation tuple. GetAttendanceInviteeStatus
// returns the projection for exactly the complete external tuple supplied in
// request; it must validate that tuple and must never select another invitee.
type CompetiosAttendanceInviteeStatusService interface {
	CompetiosAttendanceService
	GetAttendanceEventStatus(ctx context.Context, servicePrincipalID string, eventKey CompetiosEventKey) (AttendanceStatusProjection, error)
	GetAttendanceInviteeStatus(ctx context.Context, servicePrincipalID string, request GetAttendanceInviteeStatusRequest) (AttendanceStatusProjection, error)
}
