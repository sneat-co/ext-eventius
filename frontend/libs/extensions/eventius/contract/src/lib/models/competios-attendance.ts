/**
 * Safe Eventius attendance projection for an external Competios Event. This
 * contract deliberately has no RSVP token, share URL, invitee name, email, or
 * other transport authority.
 */
export type CompetiosAttendanceResponse = 'yes' | 'no' | 'maybe';

export type CompetiosAttendanceEventState = 'active' | 'cancelled';

export type CompetiosAttendanceInvitationState = 'active' | 'revoked';

export type CompetiosAttendanceResponderKind = 'account' | 'guardian';

/**
 * IDs/keys/RequestIDs are valid UTF-8, 1..128 UTF-8 bytes, with no leading or
 * trailing Unicode whitespace. Implementations preserve exact bytes and never
 * trim or normalize. The alias documents the wire constraint; runtime adapters
 * must validate bytes before calling the provider.
 */
export type CompetiosAttendanceID = string;

/** Audit reasons use the same canonical policy with a 1..512 UTF-8 byte bound. */
export type CompetiosAttendanceReason = string;

export const COMPETIOS_ATTENDANCE_ID_MAX_UTF8_BYTES = 128;
export const COMPETIOS_ATTENDANCE_REASON_MAX_UTF8_BYTES = 512;

export type CompetiosAttendanceCommandOperation =
  | 'ensure_attendance_event'
  | 'ensure_attendance_invitee_invitation'
  | 'revoke_attendance_invitation'
  | 'cancel_attendance_event';

export type CompetiosAttendanceCommandErrorCode = 'command_conflict';
export const COMPETIOS_ATTENDANCE_COMMAND_CONFLICT_CODE =
  'command_conflict' as const;

export interface ICompetiosAttendanceCommandError {
  readonly code: CompetiosAttendanceCommandErrorCode;
}

/** Lowercase hexadecimal SHA-256 (exactly 64 ASCII bytes). */
export type CompetiosAttendanceCommandPayloadFingerprint = string;

export interface ICompetiosAttendanceResponderRef {
  readonly kind: CompetiosAttendanceResponderKind;
  readonly accountID: CompetiosAttendanceID;
}

/** Matches the frozen CalendarEventRef TypeSpec model. */
export interface ICompetiosAttendanceCalendarEventRef {
  readonly spaceID: CompetiosAttendanceID;
  readonly happeningID: CompetiosAttendanceID;
}

export interface IEnsureCompetiosAttendanceEventRequest {
  readonly requestID: CompetiosAttendanceID;
  readonly competiosEventKey: CompetiosAttendanceID;
  readonly calendarEvent: ICompetiosAttendanceCalendarEventRef;
}

export interface IEnsureCompetiosAttendanceInvitationRequest {
  readonly requestID: CompetiosAttendanceID;
  readonly attendanceEventID: CompetiosAttendanceID;
  readonly competiosRegistrationKey: CompetiosAttendanceID;
  readonly competiosTournamentKey: CompetiosAttendanceID;
  readonly competiosCompetitionKey: CompetiosAttendanceID;
  readonly competiosEntryKey: CompetiosAttendanceID;
  /**
   * Legacy shape: providers must reject it because it cannot identify one
   * invitee lifecycle. Use IEnsureCompetiosAttendanceInviteeInvitationRequest.
   */
  readonly responder: ICompetiosAttendanceResponderRef;
}

/** Exact idempotent ensure command for one invitee lifecycle. */
export interface IEnsureCompetiosAttendanceInviteeInvitationRequest {
  readonly requestID: CompetiosAttendanceID;
  readonly attendanceEventID: CompetiosAttendanceID;
  /** The provider must verify this correlates with attendanceEventID. */
  readonly competiosEventKey: CompetiosAttendanceID;
  readonly competiosTournamentKey: CompetiosAttendanceID;
  readonly competiosCompetitionKey: CompetiosAttendanceID;
  readonly competiosEntryKey: CompetiosAttendanceID;
  readonly competiosRegistrationKey: CompetiosAttendanceID;
  /** Opaque identity of one invitee. */
  readonly competiosInviteeKey: CompetiosAttendanceID;
  /** Opaque Competios revision for this Entry invitation lifecycle. */
  readonly competiosEntryLifecycleRevision: CompetiosAttendanceID;
  readonly responder: ICompetiosAttendanceResponderRef;
}

/** Complete, safe lookup key for precisely one invitee invitation status. */
export interface IGetCompetiosAttendanceInviteeStatusRequest {
  readonly competiosEventKey: CompetiosAttendanceID;
  readonly competiosTournamentKey: CompetiosAttendanceID;
  readonly competiosCompetitionKey: CompetiosAttendanceID;
  readonly competiosEntryKey: CompetiosAttendanceID;
  readonly competiosRegistrationKey: CompetiosAttendanceID;
  readonly competiosInviteeKey: CompetiosAttendanceID;
  readonly competiosEntryLifecycleRevision: CompetiosAttendanceID;
}

/**
 * Exact, auditable revoke command. Before mutation, the provider atomically
 * verifies event correlation, invitation parent, and full stored tuple equality.
 */
export interface IRevokeCompetiosAttendanceInvitationCommand
  extends IGetCompetiosAttendanceInviteeStatusRequest {
  readonly requestID: CompetiosAttendanceID;
  readonly attendanceEventID: CompetiosAttendanceID;
  readonly attendanceInvitationID: CompetiosAttendanceID;
  readonly reason: CompetiosAttendanceReason;
}

/** Exact, auditable cancel command for one attendance bridge. */
export interface ICancelCompetiosAttendanceEventCommand {
  readonly requestID: CompetiosAttendanceID;
  readonly attendanceEventID: CompetiosAttendanceID;
  readonly competiosEventKey: CompetiosAttendanceID;
  readonly reason: CompetiosAttendanceReason;
}

export interface ICompetiosAttendanceStatus {
  readonly competiosEventKey: CompetiosAttendanceID;
  readonly competiosRegistrationKey?: CompetiosAttendanceID;
  readonly competiosTournamentKey?: CompetiosAttendanceID;
  readonly competiosCompetitionKey?: CompetiosAttendanceID;
  readonly competiosEntryKey?: CompetiosAttendanceID;
  readonly competiosInviteeKey?: CompetiosAttendanceID;
  readonly competiosEntryLifecycleRevision?: CompetiosAttendanceID;
  readonly attendanceEventID: CompetiosAttendanceID;
  readonly attendanceInvitationID?: CompetiosAttendanceID;
  readonly eventState: CompetiosAttendanceEventState;
  readonly invitationState?: CompetiosAttendanceInvitationState;
  readonly response?: CompetiosAttendanceResponse;
  readonly respondedAt?: string;
}

/**
 * Durable global binding for (authenticated service principal ID, requestID).
 * Fingerprinting is lowercase SHA-256 over operation then every field in model
 * declaration order, encoded as uint32 big-endian UTF-8 byte length plus exact
 * UTF-8 bytes. No JSON, delimiter, trim, normalization, or omission is used.
 */
export interface ICompetiosAttendanceCommandBinding {
  readonly servicePrincipalID: CompetiosAttendanceID;
  readonly requestID: CompetiosAttendanceID;
  readonly operation: CompetiosAttendanceCommandOperation;
  readonly payloadFingerprint: CompetiosAttendanceCommandPayloadFingerprint;
  readonly projection: ICompetiosAttendanceStatus;
}
