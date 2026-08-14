package facade4eventius

import (
	"errors"
	"testing"
	"time"

	"github.com/sneat-co/ext-eventius/backend/participation"
)

func TestCompetiosAttendanceContractsRejectUnsafeOrIncompleteValues(t *testing.T) {
	validEvent := EnsureAttendanceEventRequest{RequestID: "request-1", CompetiosEventKey: "competios:event-1", Title: "Chess cup", StartsAt: time.Now().UTC(), Location: "Clubhouse"}
	if err := ValidateEnsureAttendanceEventRequest(validEvent); err != nil {
		t.Fatalf("valid event request: %v", err)
	}
	validInvitation := EnsureAttendanceInvitationRequest{RequestID: "request-2", AttendanceEventID: "attendance-1", CompetiosRegistrationKey: "competios:registration-1", CompetiosTournamentKey: "competios:tournament-1", CompetiosCompetitionKey: "competios:competition-1", CompetiosEntryKey: "competios:entry-1"}
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
		"blank external event":        ValidateEnsureAttendanceEventRequest(EnsureAttendanceEventRequest{RequestID: "x", Title: "Cup", StartsAt: time.Now().UTC(), Location: "Clubhouse"}),
		"blank registration":          ValidateEnsureAttendanceInvitationRequest(EnsureAttendanceInvitationRequest{RequestID: "x", AttendanceEventID: "a"}),
		"response without invitation": ValidateAttendanceStatusProjection(AttendanceStatusProjection{CompetiosEventKey: "e", AttendanceEventID: "a", EventState: AttendanceEventActive, Response: &answer}),
		"unknown response":            ValidateAttendanceStatusProjection(AttendanceStatusProjection{CompetiosEventKey: "e", CompetiosRegistrationKey: "r", AttendanceEventID: "a", AttendanceInvitationID: "i", EventState: AttendanceEventActive, InvitationState: AttendanceInvitationActive, Response: ptr(participation.Coarse("later"))}),
	} {
		if !errors.Is(err, ErrInvalidCompetiosAttendanceRequest) {
			t.Errorf("%s error = %v, want ErrInvalidCompetiosAttendanceRequest", name, err)
		}
	}
}

func ptr(value participation.Coarse) *participation.Coarse { return &value }
