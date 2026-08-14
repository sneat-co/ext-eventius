/**
 * Safe Eventius attendance projection for an external Competios Event. This
 * contract deliberately has no RSVP token, share URL, invitee name, email, or
 * other transport authority.
 */
export type CompetiosAttendanceResponse = 'yes' | 'no' | 'maybe';

export type CompetiosAttendanceEventState = 'active' | 'cancelled';

export type CompetiosAttendanceInvitationState = 'active' | 'revoked';

export type CompetiosAttendanceResponderKind = 'account' | 'guardian';

export interface ICompetiosAttendanceResponderRef {
  readonly kind: CompetiosAttendanceResponderKind;
  readonly accountID: string;
}

/** Matches the frozen CalendarEventRef TypeSpec model. */
export interface ICompetiosAttendanceCalendarEventRef {
  readonly spaceID: string;
  readonly happeningID: string;
}

export interface IEnsureCompetiosAttendanceEventRequest {
  readonly requestID: string;
  readonly competiosEventKey: string;
  readonly calendarEvent: ICompetiosAttendanceCalendarEventRef;
}

export interface IEnsureCompetiosAttendanceInvitationRequest {
  readonly requestID: string;
  readonly attendanceEventID: string;
  readonly competiosRegistrationKey: string;
  readonly competiosTournamentKey: string;
  readonly competiosCompetitionKey: string;
  readonly competiosEntryKey: string;
  /**
   * Legacy shape: providers must reject it because it cannot identify one
   * invitee lifecycle. Use IEnsureCompetiosAttendanceInviteeInvitationRequest.
   */
  readonly responder: ICompetiosAttendanceResponderRef;
}

/** Exact idempotent ensure command for one invitee lifecycle. */
export interface IEnsureCompetiosAttendanceInviteeInvitationRequest {
  readonly requestID: string;
  readonly attendanceEventID: string;
  /** The provider must verify this correlates with attendanceEventID. */
  readonly competiosEventKey: string;
  readonly competiosTournamentKey: string;
  readonly competiosCompetitionKey: string;
  readonly competiosEntryKey: string;
  readonly competiosRegistrationKey: string;
  /** Opaque identity of one invitee. */
  readonly competiosInviteeKey: string;
  /** Opaque Competios revision for this Entry invitation lifecycle. */
  readonly competiosEntryLifecycleRevision: string;
  readonly responder: ICompetiosAttendanceResponderRef;
}

/** Complete, safe lookup key for precisely one invitee invitation status. */
export interface IGetCompetiosAttendanceInviteeStatusRequest {
  readonly competiosEventKey: string;
  readonly competiosTournamentKey: string;
  readonly competiosCompetitionKey: string;
  readonly competiosEntryKey: string;
  readonly competiosRegistrationKey: string;
  readonly competiosInviteeKey: string;
  readonly competiosEntryLifecycleRevision: string;
}

/** Exact, auditable revoke command for one invitee lifecycle. */
export interface IRevokeCompetiosAttendanceInvitationCommand
  extends IGetCompetiosAttendanceInviteeStatusRequest {
  readonly requestID: string;
  readonly attendanceEventID: string;
  readonly attendanceInvitationID: string;
  readonly reason: string;
}

/** Exact, auditable cancel command for one attendance bridge. */
export interface ICancelCompetiosAttendanceEventCommand {
  readonly requestID: string;
  readonly attendanceEventID: string;
  readonly competiosEventKey: string;
  readonly reason: string;
}

export interface ICompetiosAttendanceStatus {
  readonly competiosEventKey: string;
  readonly competiosRegistrationKey?: string;
  readonly competiosTournamentKey?: string;
  readonly competiosCompetitionKey?: string;
  readonly competiosEntryKey?: string;
  readonly competiosInviteeKey?: string;
  readonly competiosEntryLifecycleRevision?: string;
  readonly attendanceEventID: string;
  readonly attendanceInvitationID?: string;
  readonly eventState: CompetiosAttendanceEventState;
  readonly invitationState?: CompetiosAttendanceInvitationState;
  readonly response?: CompetiosAttendanceResponse;
  readonly respondedAt?: string;
}
