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

enum AttendanceEventState { active, cancelled }
enum AttendanceInvitationState { active, revoked }
enum AttendanceResponse { yes, no, maybe }
enum AttendanceResponderKind { account, guardian }

model CalendarEventRef {
  spaceID: string;
  happeningID: string;
}

model AttendanceResponderRef {
  kind: AttendanceResponderKind;
  accountID: string;
}

/** Ensures one Eventius attendance occasion for the already scheduled event. */
model EnsureAttendanceEventRequest {
  requestID: string;
  competiosEventKey: string;
  calendarEvent: CalendarEventRef;
}

/**
 * Retained legacy shape. Providers reject it because no invitee/lifecycle
 * identity is present; use EnsureAttendanceInviteeInvitationRequest instead.
 */
model EnsureAttendanceInvitationRequest {
  requestID: string;
  attendanceEventID: string;
  competiosRegistrationKey: string;
  competiosTournamentKey: string;
  competiosCompetitionKey: string;
  competiosEntryKey: string;
  responder: AttendanceResponderRef;
}

/** Exact idempotent ensure command for one invitee lifecycle. The provider
 * verifies that competiosEventKey is correlated to attendanceEventID. */
model EnsureAttendanceInviteeInvitationRequest {
  requestID: string;
  attendanceEventID: string;
  competiosEventKey: string;
  competiosTournamentKey: string;
  competiosCompetitionKey: string;
  competiosEntryKey: string;
  competiosRegistrationKey: string;
  competiosInviteeKey: string;
  /** Opaque Competios lifecycle revision. Eventius never parses it. */
  competiosEntryLifecycleRevision: string;
  responder: AttendanceResponderRef;
}

/** Complete key for one invitation lookup; no contact, token, or payment data. */
model GetAttendanceInviteeStatusRequest {
  competiosEventKey: string;
  competiosTournamentKey: string;
  competiosCompetitionKey: string;
  competiosEntryKey: string;
  competiosRegistrationKey: string;
  competiosInviteeKey: string;
  competiosEntryLifecycleRevision: string;
}

/** Event-only projections omit every invitation tuple field; an invitation
 * projection contains the complete external tuple. answered responses have a
 * non-zero UTC respondedAt; revoked answered responses remain valid. */
model AttendanceStatusProjection {
  competiosEventKey: string;
  competiosTournamentKey?: string;
  competiosCompetitionKey?: string;
  competiosEntryKey?: string;
  competiosRegistrationKey?: string;
  competiosInviteeKey?: string;
  competiosEntryLifecycleRevision?: string;
  attendanceEventID: string;
  attendanceInvitationID?: string;
  eventState: AttendanceEventState;
  invitationState?: AttendanceInvitationState;
  response?: AttendanceResponse;
  respondedAt?: utcDateTime;
}
