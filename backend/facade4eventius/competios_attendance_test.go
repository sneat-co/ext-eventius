package facade4eventius

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sneat-co/ext-eventius/backend/participation"
)

func TestLegacyEnsureFailsClosed(t *testing.T) {
	err := ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest{})
	if !errors.Is(err, ErrLegacyCompetiosAttendanceEnsureUnsupported) {
		t.Fatalf("legacy ensure error = %v", err)
	}
}

func TestAttendanceStatusProjectionHasNoTokenContactOrPaymentFields(t *testing.T) {
	for _, forbidden := range []string{"Token", "Contact", "Payment"} {
		for i := range reflect.TypeFor[AttendanceStatusProjection]().NumField() {
			if strings.Contains(reflect.TypeFor[AttendanceStatusProjection]().Field(i).Name, forbidden) {
				t.Fatalf("safe projection must not expose %s fields", forbidden)
			}
		}
	}
}

func TestExactInviteeRequestsAndProjectionValidate(t *testing.T) {
	if err := ValidateEnsureAttendanceEventRequest(validEnsureAttendanceEventRequest()); err != nil {
		t.Fatalf("valid event request: %v", err)
	}
	if err := ValidateEnsureAttendanceInviteeInvitationRequest(validEnsureAttendanceInviteeInvitationRequest()); err != nil {
		t.Fatalf("valid exact ensure request: %v", err)
	}
	if err := ValidateGetAttendanceInviteeStatusRequest(validGetAttendanceInviteeStatusRequest()); err != nil {
		t.Fatalf("valid exact status query: %v", err)
	}
	if err := ValidateRevokeAttendanceInvitationCommand(validRevokeAttendanceInvitationCommand()); err != nil {
		t.Fatalf("valid revoke command: %v", err)
	}
	if err := ValidateCancelAttendanceEventCommand(validCancelAttendanceEventCommand()); err != nil {
		t.Fatalf("valid cancel command: %v", err)
	}
	answer, at := participation.CoarseYes, time.Now().UTC()
	status := validAttendanceStatusProjection()
	status.Response, status.RespondedAt = &answer, &at
	if err := ValidateAttendanceStatusProjection(status); err != nil {
		t.Fatalf("valid answered projection: %v", err)
	}
}

func TestExactAttendanceMutationCommandsRejectEveryRequiredField(t *testing.T) {
	for _, field := range []string{"request", "attendance event", "attendance invitation", "event", "tournament", "competition", "entry", "registration", "invitee", "lifecycle revision", "reason"} {
		t.Run("revoke rejects "+field, func(t *testing.T) {
			command := validRevokeAttendanceInvitationCommand()
			setRevokeField(&command, field, " ")
			assertInvalid(t, ValidateRevokeAttendanceInvitationCommand(command))
		})
	}
	for _, field := range []string{"request", "attendance event", "event", "reason"} {
		t.Run("cancel rejects "+field, func(t *testing.T) {
			command := validCancelAttendanceEventCommand()
			setCancelField(&command, field, " ")
			assertInvalid(t, ValidateCancelAttendanceEventCommand(command))
		})
	}
}

func TestAttendanceValidatorsRejectEventAndResponderBoundaries(t *testing.T) {
	for name, request := range map[string]EnsureAttendanceEventRequest{
		"blank request":   {CompetiosEventKey: "event", CalendarEvent: CalendarEventRef{SpaceID: "space", HappeningID: "happening"}},
		"blank event":     {RequestID: "request", CalendarEvent: CalendarEventRef{SpaceID: "space", HappeningID: "happening"}},
		"blank space":     {RequestID: "request", CompetiosEventKey: "event", CalendarEvent: CalendarEventRef{HappeningID: "happening"}},
		"blank happening": {RequestID: "request", CompetiosEventKey: "event", CalendarEvent: CalendarEventRef{SpaceID: "space"}},
	} {
		t.Run("event "+name, func(t *testing.T) {
			assertInvalid(t, ValidateEnsureAttendanceEventRequest(request))
		})
	}

	for name, modify := range map[string]func(*EnsureAttendanceInviteeInvitationRequest){
		"blank request":     func(v *EnsureAttendanceInviteeInvitationRequest) { v.RequestID = " " },
		"blank account":     func(v *EnsureAttendanceInviteeInvitationRequest) { v.Responder.AccountID = " " },
		"unknown responder": func(v *EnsureAttendanceInviteeInvitationRequest) { v.Responder.Kind = "delegate" },
	} {
		t.Run("exact ensure "+name, func(t *testing.T) {
			request := validEnsureAttendanceInviteeInvitationRequest()
			modify(&request)
			assertInvalid(t, ValidateEnsureAttendanceInviteeInvitationRequest(request))
		})
	}
}

func TestExactInviteeTupleRejectsEveryMissingField(t *testing.T) {
	for _, name := range []string{"attendance event", "event", "tournament", "competition", "entry", "registration", "invitee", "lifecycle revision"} {
		t.Run("ensure rejects "+name, func(t *testing.T) {
			request := validEnsureAttendanceInviteeInvitationRequest()
			setEnsureField(&request, name, " ")
			assertInvalid(t, ValidateEnsureAttendanceInviteeInvitationRequest(request))
		})
	}
	for _, name := range []string{"event", "tournament", "competition", "entry", "registration", "invitee", "lifecycle revision"} {
		t.Run("query rejects "+name, func(t *testing.T) {
			query := validGetAttendanceInviteeStatusRequest()
			setQueryField(&query, name, " ")
			assertInvalid(t, ValidateGetAttendanceInviteeStatusRequest(query))
		})
		t.Run("projection rejects "+name, func(t *testing.T) {
			projection := validAttendanceStatusProjection()
			setProjectionField(&projection, name, " ")
			assertInvalid(t, ValidateAttendanceStatusProjection(projection))
		})
	}
}

func TestAttendanceStatusProjectionStateAndResponseBoundaries(t *testing.T) {
	answer, utc := participation.CoarseYes, time.Now().UTC()
	local := utc.In(time.FixedZone("review", 3600))
	zero := time.Time{}

	for name, projection := range map[string]AttendanceStatusProjection{
		"event only is valid":                   {CompetiosEventKey: "event", AttendanceEventID: "eventius-event", EventState: AttendanceEventCancelled},
		"answered revoked is valid":             answeredProjection(AttendanceEventActive, AttendanceInvitationRevoked, answer, utc),
		"answered cancelled revoked is valid":   answeredProjection(AttendanceEventCancelled, AttendanceInvitationRevoked, answer, utc),
		"cancelled revoked unanswered is valid": withStates(AttendanceEventCancelled, AttendanceInvitationRevoked),
		"cancelled active invitation":           withStates(AttendanceEventCancelled, AttendanceInvitationActive),
		"response without timestamp":            withResponse(&answer, nil),
		"timestamp without response":            withResponse(nil, &utc),
		"response with zero timestamp":          withResponse(&answer, &zero),
		"response with non UTC timestamp":       withResponse(&answer, &local),
		"event only with invitation tuple": func() AttendanceStatusProjection {
			p := AttendanceStatusProjection{CompetiosEventKey: "event", AttendanceEventID: "eventius-event", EventState: AttendanceEventActive}
			p.CompetiosInviteeKey = "invitee"
			return p
		}(),
		"invitation missing lifecycle revision": func() AttendanceStatusProjection {
			p := validAttendanceStatusProjection()
			p.CompetiosEntryLifecycleRevision = ""
			return p
		}(),
		"unknown response": func() AttendanceStatusProjection {
			p := validAttendanceStatusProjection()
			p.Response, p.RespondedAt = ptr(participation.Coarse("later")), &utc
			return p
		}(),
		"unknown event state": func() AttendanceStatusProjection {
			p := validAttendanceStatusProjection()
			p.EventState = "gone"
			return p
		}(),
		"unknown invitation state": func() AttendanceStatusProjection {
			p := validAttendanceStatusProjection()
			p.InvitationState = "sent"
			return p
		}(),
		"blank attendance event": func() AttendanceStatusProjection {
			p := validAttendanceStatusProjection()
			p.AttendanceEventID = " "
			return p
		}(),
	} {
		t.Run(name, func(t *testing.T) {
			err := ValidateAttendanceStatusProjection(projection)
			switch name {
			case "event only is valid", "answered revoked is valid", "answered cancelled revoked is valid", "cancelled revoked unanswered is valid":
				if err != nil {
					t.Fatalf("ValidateAttendanceStatusProjection() = %v", err)
				}
			default:
				assertInvalid(t, err)
			}
		})
	}
}

func validEnsureAttendanceEventRequest() EnsureAttendanceEventRequest {
	return EnsureAttendanceEventRequest{RequestID: "request-event", CompetiosEventKey: "event", CalendarEvent: CalendarEventRef{SpaceID: "space", HappeningID: "happening"}}
}

func validEnsureAttendanceInviteeInvitationRequest() EnsureAttendanceInviteeInvitationRequest {
	return EnsureAttendanceInviteeInvitationRequest{RequestID: "request-invitee", AttendanceEventID: "eventius-event", CompetiosEventKey: "event", CompetiosTournamentKey: "tournament", CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry", CompetiosRegistrationKey: "registration", CompetiosInviteeKey: "invitee", CompetiosEntryLifecycleRevision: "revision-2", Responder: AttendanceResponderRef{Kind: AttendanceResponderAccount, AccountID: "account"}}
}

func validGetAttendanceInviteeStatusRequest() GetAttendanceInviteeStatusRequest {
	return GetAttendanceInviteeStatusRequest{CompetiosEventKey: "event", CompetiosTournamentKey: "tournament", CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry", CompetiosRegistrationKey: "registration", CompetiosInviteeKey: "invitee", CompetiosEntryLifecycleRevision: "revision-2"}
}

func validRevokeAttendanceInvitationCommand() RevokeAttendanceInvitationCommand {
	return RevokeAttendanceInvitationCommand{RequestID: "revoke-1", AttendanceEventID: "eventius-event", AttendanceInvitationID: "invitation", CompetiosEventKey: "event", CompetiosTournamentKey: "tournament", CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry", CompetiosRegistrationKey: "registration", CompetiosInviteeKey: "invitee", CompetiosEntryLifecycleRevision: "revision-2", Reason: "entry withdrawn"}
}

func validCancelAttendanceEventCommand() CancelAttendanceEventCommand {
	return CancelAttendanceEventCommand{RequestID: "cancel-1", AttendanceEventID: "eventius-event", CompetiosEventKey: "event", Reason: "competition cancelled"}
}

func validAttendanceStatusProjection() AttendanceStatusProjection {
	return AttendanceStatusProjection{CompetiosEventKey: "event", CompetiosTournamentKey: "tournament", CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry", CompetiosRegistrationKey: "registration", CompetiosInviteeKey: "invitee", CompetiosEntryLifecycleRevision: "revision-2", AttendanceEventID: "eventius-event", AttendanceInvitationID: "invitation", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive}
}

func answeredProjection(eventState AttendanceEventState, invitationState AttendanceInvitationState, answer participation.Coarse, at time.Time) AttendanceStatusProjection {
	p := withStates(eventState, invitationState)
	p.Response, p.RespondedAt = &answer, &at
	return p
}

func withStates(eventState AttendanceEventState, invitationState AttendanceInvitationState) AttendanceStatusProjection {
	p := validAttendanceStatusProjection()
	p.EventState, p.InvitationState = eventState, invitationState
	return p
}

func withResponse(answer *participation.Coarse, at *time.Time) AttendanceStatusProjection {
	p := validAttendanceStatusProjection()
	p.Response, p.RespondedAt = answer, at
	return p
}

func setEnsureField(value *EnsureAttendanceInviteeInvitationRequest, name, replacement string) {
	switch name {
	case "attendance event":
		value.AttendanceEventID = AttendanceEventID(replacement)
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
	case "lifecycle revision":
		value.CompetiosEntryLifecycleRevision = CompetiosEntryLifecycleRevision(replacement)
	}
}

func setQueryField(value *GetAttendanceInviteeStatusRequest, name, replacement string) {
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
	case "lifecycle revision":
		value.CompetiosEntryLifecycleRevision = CompetiosEntryLifecycleRevision(replacement)
	}
}

func setRevokeField(value *RevokeAttendanceInvitationCommand, name, replacement string) {
	switch name {
	case "request":
		value.RequestID = replacement
	case "attendance event":
		value.AttendanceEventID = AttendanceEventID(replacement)
	case "attendance invitation":
		value.AttendanceInvitationID = AttendanceInvitationID(replacement)
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
	case "lifecycle revision":
		value.CompetiosEntryLifecycleRevision = CompetiosEntryLifecycleRevision(replacement)
	case "reason":
		value.Reason = replacement
	}
}

func setCancelField(value *CancelAttendanceEventCommand, name, replacement string) {
	switch name {
	case "request":
		value.RequestID = replacement
	case "attendance event":
		value.AttendanceEventID = AttendanceEventID(replacement)
	case "event":
		value.CompetiosEventKey = CompetiosEventKey(replacement)
	case "reason":
		value.Reason = replacement
	}
}

func setProjectionField(value *AttendanceStatusProjection, name, replacement string) {
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
	case "lifecycle revision":
		value.CompetiosEntryLifecycleRevision = CompetiosEntryLifecycleRevision(replacement)
	}
}

func assertInvalid(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, ErrInvalidCompetiosAttendanceRequest) {
		t.Fatalf("error = %v, want ErrInvalidCompetiosAttendanceRequest", err)
	}
}

func ptr(value participation.Coarse) *participation.Coarse { return &value }
