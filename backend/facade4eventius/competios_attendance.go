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

// AttendanceEventID and AttendanceInvitationID are Eventius-owned opaque
// identities. They deliberately cannot be substituted for one another.
type AttendanceEventID string
type AttendanceInvitationID string

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

// EnsureAttendanceEventRequest describes the presentation-safe attendance
// event that Eventius should maintain for a Competios Event. RequestID makes a
// retried command idempotent; CompetiosEventKey prevents duplicate events when
// a different delivery attempt retries the same external Event.
type EnsureAttendanceEventRequest struct {
	RequestID         string            `json:"requestID"`
	CompetiosEventKey CompetiosEventKey `json:"competiosEventKey"`
	Title             string            `json:"title"`
	StartsAt          time.Time         `json:"startsAt"`
	Location          string            `json:"location"`
}

// EnsureAttendanceInvitationRequest makes one attendance invitation for a
// confirmed registration. No RSVP token belongs to this cross-product
// contract: any token is Eventius transport authority and must remain private
// to Eventius.
type EnsureAttendanceInvitationRequest struct {
	RequestID                string                   `json:"requestID"`
	AttendanceEventID        AttendanceEventID        `json:"attendanceEventID"`
	CompetiosRegistrationKey CompetiosRegistrationKey `json:"competiosRegistrationKey"`
	CompetiosTournamentKey   CompetiosTournamentKey   `json:"competiosTournamentKey"`
	CompetiosCompetitionKey  CompetiosCompetitionKey  `json:"competiosCompetitionKey"`
	CompetiosEntryKey        CompetiosEntryKey        `json:"competiosEntryKey"`
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
	AttendanceEventID        AttendanceEventID         `json:"attendanceEventID"`
	AttendanceInvitationID   AttendanceInvitationID    `json:"attendanceInvitationID,omitempty"`
	EventState               AttendanceEventState      `json:"eventState"`
	InvitationState          AttendanceInvitationState `json:"invitationState,omitempty"`
	Response                 *participation.Coarse     `json:"response,omitempty"`
	RespondedAt              *time.Time                `json:"respondedAt,omitempty"`
}

func ValidateEnsureAttendanceEventRequest(value EnsureAttendanceEventRequest) error {
	if strings.TrimSpace(value.RequestID) == "" || strings.TrimSpace(string(value.CompetiosEventKey)) == "" || strings.TrimSpace(value.Title) == "" || value.StartsAt.IsZero() || strings.TrimSpace(value.Location) == "" {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

func ValidateEnsureAttendanceInvitationRequest(value EnsureAttendanceInvitationRequest) error {
	if strings.TrimSpace(value.RequestID) == "" || value.AttendanceEventID == "" || strings.TrimSpace(string(value.CompetiosRegistrationKey)) == "" || value.CompetiosTournamentKey == "" || value.CompetiosCompetitionKey == "" || value.CompetiosEntryKey == "" {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

func ValidateAttendanceStatusProjection(value AttendanceStatusProjection) error {
	if value.CompetiosEventKey == "" || value.AttendanceEventID == "" {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.EventState != AttendanceEventActive && value.EventState != AttendanceEventCancelled {
		return ErrInvalidCompetiosAttendanceRequest
	}
	if value.AttendanceInvitationID == "" {
		if value.CompetiosRegistrationKey != "" || value.CompetiosTournamentKey != "" || value.CompetiosCompetitionKey != "" || value.CompetiosEntryKey != "" || value.InvitationState != "" || value.Response != nil || value.RespondedAt != nil {
			return ErrInvalidCompetiosAttendanceRequest
		}
		return nil
	}
	if value.CompetiosRegistrationKey == "" || value.CompetiosTournamentKey == "" || value.CompetiosCompetitionKey == "" || value.CompetiosEntryKey == "" || (value.InvitationState != AttendanceInvitationActive && value.InvitationState != AttendanceInvitationRevoked) {
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
type CompetiosAttendanceService interface {
	EnsureAttendanceEvent(context.Context, string, EnsureAttendanceEventRequest) (AttendanceStatusProjection, error)
	EnsureAttendanceInvitation(context.Context, string, EnsureAttendanceInvitationRequest) (AttendanceStatusProjection, error)
	GetAttendanceStatus(context.Context, string, CompetiosEventKey, CompetiosRegistrationKey) (AttendanceStatusProjection, error)
	RevokeAttendanceInvitation(context.Context, string, AttendanceInvitationID, string) (AttendanceStatusProjection, error)
	CancelAttendanceEvent(context.Context, string, AttendanceEventID, string) (AttendanceStatusProjection, error)
}
