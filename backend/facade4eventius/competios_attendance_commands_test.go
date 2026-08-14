// Copyright 2026 Sneat.app

package facade4eventius

import (
	"errors"
	"strings"
	"testing"
)

func TestAttendanceIdentifiersHaveFiniteCanonicalUTF8ByteBounds(t *testing.T) {
	maxID := strings.Repeat("x", CompetiosAttendanceIDMaxBytes)
	tooLongID := maxID + "x"
	invalidUTF8 := string([]byte{0xff})

	for name, testCase := range map[string]boundedFieldValidator{
		"service principal": {
			valid:   ValidateCompetiosAttendanceServicePrincipalID,
			invalid: ValidateCompetiosAttendanceServicePrincipalID,
		},
		"event request requestID": requestEventFieldValidator(func(value *EnsureAttendanceEventRequest, replacement string) { value.RequestID = replacement }),
		"calendar spaceID": requestEventFieldValidator(func(value *EnsureAttendanceEventRequest, replacement string) {
			value.CalendarEvent.SpaceID = replacement
		}),
		"calendar happeningID": requestEventFieldValidator(func(value *EnsureAttendanceEventRequest, replacement string) {
			value.CalendarEvent.HappeningID = replacement
		}),
		"Competios event key": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.CompetiosEventKey = CompetiosEventKey(replacement)
		}),
		"Competios tournament key": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.CompetiosTournamentKey = CompetiosTournamentKey(replacement)
		}),
		"Competios competition key": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.CompetiosCompetitionKey = CompetiosCompetitionKey(replacement)
		}),
		"Competios entry key": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.CompetiosEntryKey = CompetiosEntryKey(replacement)
		}),
		"Competios registration key": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.CompetiosRegistrationKey = CompetiosRegistrationKey(replacement)
		}),
		"Competios invitee key": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.CompetiosInviteeKey = CompetiosInviteeKey(replacement)
		}),
		"Competios lifecycle revision": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.CompetiosEntryLifecycleRevision = CompetiosEntryLifecycleRevision(replacement)
		}),
		"attendance event ID": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.AttendanceEventID = AttendanceEventID(replacement)
		}),
		"responder account ID": exactInviteeFieldValidator(func(value *EnsureAttendanceInviteeInvitationRequest, replacement string) {
			value.Responder.AccountID = replacement
		}),
		"attendance invitation ID": revokeFieldValidator(func(value *RevokeAttendanceInvitationCommand, replacement string) {
			value.AttendanceInvitationID = AttendanceInvitationID(replacement)
		}),
	} {
		t.Run(name+" accepts max bytes", func(t *testing.T) {
			if err := testCase.valid(maxID); err != nil {
				t.Fatalf("max-byte value rejected: %v", err)
			}
		})
		for invalidName, value := range map[string]string{
			"max plus one":        tooLongID,
			"invalid UTF-8":       invalidUTF8,
			"leading whitespace":  " " + maxID[:1],
			"trailing whitespace": maxID[:1] + " ",
		} {
			t.Run(name+" rejects "+invalidName, func(t *testing.T) {
				assertInvalid(t, testCase.invalid(value))
			})
		}
	}

	multibyteMax := strings.Repeat("🙂", CompetiosAttendanceIDMaxBytes/4)
	if len(multibyteMax) != CompetiosAttendanceIDMaxBytes || ValidateCompetiosAttendanceServicePrincipalID(multibyteMax) != nil {
		t.Fatal("exact UTF-8 byte limit must accept a multibyte value at 128 bytes")
	}
	assertInvalid(t, ValidateCompetiosAttendanceServicePrincipalID(multibyteMax+"🙂"))
}

func TestAttendanceReasonHasFiniteCanonicalUTF8ByteBounds(t *testing.T) {
	for name, reason := range map[string]string{
		"max":           strings.Repeat("r", CompetiosAttendanceReasonMaxBytes),
		"multibyte max": strings.Repeat("🙂", CompetiosAttendanceReasonMaxBytes/4),
	} {
		t.Run(name, func(t *testing.T) {
			command := validRevokeAttendanceInvitationCommand()
			command.Reason = reason
			if err := ValidateRevokeAttendanceInvitationCommand(command); err != nil {
				t.Fatalf("reason at byte limit rejected: %v", err)
			}
		})
	}
	for name, reason := range map[string]string{
		"max plus one":        strings.Repeat("r", CompetiosAttendanceReasonMaxBytes+1),
		"multibyte over max":  strings.Repeat("🙂", CompetiosAttendanceReasonMaxBytes/4+1),
		"invalid UTF-8":       string([]byte{0xff}),
		"leading whitespace":  " withdrawn",
		"trailing whitespace": "withdrawn ",
	} {
		t.Run(name, func(t *testing.T) {
			command := validRevokeAttendanceInvitationCommand()
			command.Reason = reason
			assertInvalid(t, ValidateRevokeAttendanceInvitationCommand(command))
		})
	}
}

func TestCanonicalAttendanceCommandFingerprintIsStableAndFullPayload(t *testing.T) {
	command := validRevokeAttendanceInvitationCommand()
	fingerprint, err := FingerprintRevokeAttendanceInvitationCommand(command)
	if err != nil {
		t.Fatal(err)
	}
	const expected AttendanceCommandPayloadFingerprint = "7289e945623db4bdbf7a1b41e3b622bd9735fe9f45e5f0cc11bd6ec7de3d8345"
	if fingerprint != expected {
		t.Fatalf("fingerprint = %q, want %q", fingerprint, expected)
	}

	mutations := map[string]func(*RevokeAttendanceInvitationCommand){
		"request":       func(v *RevokeAttendanceInvitationCommand) { v.RequestID += "-changed" },
		"event ID":      func(v *RevokeAttendanceInvitationCommand) { v.AttendanceEventID += "-changed" },
		"invitation ID": func(v *RevokeAttendanceInvitationCommand) { v.AttendanceInvitationID += "-changed" },
		"event":         func(v *RevokeAttendanceInvitationCommand) { v.CompetiosEventKey += "-changed" },
		"tournament":    func(v *RevokeAttendanceInvitationCommand) { v.CompetiosTournamentKey += "-changed" },
		"competition":   func(v *RevokeAttendanceInvitationCommand) { v.CompetiosCompetitionKey += "-changed" },
		"entry":         func(v *RevokeAttendanceInvitationCommand) { v.CompetiosEntryKey += "-changed" },
		"registration":  func(v *RevokeAttendanceInvitationCommand) { v.CompetiosRegistrationKey += "-changed" },
		"invitee":       func(v *RevokeAttendanceInvitationCommand) { v.CompetiosInviteeKey += "-changed" },
		"lifecycle":     func(v *RevokeAttendanceInvitationCommand) { v.CompetiosEntryLifecycleRevision += "-changed" },
		"reason":        func(v *RevokeAttendanceInvitationCommand) { v.Reason += " changed" },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			changed := command
			mutate(&changed)
			changedFingerprint, changedErr := FingerprintRevokeAttendanceInvitationCommand(changed)
			if changedErr != nil {
				t.Fatal(changedErr)
			}
			if changedFingerprint == fingerprint {
				t.Fatal("changed full-payload field must change fingerprint")
			}
		})
	}
}

func TestAttendanceCommandBindingReplayAndConflictConformance(t *testing.T) {
	command := validRevokeAttendanceInvitationCommand()
	fingerprint, err := FingerprintRevokeAttendanceInvitationCommand(command)
	if err != nil {
		t.Fatal(err)
	}
	projection := validAttendanceStatusProjection()
	projection.InvitationState = AttendanceInvitationRevoked
	binding, err := NewAttendanceCommandBinding("competios-production", command.RequestID, AttendanceCommandRevokeInvitation, fingerprint, projection)
	if err != nil {
		t.Fatal(err)
	}

	replayed, err := ResolveAttendanceCommandReplay(binding, binding.ServicePrincipalID, binding.RequestID, binding.Operation, binding.PayloadFingerprint)
	if err != nil || replayed.AttendanceInvitationID != projection.AttendanceInvitationID {
		t.Fatalf("byte-identical replay = (%+v, %v)", replayed, err)
	}

	for name, operationAndFingerprint := range map[string]struct {
		operation   AttendanceCommandOperation
		fingerprint AttendanceCommandPayloadFingerprint
	}{
		"changed payload":    {binding.Operation, AttendanceCommandPayloadFingerprint(strings.Repeat("0", 64))},
		"cross-method reuse": {AttendanceCommandCancelEvent, binding.PayloadFingerprint},
	} {
		t.Run(name, func(t *testing.T) {
			_, conflictErr := ResolveAttendanceCommandReplay(binding, binding.ServicePrincipalID, binding.RequestID, operationAndFingerprint.operation, operationAndFingerprint.fingerprint)
			if !errors.Is(conflictErr, ErrCompetiosAttendanceCommandConflict) {
				t.Fatalf("error = %v, want conflict sentinel", conflictErr)
			}
			var typed *CompetiosAttendanceCommandConflictError
			if !errors.As(conflictErr, &typed) || typed.ServicePrincipalID != binding.ServicePrincipalID || typed.RequestID != binding.RequestID {
				t.Fatalf("typed conflict = %#v", typed)
			}
		})
	}
}

func TestAttendanceCommandConformanceReplayCreatesOneMutationAndAudit(t *testing.T) {
	command := validRevokeAttendanceInvitationCommand()
	fingerprint, err := FingerprintRevokeAttendanceInvitationCommand(command)
	if err != nil {
		t.Fatal(err)
	}
	projection := validAttendanceStatusProjection()
	projection.InvitationState = AttendanceInvitationRevoked
	principal := "competios-production"

	var binding *AttendanceCommandBinding
	mutations, audits := 0, 0
	execute := func(operation AttendanceCommandOperation, incoming AttendanceCommandPayloadFingerprint) (AttendanceStatusProjection, error) {
		if binding != nil {
			return ResolveAttendanceCommandReplay(*binding, principal, command.RequestID, operation, incoming)
		}
		mutations++
		audits++
		created, createErr := NewAttendanceCommandBinding(principal, command.RequestID, operation, incoming, projection)
		if createErr != nil {
			return AttendanceStatusProjection{}, createErr
		}
		binding = &created
		return projection, nil
	}

	if _, err = execute(AttendanceCommandRevokeInvitation, fingerprint); err != nil {
		t.Fatal(err)
	}
	if _, err = execute(AttendanceCommandRevokeInvitation, fingerprint); err != nil {
		t.Fatal(err)
	}
	changed := command
	changed.Reason = "different reason"
	changedFingerprint, err := FingerprintRevokeAttendanceInvitationCommand(changed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = execute(AttendanceCommandRevokeInvitation, changedFingerprint); !errors.Is(err, ErrCompetiosAttendanceCommandConflict) {
		t.Fatalf("changed reason error = %v", err)
	}
	if _, err = execute(AttendanceCommandCancelEvent, fingerprint); !errors.Is(err, ErrCompetiosAttendanceCommandConflict) {
		t.Fatalf("cross-method error = %v", err)
	}
	if mutations != 1 || audits != 1 {
		t.Fatalf("mutations=%d audits=%d, want one each", mutations, audits)
	}
}

func TestRevokeTargetRequiresEventParentCorrelationAndFullStoredTuple(t *testing.T) {
	command := validRevokeAttendanceInvitationCommand()
	stored := validAttendanceStatusProjection()
	if err := ValidateRevokeAttendanceInvitationTarget(command, stored); err != nil {
		t.Fatalf("matching target rejected: %v", err)
	}

	mutations := map[string]func(*AttendanceStatusProjection){
		"invitation parent event":     func(v *AttendanceStatusProjection) { v.AttendanceEventID = "other-event" },
		"invitation ID":               func(v *AttendanceStatusProjection) { v.AttendanceInvitationID = "other-invitation" },
		"Competios event correlation": func(v *AttendanceStatusProjection) { v.CompetiosEventKey = "other-event" },
		"tournament":                  func(v *AttendanceStatusProjection) { v.CompetiosTournamentKey = "other-tournament" },
		"competition":                 func(v *AttendanceStatusProjection) { v.CompetiosCompetitionKey = "other-competition" },
		"entry":                       func(v *AttendanceStatusProjection) { v.CompetiosEntryKey = "other-entry" },
		"registration":                func(v *AttendanceStatusProjection) { v.CompetiosRegistrationKey = "other-registration" },
		"invitee":                     func(v *AttendanceStatusProjection) { v.CompetiosInviteeKey = "other-invitee" },
		"lifecycle":                   func(v *AttendanceStatusProjection) { v.CompetiosEntryLifecycleRevision = "other-revision" },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			mismatch := stored
			mutate(&mismatch)
			assertInvalid(t, ValidateRevokeAttendanceInvitationTarget(command, mismatch))
		})
	}
}

func TestCancelTargetRequiresExactEventCorrelation(t *testing.T) {
	command := validCancelAttendanceEventCommand()
	stored := AttendanceStatusProjection{
		CompetiosEventKey: command.CompetiosEventKey,
		AttendanceEventID: command.AttendanceEventID,
		EventState:        AttendanceEventActive,
	}
	if err := ValidateCancelAttendanceEventTarget(command, stored); err != nil {
		t.Fatalf("matching target rejected: %v", err)
	}
	for name, mutate := range map[string]func(*AttendanceStatusProjection){
		"event ID":        func(v *AttendanceStatusProjection) { v.AttendanceEventID = "other-event" },
		"Competios event": func(v *AttendanceStatusProjection) { v.CompetiosEventKey = "other-event" },
		"invitation projection": func(v *AttendanceStatusProjection) {
			*v = validAttendanceStatusProjection()
		},
	} {
		t.Run(name, func(t *testing.T) {
			mismatch := stored
			mutate(&mismatch)
			assertInvalid(t, ValidateCancelAttendanceEventTarget(command, mismatch))
		})
	}
}

type boundedFieldValidator struct {
	valid   func(string) error
	invalid func(string) error
}

func requestEventFieldValidator(set func(*EnsureAttendanceEventRequest, string)) boundedFieldValidator {
	validate := func(replacement string) error {
		value := validEnsureAttendanceEventRequest()
		set(&value, replacement)
		return ValidateEnsureAttendanceEventRequest(value)
	}
	return boundedFieldValidator{valid: validate, invalid: validate}
}

func exactInviteeFieldValidator(set func(*EnsureAttendanceInviteeInvitationRequest, string)) boundedFieldValidator {
	validate := func(replacement string) error {
		value := validEnsureAttendanceInviteeInvitationRequest()
		set(&value, replacement)
		return ValidateEnsureAttendanceInviteeInvitationRequest(value)
	}
	return boundedFieldValidator{valid: validate, invalid: validate}
}

func revokeFieldValidator(set func(*RevokeAttendanceInvitationCommand, string)) boundedFieldValidator {
	validate := func(replacement string) error {
		value := validRevokeAttendanceInvitationCommand()
		set(&value, replacement)
		return ValidateRevokeAttendanceInvitationCommand(value)
	}
	return boundedFieldValidator{valid: validate, invalid: validate}
}
