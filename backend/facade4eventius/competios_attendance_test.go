package facade4eventius

import (
	"errors"
	"testing"
	"time"

	"github.com/sneat-co/ext-eventius/backend/participation"
)

func TestCompetiosAttendanceContractsRejectUnsafeOrIncompleteValues(t *testing.T) {
	validEvent := EnsureAttendanceEventRequest{RequestID: "request-1", CompetiosEventKey: "competios:event-1", CalendarEvent: CalendarEventRef{SpaceID: "space-1", HappeningID: "happening-1"}}
	if err := ValidateEnsureAttendanceEventRequest(validEvent); err != nil {
		t.Fatalf("valid event request: %v", err)
	}
	validInvitation := EnsureAttendanceInvitationRequest{RequestID: "request-2", AttendanceEventID: "attendance-1", CompetiosRegistrationKey: "competios:registration-1", CompetiosTournamentKey: "competios:tournament-1", CompetiosCompetitionKey: "competios:competition-1", CompetiosEntryKey: "competios:entry-1", Responder: AttendanceResponderRef{Kind: AttendanceResponderAccount, AccountID: "user-1"}}
	if err := ValidateEnsureAttendanceInvitationRequest(validInvitation); err != nil {
		t.Fatalf("valid invitation request: %v", err)
	}
	answer := participation.CoarseYes
	validStatus := AttendanceStatusProjection{
		CompetiosEventKey: "competios:event-1", CompetiosRegistrationKey: "competios:registration-1", CompetiosTournamentKey: "competios:tournament-1", CompetiosCompetitionKey: "competios:competition-1", CompetiosEntryKey: "competios:entry-1",
		AttendanceEventID: "attendance-1", AttendanceInvitationID: "attendance-invitation-1",
		EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive, Response: &answer,
	}
	if err := ValidateAttendanceStatusProjection(validStatus); err != nil {
		t.Fatalf("valid safe projection: %v", err)
	}

	for name, err := range map[string]error{
		"blank external event":        ValidateEnsureAttendanceEventRequest(EnsureAttendanceEventRequest{RequestID: "x", CalendarEvent: CalendarEventRef{SpaceID: "space-1", HappeningID: "happening-1"}}),
		"blank calendar event":        ValidateEnsureAttendanceEventRequest(EnsureAttendanceEventRequest{RequestID: "x", CompetiosEventKey: "e"}),
		"blank registration":          ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest{RequestID: "x", AttendanceEventID: "a"}),
		"blank responder account":     ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest{RequestID: "x", AttendanceEventID: "a", CompetiosRegistrationKey: "r", CompetiosTournamentKey: "t", CompetiosCompetitionKey: "c", CompetiosEntryKey: "e", Responder: AttendanceResponderRef{Kind: AttendanceResponderAccount}}),
		"unknown responder kind":      ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest{RequestID: "x", AttendanceEventID: "a", CompetiosRegistrationKey: "r", CompetiosTournamentKey: "t", CompetiosCompetitionKey: "c", CompetiosEntryKey: "e", Responder: AttendanceResponderRef{Kind: "delegate", AccountID: "user-1"}}),
		"response without invitation": ValidateAttendanceStatusProjection(AttendanceStatusProjection{CompetiosEventKey: "e", AttendanceEventID: "a", EventState: AttendanceEventActive, Response: &answer}),
		"unknown response":            ValidateAttendanceStatusProjection(AttendanceStatusProjection{CompetiosEventKey: "e", CompetiosRegistrationKey: "r", AttendanceEventID: "a", AttendanceInvitationID: "i", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive, Response: ptr(participation.Coarse("later"))}),
	} {
		if !errors.Is(err, ErrInvalidCompetiosAttendanceRequest) {
			t.Errorf("%s error = %v, want ErrInvalidCompetiosAttendanceRequest", name, err)
		}
	}
}

func TestAttendanceStatusProjectionValidationBoundaries(t *testing.T) {
	answer := participation.CoarseYes
	at := time.Now().UTC()
	valid := AttendanceStatusProjection{
		CompetiosEventKey: "event", CompetiosRegistrationKey: "registration", CompetiosTournamentKey: "tournament", CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry",
		AttendanceEventID: "eventius-event", AttendanceInvitationID: "invitation", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive,
	}
	for name, projection := range map[string]AttendanceStatusProjection{
		"cancelled event without invitation is valid":   {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", EventState: AttendanceEventCancelled},
		"answered invitation is valid":                  func() AttendanceStatusProjection { p := valid; p.Response, p.RespondedAt = &answer, &at; return p }(),
		"blank event key":                               {AttendanceEventID: "eventius-event", EventState: AttendanceEventActive},
		"blank attendance event":                        {CompetiosEventKey: "event", EventState: AttendanceEventActive},
		"unknown event state":                           {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", EventState: "gone"},
		"event only carries invitation metadata":        {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive},
		"invitation missing external registration data": {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", AttendanceInvitationID: "invitation", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive},
		"unknown invitation state":                      {CompetiosEventKey: "event", CompetiosRegistrationKey: "registration", CompetiosTournamentKey: "tournament", CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry", AttendanceEventID: "eventius-event", AttendanceInvitationID: "invitation", EventState: AttendanceEventActive, InvitationState: "sent"},
		"response timestamp without response":           func() AttendanceStatusProjection { p := valid; p.RespondedAt = &at; return p }(),
	} {
		t.Run(name, func(t *testing.T) {
			err := ValidateAttendanceStatusProjection(projection)
			if name == "cancelled event without invitation is valid" || name == "answered invitation is valid" {
				if err != nil {
					t.Fatalf("ValidateAttendanceStatusProjection() = %v", err)
				}
				return
			}
			if !errors.Is(err, ErrInvalidCompetiosAttendanceRequest) {
				t.Fatalf("ValidateAttendanceStatusProjection() = %v", err)
			}
		})
	}
}

func ptr(value participation.Coarse) *participation.Coarse { return &value }
