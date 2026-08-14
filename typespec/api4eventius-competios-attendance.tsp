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

model AttendanceResponderRef {
  kind: AttendanceResponderKind;
  accountID: string;
}

model EnsureAttendanceInvitationRequest {
  requestID: string;
  attendanceEventID: string;
  competiosRegistrationKey: string;
  competiosTournamentKey: string;
  competiosCompetitionKey: string;
  competiosEntryKey: string;

  /** Required opaque Competios identity for one invitee lifecycle. It MAY
   * encode the Entry lifecycle revision; Eventius never parses it. */
  competiosInviteeKey: string;
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
}

/** Event-only projections omit every invitation tuple field; an invitation
 * projection contains the complete external tuple. */
model AttendanceStatusProjection {
  competiosEventKey: string;
  competiosTournamentKey?: string;
  competiosCompetitionKey?: string;
  competiosEntryKey?: string;
  competiosRegistrationKey?: string;
  competiosInviteeKey?: string;
  attendanceEventID: string;
  attendanceInvitationID?: string;
  eventState: AttendanceEventState;
  invitationState?: AttendanceInvitationState;
  response?: AttendanceResponse;
  respondedAt?: utcDateTime;
}
