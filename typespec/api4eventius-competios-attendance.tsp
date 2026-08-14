/**
 * Frozen safe wire shapes for the Competios-to-Eventius attendance provider
 * facade. The concrete provider transport is deliberately out of scope: this
 * public extension contract defines models only and does not implement an
 * Eventius consumer.
 *
 * The matching Go facade is ../backend/facade4eventius/competios_attendance.go
 * and the published TypeScript mirror is
 * ../frontend/libs/extensions/eventius/contract/src/lib/models/competios-attendance.ts.
 */
namespace Eventius.CompetiosAttendance;

/** Canonical opaque identifier/key: valid UTF-8, 1..128 UTF-8 bytes, with no
 * leading or trailing Unicode whitespace. Implementations measure bytes; the
 * TypeSpec character decorators are an additional transport-level bound. */
@minLength(1)
@maxLength(128)
scalar AttendanceID extends string;

/** Canonical audit reason: valid UTF-8, 1..512 UTF-8 bytes, with no leading or
 * trailing Unicode whitespace. Implementations preserve exact bytes. */
@minLength(1)
@maxLength(512)
scalar AttendanceReason extends string;

/** Lowercase hexadecimal SHA-256 over the canonical full command payload. */
@minLength(64)
@maxLength(64)
@pattern("^[0-9a-f]{64}$")
scalar AttendanceCommandPayloadFingerprint extends string;

enum AttendanceEventState { active, cancelled }
enum AttendanceInvitationState { active, revoked }
enum AttendanceResponse { yes, no, maybe }
enum AttendanceResponderKind { account, guardian }
enum AttendanceCommandOperation {
  ensure_attendance_event,
  ensure_attendance_invitee_invitation,
  revoke_attendance_invitation,
  cancel_attendance_event
}
enum AttendanceCommandErrorCode { command_conflict }

/** Safe error vocabulary for changed-payload or cross-method RequestID reuse. */
model AttendanceCommandError {
  code: AttendanceCommandErrorCode;
}

model CalendarEventRef {
  spaceID: AttendanceID;
  happeningID: AttendanceID;
}

model AttendanceResponderRef {
  kind: AttendanceResponderKind;
  accountID: AttendanceID;
}

/** Ensures one Eventius attendance occasion for the already scheduled event. */
model EnsureAttendanceEventRequest {
  requestID: AttendanceID;
  competiosEventKey: AttendanceID;
  calendarEvent: CalendarEventRef;
}

/**
 * Retained legacy shape. Providers reject it because no invitee/lifecycle
 * identity is present; use EnsureAttendanceInviteeInvitationRequest instead.
 */
model EnsureAttendanceInvitationRequest {
  requestID: AttendanceID;
  attendanceEventID: AttendanceID;
  competiosRegistrationKey: AttendanceID;
  competiosTournamentKey: AttendanceID;
  competiosCompetitionKey: AttendanceID;
  competiosEntryKey: AttendanceID;
  responder: AttendanceResponderRef;
}

/** Exact idempotent ensure command for one invitee lifecycle. The provider
 * verifies that competiosEventKey is correlated to attendanceEventID. */
model EnsureAttendanceInviteeInvitationRequest {
  requestID: AttendanceID;
  attendanceEventID: AttendanceID;
  competiosEventKey: AttendanceID;
  competiosTournamentKey: AttendanceID;
  competiosCompetitionKey: AttendanceID;
  competiosEntryKey: AttendanceID;
  competiosRegistrationKey: AttendanceID;
  competiosInviteeKey: AttendanceID;
  /** Opaque Competios lifecycle revision. Eventius never parses it. */
  competiosEntryLifecycleRevision: AttendanceID;
  responder: AttendanceResponderRef;
}

/** Complete key for one invitation lookup; no contact, token, or payment data. */
model GetAttendanceInviteeStatusRequest {
  competiosEventKey: AttendanceID;
  competiosTournamentKey: AttendanceID;
  competiosCompetitionKey: AttendanceID;
  competiosEntryKey: AttendanceID;
  competiosRegistrationKey: AttendanceID;
  competiosInviteeKey: AttendanceID;
  competiosEntryLifecycleRevision: AttendanceID;
}

/** Exact durable revoke command. Before mutation the provider atomically proves
 * attendanceEventID↔competiosEventKey correlation, the stored invitation's
 * parent==attendanceEventID, and exact equality of the complete lifecycle tuple.
 * Any mismatch fails closed without mutation or audit. */
model RevokeAttendanceInvitationCommand {
  requestID: AttendanceID;
  attendanceEventID: AttendanceID;
  attendanceInvitationID: AttendanceID;
  competiosEventKey: AttendanceID;
  competiosTournamentKey: AttendanceID;
  competiosCompetitionKey: AttendanceID;
  competiosEntryKey: AttendanceID;
  competiosRegistrationKey: AttendanceID;
  competiosInviteeKey: AttendanceID;
  competiosEntryLifecycleRevision: AttendanceID;
  reason: AttendanceReason;
}

/** Exact durable cancel command. It suppresses active attendance invitations
 * but never creates, updates, or cancels the canonical Calendarius Happening. */
model CancelAttendanceEventCommand {
  requestID: AttendanceID;
  attendanceEventID: AttendanceID;
  competiosEventKey: AttendanceID;
  reason: AttendanceReason;
}

/** Event-only projections omit every invitation tuple field; an invitation
 * projection contains the complete external tuple. answered responses have a
 * non-zero UTC respondedAt; revoked answered responses remain valid. */
model AttendanceStatusProjection {
  competiosEventKey: AttendanceID;
  competiosTournamentKey?: AttendanceID;
  competiosCompetitionKey?: AttendanceID;
  competiosEntryKey?: AttendanceID;
  competiosRegistrationKey?: AttendanceID;
  competiosInviteeKey?: AttendanceID;
  competiosEntryLifecycleRevision?: AttendanceID;
  attendanceEventID: AttendanceID;
  attendanceInvitationID?: AttendanceID;
  eventState: AttendanceEventState;
  invitationState?: AttendanceInvitationState;
  response?: AttendanceResponse;
  respondedAt?: utcDateTime;
}

/**
 * Normative durable semantics: exact methods bind (authenticated service
 * principal ID, requestID) globally across all four operation names. The full
 * payload fingerprint is SHA-256 over the operation followed by every model
 * field in declaration order, each encoded as uint32 big-endian UTF-8 byte
 * length plus exact UTF-8 bytes. There is no JSON, delimiter, trimming, Unicode
 * normalization, or omission. A byte-identical replay returns the originally
 * recorded safe projection and writes no mutation/audit; a changed payload,
 * target, reason, or cross-method reuse returns command_conflict and writes
 * nothing. Binding, mutation, audit, and result projection are atomic.
 */
model AttendanceCommandBinding {
  servicePrincipalID: AttendanceID;
  requestID: AttendanceID;
  `operation`: AttendanceCommandOperation;
  `payloadFingerprint`: AttendanceCommandPayloadFingerprint;
  `projection`: AttendanceStatusProjection;
}
