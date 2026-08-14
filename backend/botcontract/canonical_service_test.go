package botcontract

import (
	"errors"
	"testing"

	"github.com/sneat-co/ext-eventius/backend/facade4eventius"
)

func TestCanonicalAliasesPreserveTheFacadeVocabulary(t *testing.T) {
	request := CreateEventRequest{RequestID: "request", SpaceID: "space", Title: "Night"}
	if facade4eventius.CreateEventRequest(request).Title != "Night" {
		t.Fatal("bot contract must retain the canonical create-event request")
	}
	if ParticipationYes != "yes" || ParticipationMaybe != "maybe" || ParticipationNo != "no" {
		t.Fatal("bot participation values must reuse canonical RSVP values")
	}
	inviteeKey := CompetiosInviteeKey("competios:invitee@entry-revision")
	if facade4eventius.CompetiosInviteeKey(inviteeKey) != "competios:invitee@entry-revision" {
		t.Fatal("bot contract must retain the canonical Competios invitee key")
	}
	lifecycleRevision := CompetiosEntryLifecycleRevision("entry-revision")
	if facade4eventius.CompetiosEntryLifecycleRevision(lifecycleRevision) != "entry-revision" {
		t.Fatal("bot contract must retain the canonical Entry lifecycle revision")
	}
}

func TestCanonicalAttendanceCommandAliasesCompileAgainstFacadeVocabulary(t *testing.T) {
	var commandService CompetiosAttendanceCommandService
	var canonicalService facade4eventius.CompetiosAttendanceCommandService = commandService
	_ = canonicalService

	revoke := RevokeAttendanceInvitationCommand{
		RequestID: "request", AttendanceEventID: "event", AttendanceInvitationID: "invitation",
		CompetiosEventKey: "competios-event", CompetiosTournamentKey: "tournament",
		CompetiosCompetitionKey: "competition", CompetiosEntryKey: "entry",
		CompetiosRegistrationKey: "registration", CompetiosInviteeKey: "invitee",
		CompetiosEntryLifecycleRevision: "revision", Reason: "withdrawn",
	}
	if facade4eventius.RevokeAttendanceInvitationCommand(revoke).Reason != "withdrawn" {
		t.Fatal("bot contract must retain canonical revoke command")
	}
	cancel := CancelAttendanceEventCommand{RequestID: "request", AttendanceEventID: "event", CompetiosEventKey: "competios-event", Reason: "cancelled"}
	if facade4eventius.CancelAttendanceEventCommand(cancel).Reason != "cancelled" {
		t.Fatal("bot contract must retain canonical cancel command")
	}
	binding := AttendanceCommandBinding{Operation: AttendanceCommandRevokeInvitation, PayloadFingerprint: AttendanceCommandPayloadFingerprint("fingerprint")}
	if facade4eventius.AttendanceCommandBinding(binding).Operation != facade4eventius.AttendanceCommandRevokeInvitation {
		t.Fatal("bot contract must retain canonical command binding vocabulary")
	}
	if !errors.Is(ErrCompetiosAttendanceCommandConflict, facade4eventius.ErrCompetiosAttendanceCommandConflict) {
		t.Fatal("bot contract must retain canonical command conflict sentinel")
	}
	if AttendanceCommandErrorCodeConflict != "command_conflict" {
		t.Fatal("bot contract must retain canonical command conflict code")
	}
	var typed *CompetiosAttendanceCommandConflictError
	_ = typed
}
