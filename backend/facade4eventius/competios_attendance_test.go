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
	if err := ValidateEnsureAttendanceInvitationRequest(validEnsureAttendanceInvitationRequest()); err != nil {
		t.Fatalf("valid invitation request: %v", err)
	}
	if err := ValidateGetAttendanceInviteeStatusRequest(validGetAttendanceInviteeStatusRequest()); err != nil {
		t.Fatalf("valid invitee status query: %v", err)
	}
	answer := participation.CoarseYes
	validStatus := validAttendanceStatusProjection()
	validStatus.Response = &answer
	if err := ValidateAttendanceStatusProjection(validStatus); err != nil {
		t.Fatalf("valid safe projection: %v", err)
	}

	for name, err := range map[string]error{
		"blank external event":        ValidateEnsureAttendanceEventRequest(EnsureAttendanceEventRequest{RequestID: "x", CalendarEvent: CalendarEventRef{SpaceID: "space-1", HappeningID: "happening-1"}}),
		"blank calendar event":        ValidateEnsureAttendanceEventRequest(EnsureAttendanceEventRequest{RequestID: "x", CompetiosEventKey: "e"}),
		"blank registration":          ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest{RequestID: "x", AttendanceEventID: "a"}),
		"blank responder account":     ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest{RequestID: "x", AttendanceEventID: "a", CompetiosRegistrationKey: "r", CompetiosTournamentKey: "t", CompetiosCompetitionKey: "c", CompetiosEntryKey: "e", CompetiosInviteeKey: "i", Responder: AttendanceResponderRef{Kind: AttendanceResponderAccount}}),
		"unknown responder kind":      ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest{RequestID: "x", AttendanceEventID: "a", CompetiosRegistrationKey: "r", CompetiosTournamentKey: "t", CompetiosCompetitionKey: "c", CompetiosEntryKey: "e", CompetiosInviteeKey: "i", Responder: AttendanceResponderRef{Kind: "delegate", AccountID: "user-1"}}),
		"response without invitation": ValidateAttendanceStatusProjection(AttendanceStatusProjection{CompetiosEventKey: "e", AttendanceEventID: "a", EventState: AttendanceEventActive, Response: &answer}),
		"unknown response":            ValidateAttendanceStatusProjection(AttendanceStatusProjection{CompetiosEventKey: "e", CompetiosRegistrationKey: "r", CompetiosTournamentKey: "t", CompetiosCompetitionKey: "c", CompetiosEntryKey: "e", CompetiosInviteeKey: "i", AttendanceEventID: "a", AttendanceInvitationID: "i", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive, Response: ptr(participation.Coarse("later"))}),
	} {
		if !errors.Is(err, ErrInvalidCompetiosAttendanceRequest) {
			t.Errorf("%s error = %v, want ErrInvalidCompetiosAttendanceRequest", name, err)
		}
	}
}

func TestAttendanceStatusProjectionValidationBoundaries(t *testing.T) {
	answer := participation.CoarseYes
	at := time.Now().UTC()
	valid := validAttendanceStatusProjection()
	for name, projection := range map[string]AttendanceStatusProjection{
		"cancelled event without invitation is valid":   {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", EventState: AttendanceEventCancelled},
		"answered invitation is valid":                  func() AttendanceStatusProjection { p := valid; p.Response, p.RespondedAt = &answer, &at; return p }(),
		"blank event key":                               {AttendanceEventID: "eventius-event", EventState: AttendanceEventActive},
		"blank attendance event":                        {CompetiosEventKey: "event", EventState: AttendanceEventActive},
		"unknown event state":                           {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", EventState: "gone"},
		"event only carries invitation metadata":        {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive},
		"invitation missing external registration data": {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", AttendanceInvitationID: "invitation", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive},
		"unknown invitation state":                      func() AttendanceStatusProjection { p := valid; p.InvitationState = "sent"; return p }(),
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

func TestInviteeAttendanceTupleValidationIsExhaustive(t *testing.T) {
	for name := range map[string]struct{}{"attendance event": {}, "tournament": {}, "competition": {}, "entry": {}, "registration": {}, "invitee": {}} {
		t.Run("ensure rejects missing "+name, func(t *testing.T) {
			request := validEnsureAttendanceInvitationRequest()
			setInvitationTupleField(&request, name, " ")
			if !errors.Is(ValidateEnsureAttendanceInvitationRequest(request), ErrInvalidCompetiosAttendanceRequest) {
				t.Fatal("expected invalid ensure request")
			}
		})
	}

	for name := range map[string]struct{}{"event": {}, "tournament": {}, "competition": {}, "entry": {}, "registration": {}, "invitee": {}} {
		t.Run("query rejects missing "+name, func(t *testing.T) {
			query := validGetAttendanceInviteeStatusRequest()
			setQueryTupleField(&query, name, " ")
			if !errors.Is(ValidateGetAttendanceInviteeStatusRequest(query), ErrInvalidCompetiosAttendanceRequest) {
				t.Fatal("expected invalid invitee query")
			}
		})
		t.Run("projection rejects missing "+name, func(t *testing.T) {
			projection := validAttendanceStatusProjection()
			setProjectionTupleField(&projection, name, " ")
			if !errors.Is(ValidateAttendanceStatusProjection(projection), ErrInvalidCompetiosAttendanceRequest) {
				t.Fatal("expected invalid invitation projection")
			}
		})
	}
}

func validEnsureAttendanceInvitationRequest() EnsureAttendanceInvitationRequest {
	return EnsureAttendanceInvitationRequest{RequestID: "request-2", AttendanceEventID: "attendance-1", CompetiosRegistrationKey: "competios:registration-1", CompetiosTournamentKey: "competios:tournament-1", CompetiosCompetitionKey: "competios:competition-1", CompetiosEntryKey: "competios:entry-1", CompetiosInviteeKey: "competios:invitee-1@revision-1", Responder: AttendanceResponderRef{Kind: AttendanceResponderAccount, AccountID: "user-1"}}
}

func validGetAttendanceInviteeStatusRequest() GetAttendanceInviteeStatusRequest {
	return GetAttendanceInviteeStatusRequest{CompetiosEventKey: "event", CompetiosTournamentKey: "tournament", CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry", CompetiosRegistrationKey: "registration", CompetiosInviteeKey: "invitee@revision-1"}
}

func validAttendanceStatusProjection() AttendanceStatusProjection {
	return AttendanceStatusProjection{CompetiosEventKey: "event", CompetiosTournamentKey: "tournament", CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry", CompetiosRegistrationKey: "registration", CompetiosInviteeKey: "invitee@revision-1", AttendanceEventID: "eventius-event", AttendanceInvitationID: "invitation", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive}
}

func setInvitationTupleField(value *EnsureAttendanceInvitationRequest, name, replacement string) {
	switch name {
	case "attendance event":
		value.AttendanceEventID = AttendanceEventID(replacement)
	case "tournament":
		value.CompetiosTournamentKey = CompetiosTournamentKey(replacement)
	case "competition":
		value.CompetiosCompetitionKey = CompetiosCompetitionKey(replacement)
	case "entry":
		value.CompetiosEntryKey = CompetiosEntryKey(replacement)
	case "registration":
		value.CompetiosRegistrationKey = CompetiosRegistrationKey(replacement)
	case "invitee":
		value.CompetiosInviteeKey = CompetiosInviteeKey(replacement)
	}
}

func setQueryTupleField(value *GetAttendanceInviteeStatusRequest, name, replacement string) {
	switch name {
	case "event":
		value.CompetiosEventKey = CompetiosEventKey(replacement)
	case "tournament":
		value.CompetiosTournamentKey = CompetiosTournamentKey(replacement)
	case "competition":
		value.CompetiosCompetitionKey = CompetiosCompetitionKey(replacement)
	case "entry":
		value.CompetiosEntryKey = CompetiosEntryKey(replacement)
	case "registration":
		value.CompetiosRegistrationKey = CompetiosRegistrationKey(replacement)
	case "invitee":
		value.CompetiosInviteeKey = CompetiosInviteeKey(replacement)
	}
}

func setProjectionTupleField(value *AttendanceStatusProjection, name, replacement string) {
	switch name {
	case "event":
		value.CompetiosEventKey = CompetiosEventKey(replacement)
	case "tournament":
		value.CompetiosTournamentKey = CompetiosTournamentKey(replacement)
	case "competition":
		value.CompetiosCompetitionKey = CompetiosCompetitionKey(replacement)
	case "entry":
		value.CompetiosEntryKey = CompetiosEntryKey(replacement)
	case "registration":
		value.CompetiosRegistrationKey = CompetiosRegistrationKey(replacement)
	case "invitee":
		value.CompetiosInviteeKey = CompetiosInviteeKey(replacement)
	}
}

func ptr(value participation.Coarse) *participation.Coarse { return &value }
