import { describe, expect, expectTypeOf, it } from 'vitest';

import type {
  ICompetiosAttendanceStatus,
  IEnsureCompetiosAttendanceInvitationRequest,
  IGetCompetiosAttendanceInviteeStatusRequest,
} from './competios-attendance';

describe('Competios attendance invitee contract', () => {
  it('requires the opaque invitee lifecycle key on an ensure command', () => {
    const request: IEnsureCompetiosAttendanceInvitationRequest = {
      requestID: 'request-1',
      attendanceEventID: 'eventius-event-1',
      competiosRegistrationKey: 'registration-1',
      competiosTournamentKey: 'tournament-1',
      competiosCompetitionKey: 'competition-1',
      competiosEntryKey: 'entry-1',
      competiosInviteeKey: 'invitee-1@entry-revision-2',
      responder: { kind: 'account', accountID: 'account-1' },
    };

    expect(Object.keys(request)).toEqual([
      'requestID',
      'attendanceEventID',
      'competiosRegistrationKey',
      'competiosTournamentKey',
      'competiosCompetitionKey',
      'competiosEntryKey',
      'competiosInviteeKey',
      'responder',
    ]);
  });

  it('keeps the exact lookup and safe projection tuple aligned', () => {
    const lookup: IGetCompetiosAttendanceInviteeStatusRequest = {
      competiosEventKey: 'event-1',
      competiosTournamentKey: 'tournament-1',
      competiosCompetitionKey: 'competition-1',
      competiosEntryKey: 'entry-1',
      competiosRegistrationKey: 'registration-1',
      competiosInviteeKey: 'invitee-1@entry-revision-2',
    };
    const status: ICompetiosAttendanceStatus = {
      ...lookup,
      attendanceEventID: 'eventius-event-1',
      attendanceInvitationID: 'eventius-invitation-1',
      eventState: 'active',
      invitationState: 'active',
    };

    expect(status.competiosInviteeKey).toBe(lookup.competiosInviteeKey);
    expectTypeOf<Extract<keyof ICompetiosAttendanceStatus, 'token' | 'contact' | 'payment'>>().toEqualTypeOf<never>();
  });
});
