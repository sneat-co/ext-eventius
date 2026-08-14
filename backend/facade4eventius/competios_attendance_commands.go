// Copyright 2026 Sneat.app

package facade4eventius

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	// CompetiosAttendanceIDMaxBytes applies to RequestID, service principal IDs,
	// Eventius IDs, Calendarius references, responder accounts, and every opaque
	// Competios tuple component. Length is measured in UTF-8 bytes, not runes.
	CompetiosAttendanceIDMaxBytes = 128
	// CompetiosAttendanceReasonMaxBytes is the UTF-8 byte limit for an auditable
	// revoke/cancel reason.
	CompetiosAttendanceReasonMaxBytes = 512
)

var (
	// ErrCompetiosAttendanceCommandConflict is returned when a previously bound
	// (servicePrincipalID, RequestID) is reused with another operation or payload.
	// Providers must return it without performing a mutation or writing an audit.
	ErrCompetiosAttendanceCommandConflict = errors.New("eventius: Competios attendance command conflict")
)

// AttendanceCommandErrorCode is the stable cross-language error vocabulary.
type AttendanceCommandErrorCode string

const AttendanceCommandErrorCodeConflict AttendanceCommandErrorCode = "command_conflict"

// CompetiosAttendanceCommandConflictError is the typed conflict returned for a
// reused exact-command key. Error text deliberately omits opaque identifiers;
// callers can inspect the typed fields or use errors.Is with the sentinel.
type CompetiosAttendanceCommandConflictError struct {
	ServicePrincipalID string
	RequestID          string
}

func (e *CompetiosAttendanceCommandConflictError) Error() string {
	return ErrCompetiosAttendanceCommandConflict.Error()
}

func (e *CompetiosAttendanceCommandConflictError) Unwrap() error {
	return ErrCompetiosAttendanceCommandConflict
}

func (e *CompetiosAttendanceCommandConflictError) Code() AttendanceCommandErrorCode {
	return AttendanceCommandErrorCodeConflict
}

// AttendanceCommandOperation is part of the durable binding. Operation names
// are stable ASCII protocol values and participate in the payload fingerprint.
type AttendanceCommandOperation string

const (
	AttendanceCommandEnsureEvent      AttendanceCommandOperation = "ensure_attendance_event"
	AttendanceCommandEnsureInvitation AttendanceCommandOperation = "ensure_attendance_invitee_invitation"
	AttendanceCommandRevokeInvitation AttendanceCommandOperation = "revoke_attendance_invitation"
	AttendanceCommandCancelEvent      AttendanceCommandOperation = "cancel_attendance_event"
)

// AttendanceCommandPayloadFingerprint is lowercase hexadecimal SHA-256.
type AttendanceCommandPayloadFingerprint string

// AttendanceCommandBinding is the durable result bound globally to one
// (ServicePrincipalID, RequestID) across all exact attendance command methods.
// Providers persist it atomically with the mutation and its audit record.
type AttendanceCommandBinding struct {
	ServicePrincipalID string                              `json:"servicePrincipalID"`
	RequestID          string                              `json:"requestID"`
	Operation          AttendanceCommandOperation          `json:"operation"`
	PayloadFingerprint AttendanceCommandPayloadFingerprint `json:"payloadFingerprint"`
	Projection         AttendanceStatusProjection          `json:"projection"`
}

// ValidateCompetiosAttendanceServicePrincipalID applies the same canonical ID
// policy as command RequestID values: valid UTF-8, 1..128 bytes, no leading or
// trailing Unicode whitespace. Values are never trimmed or normalized.
func ValidateCompetiosAttendanceServicePrincipalID(value string) error {
	if !validAttendanceID(value) {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return nil
}

// FingerprintEnsureAttendanceEventRequest returns the canonical full-payload
// fingerprint after validating the request.
func FingerprintEnsureAttendanceEventRequest(value EnsureAttendanceEventRequest) (AttendanceCommandPayloadFingerprint, error) {
	if err := ValidateEnsureAttendanceEventRequest(value); err != nil {
		return "", err
	}
	return fingerprintAttendanceCommand(AttendanceCommandEnsureEvent,
		value.RequestID,
		string(value.CompetiosEventKey),
		value.CalendarEvent.SpaceID,
		value.CalendarEvent.HappeningID,
	), nil
}

// FingerprintEnsureAttendanceInviteeInvitationRequest returns the canonical
// full-payload fingerprint after validating the exact invitee lifecycle tuple.
func FingerprintEnsureAttendanceInviteeInvitationRequest(value EnsureAttendanceInviteeInvitationRequest) (AttendanceCommandPayloadFingerprint, error) {
	if err := ValidateEnsureAttendanceInviteeInvitationRequest(value); err != nil {
		return "", err
	}
	return fingerprintAttendanceCommand(AttendanceCommandEnsureInvitation,
		value.RequestID,
		string(value.AttendanceEventID),
		string(value.CompetiosEventKey),
		string(value.CompetiosTournamentKey),
		string(value.CompetiosCompetitionKey),
		string(value.CompetiosEntryKey),
		string(value.CompetiosRegistrationKey),
		string(value.CompetiosInviteeKey),
		string(value.CompetiosEntryLifecycleRevision),
		string(value.Responder.Kind),
		value.Responder.AccountID,
	), nil
}

func FingerprintRevokeAttendanceInvitationCommand(value RevokeAttendanceInvitationCommand) (AttendanceCommandPayloadFingerprint, error) {
	if err := ValidateRevokeAttendanceInvitationCommand(value); err != nil {
		return "", err
	}
	return fingerprintAttendanceCommand(AttendanceCommandRevokeInvitation,
		value.RequestID,
		string(value.AttendanceEventID),
		string(value.AttendanceInvitationID),
		string(value.CompetiosEventKey),
		string(value.CompetiosTournamentKey),
		string(value.CompetiosCompetitionKey),
		string(value.CompetiosEntryKey),
		string(value.CompetiosRegistrationKey),
		string(value.CompetiosInviteeKey),
		string(value.CompetiosEntryLifecycleRevision),
		value.Reason,
	), nil
}

func FingerprintCancelAttendanceEventCommand(value CancelAttendanceEventCommand) (AttendanceCommandPayloadFingerprint, error) {
	if err := ValidateCancelAttendanceEventCommand(value); err != nil {
		return "", err
	}
	return fingerprintAttendanceCommand(AttendanceCommandCancelEvent,
		value.RequestID,
		string(value.AttendanceEventID),
		string(value.CompetiosEventKey),
		value.Reason,
	), nil
}

// fingerprintAttendanceCommand is the normative canonicalization algorithm.
// It hashes the operation followed by every model field in declaration order.
// Each value is encoded as a 4-byte unsigned big-endian byte length followed by
// its exact UTF-8 bytes. There is no JSON, delimiter, trimming, Unicode
// normalization, or field omission. The digest is lowercase hexadecimal
// SHA-256. RequestID is intentionally part of the full payload as well as the
// durable binding key.
func fingerprintAttendanceCommand(operation AttendanceCommandOperation, values ...string) AttendanceCommandPayloadFingerprint {
	h := sha256.New()
	writeFingerprintValue(h, string(operation))
	for _, value := range values {
		writeFingerprintValue(h, value)
	}
	return AttendanceCommandPayloadFingerprint(hex.EncodeToString(h.Sum(nil)))
}

type fingerprintWriter interface {
	Write([]byte) (int, error)
}

func writeFingerprintValue(w fingerprintWriter, value string) {
	var size [4]byte
	binary.BigEndian.PutUint32(size[:], uint32(len(value)))
	_, _ = w.Write(size[:])
	_, _ = w.Write([]byte(value))
}

// NewAttendanceCommandBinding validates and constructs the value a provider
// must persist in the same atomic unit as its mutation and audit write.
func NewAttendanceCommandBinding(servicePrincipalID, requestID string, operation AttendanceCommandOperation, fingerprint AttendanceCommandPayloadFingerprint, projection AttendanceStatusProjection) (AttendanceCommandBinding, error) {
	binding := AttendanceCommandBinding{
		ServicePrincipalID: servicePrincipalID,
		RequestID:          requestID,
		Operation:          operation,
		PayloadFingerprint: fingerprint,
		Projection:         cloneAttendanceStatusProjection(projection),
	}
	if err := ValidateAttendanceCommandBinding(binding); err != nil {
		return AttendanceCommandBinding{}, err
	}
	return binding, nil
}

func ValidateAttendanceCommandBinding(value AttendanceCommandBinding) error {
	if ValidateCompetiosAttendanceServicePrincipalID(value.ServicePrincipalID) != nil || !validAttendanceID(value.RequestID) || !value.Operation.IsValid() || !validAttendanceFingerprint(value.PayloadFingerprint) {
		return ErrInvalidCompetiosAttendanceRequest
	}
	return ValidateAttendanceStatusProjection(value.Projection)
}

func (value AttendanceCommandOperation) IsValid() bool {
	switch value {
	case AttendanceCommandEnsureEvent, AttendanceCommandEnsureInvitation, AttendanceCommandRevokeInvitation, AttendanceCommandCancelEvent:
		return true
	default:
		return false
	}
}

// ResolveAttendanceCommandReplay applies the normative replay decision after a
// provider atomically reads the binding at (servicePrincipalID, RequestID).
// An identical operation+fingerprint returns the originally recorded safe
// projection. Any changed payload, target, reason, or cross-method reuse returns
// the typed conflict and must result in no mutation or additional audit.
func ResolveAttendanceCommandReplay(binding AttendanceCommandBinding, servicePrincipalID, requestID string, operation AttendanceCommandOperation, fingerprint AttendanceCommandPayloadFingerprint) (AttendanceStatusProjection, error) {
	if err := ValidateAttendanceCommandBinding(binding); err != nil || ValidateCompetiosAttendanceServicePrincipalID(servicePrincipalID) != nil || !validAttendanceID(requestID) || !operation.IsValid() || !validAttendanceFingerprint(fingerprint) {
		return AttendanceStatusProjection{}, ErrInvalidCompetiosAttendanceRequest
	}
	if binding.ServicePrincipalID != servicePrincipalID || binding.RequestID != requestID {
		return AttendanceStatusProjection{}, ErrInvalidCompetiosAttendanceRequest
	}
	if binding.Operation != operation || binding.PayloadFingerprint != fingerprint {
		return AttendanceStatusProjection{}, &CompetiosAttendanceCommandConflictError{ServicePrincipalID: servicePrincipalID, RequestID: requestID}
	}
	return cloneAttendanceStatusProjection(binding.Projection), nil
}

func cloneAttendanceStatusProjection(value AttendanceStatusProjection) AttendanceStatusProjection {
	clone := value
	if value.Response != nil {
		response := *value.Response
		clone.Response = &response
	}
	if value.RespondedAt != nil {
		respondedAt := *value.RespondedAt
		clone.RespondedAt = &respondedAt
	}
	return clone
}

func validAttendanceID(value string) bool {
	return validBoundedCanonicalUTF8(value, CompetiosAttendanceIDMaxBytes)
}

func validAttendanceReason(value string) bool {
	return validBoundedCanonicalUTF8(value, CompetiosAttendanceReasonMaxBytes)
}

func validBoundedCanonicalUTF8(value string, maxBytes int) bool {
	return utf8.ValidString(value) && len(value) >= 1 && len(value) <= maxBytes && strings.TrimSpace(value) == value
}

func validAttendanceFingerprint(value AttendanceCommandPayloadFingerprint) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	decoded, err := hex.DecodeString(string(value))
	return err == nil && len(decoded) == sha256.Size && strings.ToLower(string(value)) == string(value)
}

func commandCorrelationError(detail string) error {
	return fmt.Errorf("%w: %s", ErrInvalidCompetiosAttendanceRequest, detail)
}
