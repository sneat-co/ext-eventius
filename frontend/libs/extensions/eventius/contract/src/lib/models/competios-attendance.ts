/**
 * Safe Eventius attendance projection for an external Competios Event. This
 * contract deliberately has no RSVP token, share URL, invitee name, email, or
 * other transport authority.
 */
export type CompetiosAttendanceResponse = 'yes' | 'no' | 'maybe';

export type CompetiosAttendanceEventState = 'active' | 'cancelled';

export type CompetiosAttendanceInvitationState = 'active' | 'revoked';

export interface IEnsureCompetiosAttendanceEventRequest {
  readonly requestID: string;
  readonly competiosEventKey: string;
  readonly title: string;
  readonly startsAt: string;
  readonly location: string;
}

export interface IEnsureCompetiosAttendanceInvitationRequest {
  readonly requestID: string;
  readonly attendanceEventID: string;
  readonly competiosRegistrationKey: string;
  readonly competiosTournamentKey: string;
  readonly competiosCompetitionKey: string;
  readonly competiosEntryKey: string;
}

export interface ICompetiosAttendanceStatus {
  readonly competiosEventKey: string;
  readonly competiosRegistrationKey?: string;
  readonly competiosTournamentKey?: string;
  readonly competiosCompetitionKey?: string;
  readonly competiosEntryKey?: string;
  readonly attendanceEventID: string;
  readonly attendanceInvitationID?: string;
  readonly eventState: CompetiosAttendanceEventState;
  readonly invitationState?: CompetiosAttendanceInvitationState;
  readonly response?: CompetiosAttendanceResponse;
  readonly respondedAt?: string;
}
