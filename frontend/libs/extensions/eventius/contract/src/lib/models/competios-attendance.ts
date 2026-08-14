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

export interface IEnsureCompetiosAttendanceEventRequest {
  readonly requestID: string;
  readonly competiosEventKey: string;
  readonly calendarEvent: {
    readonly spaceID: string;
    readonly happeningID: string;
  };
}

export interface IEnsureCompetiosAttendanceInvitationRequest {
  readonly requestID: string;
  readonly attendanceEventID: string;
  readonly competiosRegistrationKey: string;
  readonly competiosTournamentKey: string;
  readonly competiosCompetitionKey: string;
  readonly competiosEntryKey: string;
  /** Opaque identity of one invitee and its Competios lifecycle revision. */
  readonly competiosInviteeKey: string;
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
}

export interface ICompetiosAttendanceStatus {
  readonly competiosEventKey: string;
  readonly competiosRegistrationKey?: string;
  readonly competiosTournamentKey?: string;
  readonly competiosCompetitionKey?: string;
  readonly competiosEntryKey?: string;
  readonly competiosInviteeKey?: string;
  readonly attendanceEventID: string;
  readonly attendanceInvitationID?: string;
  readonly eventState: CompetiosAttendanceEventState;
  readonly invitationState?: CompetiosAttendanceInvitationState;
  readonly response?: CompetiosAttendanceResponse;
  readonly respondedAt?: string;
}
