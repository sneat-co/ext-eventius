package facade4eventius

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/sneat-co/ext-eventius/backend/participation"
)

var (
	// ErrInvalidCompetiosAttendanceRequest means a caller supplied an incomplete
	// or unsafe cross-product attendance request. Implementations must reject it
	// before allocating any Eventius record.
	ErrInvalidCompetiosAttendanceRequest = errors.New("eventius: invalid Competios attendance request")

	// ErrLegacyCompetiosAttendanceEnsureUnsupported prevents a legacy ensure
	// request from silently collapsing multiple invitee lifecycles into one
	// invitation. Callers must use the additive exact-invitee provider capability.
	ErrLegacyCompetiosAttendanceEnsureUnsupported = errors.New("eventius: legacy Competios attendance ensure is unsupported")

	// ErrAmbiguousCompetiosAttendanceLookup prevents the legacy
	// registration-only lookup from selecting an arbitrary invitee.
	ErrAmbiguousCompetiosAttendanceLookup = errors.New("eventius: ambiguous Competios attendance lookup")

	// ErrLegacyCompetiosAttendanceMutationUnsupported prevents legacy revoke
	// and cancel methods from claiming durable caller-request idempotency when
	// their source-compatible signatures do not carry a RequestID.
	ErrLegacyCompetiosAttendanceMutationUnsupported = errors.New("eventius: legacy Competios attendance mutation is unsupported")
)

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
// enrolled Competios Entry. Eventius stores and compares it but never parses it.
type CompetiosInviteeKey string

// CompetiosEntryLifecycleRevision is an opaque Competios value that changes
// whenever an Entry lifecycle creates a distinct invitation lifecycle. It is
// required separately from CompetiosInviteeKey so the idempotency tuple is
// explicit and independently auditable.
type CompetiosEntryLifecycleRevision string

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

// EnsureAttendanceInvitationRequest is the retained legacy request shape.
// It cannot identify one invitee lifecycle, so providers MUST fail closed with
// ErrLegacyCompetiosAttendanceEnsureUnsupported. New callers must use
// EnsureAttendanceInviteeInvitationRequest through
// CompetiosAttendanceInviteeStatusService.
type EnsureAttendanceInvitationRequest struct {
	RequestID                string                   `json:"requestID"`
	AttendanceEventID        AttendanceEventID        `json:"attendanceEventID"`
	CompetiosRegistrationKey CompetiosRegistrationKey `json:"competiosRegistrationKey"`
	CompetiosTournamentKey   CompetiosTournamentKey   `json:"competiosTournamentKey"`
	CompetiosCompetitionKey  CompetiosCompetitionKey  `json:"competiosCompetitionKey"`
	CompetiosEntryKey        CompetiosEntryKey        `json:"competiosEntryKey"`
	Responder                AttendanceResponderRef   `json:"responder"`
}

// EnsureAttendanceInviteeInvitationRequest makes one attendance invitation
// for exactly one enrolled invitee lifecycle. The provider MUST verify that
// CompetiosEventKey is correlated to AttendanceEventID before creating or
// returning an invitation; a mismatched pair is rejected. No RSVP token,
// contact, or payment value belongs to this cross-product contract.
type EnsureAttendanceInviteeInvitationRequest struct {
	RequestID                       string                          `json:"requestID"`
	AttendanceEventID               AttendanceEventID               `json:"attendanceEventID"`
	CompetiosEventKey               CompetiosEventKey               `json:"competiosEventKey"`
	CompetiosTournamentKey          CompetiosTournamentKey          `json:"competiosTournamentKey"`
	CompetiosCompetitionKey         CompetiosCompetitionKey         `json:"competiosCompetitionKey"`
	CompetiosEntryKey               CompetiosEntryKey               `json:"competiosEntryKey"`
	CompetiosRegistrationKey        CompetiosRegistrationKey        `json:"competiosRegistrationKey"`
	CompetiosInviteeKey             CompetiosInviteeKey             `json:"competiosInviteeKey"`
	CompetiosEntryLifecycleRevision CompetiosEntryLifecycleRevision `json:"competiosEntryLifecycleRevision"`
	Responder                       AttendanceResponderRef          `json:"responder"`
}

// GetAttendanceInviteeStatusRequest identifies exactly one safe invitation
// status projection. It has no token, contact, payment, or response-detail
// fields. Eventius MUST NOT parse its opaque Competios values.
type GetAttendanceInviteeStatusRequest struct {
	CompetiosEventKey               CompetiosEventKey               `json:"competiosEventKey"`
	CompetiosTournamentKey          CompetiosTournamentKey          `json:"competiosTournamentKey"`
	CompetiosCompetitionKey         CompetiosCompetitionKey         `json:"competiosCompetitionKey"`
	CompetiosEntryKey               CompetiosEntryKey               `json:"competiosEntryKey"`
	CompetiosRegistrationKey        CompetiosRegistrationKey        `json:"competiosRegistrationKey"`
	CompetiosInviteeKey             CompetiosInviteeKey             `json:"competiosInviteeKey"`
	CompetiosEntryLifecycleRevision CompetiosEntryLifecycleRevision `json:"competiosEntryLifecycleRevision"`
}

// RevokeAttendanceInvitationCommand is the exact, durable command for one
// invitation lifecycle. AttendanceInvitationID and the full external tuple
// bind the target twice so a provider cannot revoke an arbitrary invitation.
// Reason is retained in Eventius's audit trail and RequestID controls replay
// and conflict detection.
type RevokeAttendanceInvitationCommand struct {
	RequestID                       string                          `json:"requestID"`
	AttendanceEventID               AttendanceEventID               `json:"attendanceEventID"`
	AttendanceInvitationID          AttendanceInvitationID          `json:"attendanceInvitationID"`
	CompetiosEventKey               CompetiosEventKey               `json:"competiosEventKey"`
	CompetiosTournamentKey          CompetiosTournamentKey          `json:"competiosTournamentKey"`
	CompetiosCompetitionKey         CompetiosCompetitionKey         `json:"competiosCompetitionKey"`
	CompetiosEntryKey               CompetiosEntryKey               `json:"competiosEntryKey"`
	CompetiosRegistrationKey        CompetiosRegistrationKey        `json:"competiosRegistrationKey"`
	CompetiosInviteeKey             CompetiosInviteeKey             `json:"competiosInviteeKey"`
	CompetiosEntryLifecycleRevision CompetiosEntryLifecycleRevision `json:"competiosEntryLifecycleRevision"`
	Reason                          string                          `json:"reason"`
}

// CancelAttendanceEventCommand is the exact, durable command for a canonical
// attendance bridge. CompetiosEventKey must correlate to AttendanceEventID;
// this command never creates or modifies a Calendarius Happening.
type CancelAttendanceEventCommand struct {
	RequestID         string            `json:"requestID"`
	AttendanceEventID AttendanceEventID `json:"attendanceEventID"`
	CompetiosEventKey CompetiosEventKey `json:"competiosEventKey"`
	Reason            string            `json:"reason"`
}

// AttendanceStatusProjection is safe to return to Competios. It intentionally
// omits invitation URLs, RSVP tokens, invitee names, contact fields, and any
// other authority-bearing data.
type AttendanceStatusProjection struct {
	CompetiosEventKey               CompetiosEventKey               `json:"competiosEventKey"`
	CompetiosRegistrationKey        CompetiosRegistrationKey        `json:"competiosRegistrationKey,omitempty"`
	CompetiosTournamentKey          CompetiosTournamentKey          `json:"competiosTournamentKey,omitempty"`
	CompetiosCompetitionKey         CompetiosCompetitionKey         `json:"competiosCompetitionKey,omitempty"`
	CompetiosEntryKey               CompetiosEntryKey               `json:"competiosEntryKey,omitempty"`
	CompetiosInviteeKey             CompetiosInviteeKey             `json:"competiosInviteeKey,omitempty"`
	CompetiosEntryLifecycleRevision CompetiosEntryLifecycleRevision `json:"competiosEntryLifecycleRevision,omitempty"`
	AttendanceEventID               AttendanceEventID               `json:"attendanceEventID"`
	AttendanceInvitationID          AttendanceInvitationID          `json:"attendanceInvitationID,omitempty"`
	EventState                      AttendanceEventState            `json:"eventState"`
	InvitationState                 AttendanceInvitationState       `json:"invitationState,omitempty"`
	Response                        *participation.Coarse           `json:"response,omitempty"`
	RespondedAt                     *time.Time                      `json:"respondedAt,omitempty"`
}

func ValidateEnsureAttendanceEventRequest(value EnsureAttendanceEventRequest) error {
	if strings.TrimSpace(value.RequestID) == "" || strings.TrimSpace(string(value.CompetiosEventKey)) == "" || strings.TrimSpace(value.CalendarEvent.SpaceID) == "" || strings.TrimSpace(value.CalendarEvent.HappeningID) == "" {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

// ValidateEnsureAttendanceInvitationRequest is intentionally fail-closed;
// validating legacy field presence cannot make its omitted invitee lifecycle
// identity safe.
func ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest) error {
	return ErrLegacyCompetiosAttendanceEnsureUnsupported
}

func ValidateEnsureAttendanceInviteeInvitationRequest(value EnsureAttendanceInviteeInvitationRequest) error {
	if strings.TrimSpace(value.RequestID) == "" || strings.TrimSpace(string(value.AttendanceEventID)) == "" || !hasCompleteInviteeTuple(value.CompetiosEventKey, value.CompetiosTournamentKey, value.CompetiosCompetitionKey, value.CompetiosEntryKey, value.CompetiosRegistrationKey, value.CompetiosInviteeKey, value.CompetiosEntryLifecycleRevision) || strings.TrimSpace(value.Responder.AccountID) == "" {
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
	if !hasCompleteInviteeTuple(value.CompetiosEventKey, value.CompetiosTournamentKey, value.CompetiosCompetitionKey, value.CompetiosEntryKey, value.CompetiosRegistrationKey, value.CompetiosInviteeKey, value.CompetiosEntryLifecycleRevision) {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

func ValidateRevokeAttendanceInvitationCommand(value RevokeAttendanceInvitationCommand) error {
	if strings.TrimSpace(value.RequestID) == "" || strings.TrimSpace(string(value.AttendanceEventID)) == "" || strings.TrimSpace(string(value.AttendanceInvitationID)) == "" || strings.TrimSpace(value.Reason) == "" || !hasCompleteInviteeTuple(value.CompetiosEventKey, value.CompetiosTournamentKey, value.CompetiosCompetitionKey, value.CompetiosEntryKey, value.CompetiosRegistrationKey, value.CompetiosInviteeKey, value.CompetiosEntryLifecycleRevision) {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

func ValidateCancelAttendanceEventCommand(value CancelAttendanceEventCommand) error {
	if strings.TrimSpace(value.RequestID) == "" || strings.TrimSpace(string(value.AttendanceEventID)) == "" || strings.TrimSpace(string(value.CompetiosEventKey)) == "" || strings.TrimSpace(value.Reason) == "" {
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
	if strings.TrimSpace(string(value.AttendanceInvitationID)) == "" {
		if hasAnyInviteeTuple(value) || value.InvitationState != "" || value.Response != nil || value.RespondedAt != nil {
			return ErrInvalidCompetiosAttendanceRequest
		}
		return nil
	}
	if !hasCompleteInviteeTuple(value.CompetiosEventKey, value.CompetiosTournamentKey, value.CompetiosCompetitionKey, value.CompetiosEntryKey, value.CompetiosRegistrationKey, value.CompetiosInviteeKey, value.CompetiosEntryLifecycleRevision) || (value.InvitationState != AttendanceInvitationActive && value.InvitationState != AttendanceInvitationRevoked) {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.EventState == AttendanceEventCancelled && value.InvitationState == AttendanceInvitationActive {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.Response != nil && !value.Response.IsValid() {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if (value.Response == nil) != (value.RespondedAt == nil) {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.RespondedAt != nil && (value.RespondedAt.IsZero() || value.RespondedAt.Location() != time.UTC) {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

func hasCompleteInviteeTuple(event CompetiosEventKey, tournament CompetiosTournamentKey, competition CompetiosCompetitionKey, entry CompetiosEntryKey, registration CompetiosRegistrationKey, invitee CompetiosInviteeKey, lifecycleRevision CompetiosEntryLifecycleRevision) bool {
	return strings.TrimSpace(string(event)) != "" && strings.TrimSpace(string(tournament)) != "" && strings.TrimSpace(string(competition)) != "" && strings.TrimSpace(string(entry)) != "" && strings.TrimSpace(string(registration)) != "" && strings.TrimSpace(string(invitee)) != "" && strings.TrimSpace(string(lifecycleRevision)) != ""
}

func hasAnyInviteeTuple(value AttendanceStatusProjection) bool {
	return value.CompetiosRegistrationKey != "" || value.CompetiosTournamentKey != "" || value.CompetiosCompetitionKey != "" || value.CompetiosEntryKey != "" || value.CompetiosInviteeKey != "" || value.CompetiosEntryLifecycleRevision != ""
}

// CompetiosAttendanceService is the retained narrow provider port used by
// existing callers. Eventius remains the sole authority for attendance,
// invitations and RSVP submission.
//
// EnsureAttendanceInvitation MUST return
// ErrLegacyCompetiosAttendanceEnsureUnsupported. RevokeAttendanceInvitation
// and CancelAttendanceEvent MUST return
// ErrLegacyCompetiosAttendanceMutationUnsupported. GetAttendanceStatus MUST
// return ErrAmbiguousCompetiosAttendanceLookup whenever its registration can
// name zero or multiple invitee lifecycles. New callers MUST use the additive
// CompetiosAttendanceInviteeStatusService capability.
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
// event-only projection with no invitation tuple. The exact ensure and lookup
// methods require the complete external tuple and never select another invitee.
type CompetiosAttendanceInviteeStatusService interface {
	CompetiosAttendanceService
	EnsureAttendanceInviteeInvitation(ctx context.Context, servicePrincipalID string, request EnsureAttendanceInviteeInvitationRequest) (AttendanceStatusProjection, error)
	GetAttendanceEventStatus(ctx context.Context, servicePrincipalID string, eventKey CompetiosEventKey) (AttendanceStatusProjection, error)
	GetAttendanceInviteeStatus(ctx context.Context, servicePrincipalID string, request GetAttendanceInviteeStatusRequest) (AttendanceStatusProjection, error)
}

// CompetiosAttendanceCommandService is the additive exact-command capability.
// Unlike the retained legacy mutation methods, every state-changing operation
// carries a caller supplied RequestID and a sufficient target correlation for
// durable replay/conflict handling.
type CompetiosAttendanceCommandService interface {
	CompetiosAttendanceInviteeStatusService
	RevokeAttendanceInvitationCommand(ctx context.Context, servicePrincipalID string, command RevokeAttendanceInvitationCommand) (AttendanceStatusProjection, error)
	CancelAttendanceEventCommand(ctx context.Context, servicePrincipalID string, command CancelAttendanceEventCommand) (AttendanceStatusProjection, error)
}
